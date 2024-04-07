import { sql } from "kysely";
import { message, setError, superValidate } from "sveltekit-superforms/server";
import { z } from "zod";
import { fail } from "@sveltejs/kit";
import { db } from "$lib/db";
import { sendRegisterEmail, sendRegisterLimitEmail } from "$lib/email";

const MAX_REGISTERED = 1;

const formSchema = z.object({
  name: z.string().nonempty(),
  email: z
    .string()
    .email()
    .refine((email) => email.endsWith("@uvic.ca"), {
      message: "Must be a UVic email address",
    })
    .refine((email) => !email.includes("+"), {
      message: "Plus addresses are not allowed",
    }),
  locker: z.string(),
});

export async function load() {
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
  return { form, availableLockers };
}

export const actions = {
  default: async ({ request }) => {
    const form = await superValidate(request, formSchema);
    if (!form.valid) {
      return fail(400, { form });
    }
    const { email, locker, name } = form.data;

    const result = await db.transaction().execute(async (trx) => {
      const result = await sql`
      SELECT 
        CASE
          WHEN NOT EXISTS (
            SELECT 1 FROM locker WHERE id = ${locker}
          ) THEN 'locker-does-not-exist'
          WHEN EXISTS (
            SELECT 1 FROM registration WHERE locker = ${locker} AND expiry > NOW()
          ) THEN 'locker-taken'
          WHEN (
            SELECT COUNT(*) FROM registration WHERE user = ${email} AND expiry > NOW()
          ) >= ${MAX_REGISTERED} THEN 'limit-exceeded'
          ELSE 'ok'
        END AS status
      `.execute(trx);
      const parsed = z
        .array(
          z.object({
            status: z.union([
              z.literal("locker-does-not-exist"),
              z.literal("locker-taken"),
              z.literal("limit-exceeded"),
              z.literal("ok"),
            ]),
          })
        )
        .min(1)
        .parse(result.rows);
      return parsed[0].status;
    });

    switch (result) {
      case "locker-does-not-exist":
        return setError(form, "locker", "This locker does not exist");
      case "locker-taken":
        return setError(
          form,
          "locker",
          "This locker is taken, please register another"
        );
      case "limit-exceeded":
        sendRegisterLimitEmail(email);
        break;
      case "ok":
        sendRegisterEmail(email, locker, name);
        break;
    }
    const msg =
      "Almost done! Check your email for a link to finish registering.";
    return message(form, { msg });
  },
};
