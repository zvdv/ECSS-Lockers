import { sql } from "kysely";
import { setError, superValidate } from "sveltekit-superforms/server";
import { z } from "zod";
import { error, fail, redirect } from "@sveltejs/kit";
import { localDateToUTC } from "$lib/date.js";
import { db } from "$lib/db";

const formSchema = z.object({
  name: z.string(),
  email: z.string().email(),
  locker: z.string(),
  expiry: z.coerce.date(),
});

export async function load({ request }) {
  const availableLockers = (
    await db
      .selectFrom("locker")
      .select(["id as locker"])
      .where(({ not, exists, selectFrom }) =>
        not(
          exists(
            selectFrom("registration")
              .select(() => sql`1`.as("one"))
              .whereRef("registration.locker", "=", "locker.id")
          )
        )
      )
      .execute()
  ).map((row) => row.locker);

  const form = await superValidate(formSchema);
  const url = new URL(request.url);
  const locker = url.searchParams.get("locker");
  if (locker === null) {
    return { form, availableLockers };
  }
  form.data.locker = locker;
  availableLockers.push(locker);
  availableLockers.sort();
  const lockerData = await db
    .selectFrom("registration")
    .select(["expiry", "locker", "name", "user"])
    .where("locker", "=", locker)
    .executeTakeFirst();
  if (lockerData === undefined) {
    return { form, availableLockers };
  }
  form.data.email = lockerData.user;
  form.data.name = lockerData.name;
  form.data.expiry = lockerData.expiry;
  return { form, availableLockers };
}

export const actions = {
  default: async ({ request }) => {
    const oldLocker = new URL(request.url).searchParams.get("locker");
    if (oldLocker === null) {
      return error(400, "Missing locker url parameter");
    }
    const form = await superValidate(request, formSchema);
    if (!form.valid) {
      return fail(400, { form });
    }
    const { email, locker, name, expiry: localExpiry } = form.data;
    if (Number.isNaN(localExpiry.valueOf())) {
      return setError(form, "expiry", "Invalid date");
    }

    // convert local date to UTC
    const expiry = localDateToUTC(localExpiry);
    const result = await db.transaction().execute(async (trx) => {
      const matchingLocker = await trx
        .selectFrom("locker")
        .select(["id"])
        .where("id", "=", locker)
        .executeTakeFirst();

      if (matchingLocker === undefined) {
        return "locker-does-not-exist";
      }

      await trx
        .insertInto("user")
        .ignore()
        .columns(["email"])
        .values({ email })
        .execute();

      const owner = await trx
        .selectFrom("registration")
        .select(["user"])
        .where("locker", "=", locker)
        .executeTakeFirst();

      if (owner !== undefined && owner.user == email) {
        await trx
          .updateTable("registration")
          .set({ expiry, name, user: email })
          .where("locker", "=", locker)
          .executeTakeFirstOrThrow();

        return "ok";
      }

      const insertResult = await trx
        .insertInto("registration")
        .ignore()
        .values({ user: email, locker, name, expiry })
        .executeTakeFirstOrThrow();

      if (insertResult.numInsertedOrUpdatedRows === 0n) {
        return "locker-taken";
      }

      const deleteResult = await trx
        .deleteFrom("registration")
        .where("locker", "=", oldLocker)
        .executeTakeFirstOrThrow();

      if (deleteResult.numDeletedRows === 0n) {
        return "cannot-delete";
      }
      return "ok";
    });
    switch (result) {
      case "locker-does-not-exist":
        return setError(form, "locker", "Locker does not exist");
      case "cannot-delete":
        return setError(
          form,
          "locker",
          "Cannot delete locker. Please try again."
        );
      case "locker-taken":
        return setError(form, "locker", "Locker is taken");
    }
    throw redirect(302, `./`);
  },
};
