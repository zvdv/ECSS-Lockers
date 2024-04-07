import { base } from "$app/paths";
import { redirect } from "@sveltejs/kit";
import { logout, mustAuthorize } from "$lib/auth.server";
import { db } from "$lib/db";

export const actions = {
  default: async ({ params, cookies }) => {
    const { user } = await mustAuthorize(cookies);
    // don't tell the user if this fails, only happens if people mess with the API and they dont deserve errors
    const hasLockers = await db.transaction().execute(async (trx) => {
      await trx
        .deleteFrom("registration")
        .where("user", "=", user)
        .where("locker", "=", params.id)
        .execute();

      const { count } = await trx
        .selectFrom("registration")
        .select([db.fn.countAll<string>().as("count")])
        .where("user", "=", user)
        .executeTakeFirstOrThrow();
      return +count > 0;
    });
    if (!hasLockers) {
      // go get a new one
      logout(cookies);
      throw redirect(302, `${base}`);
    }
    throw redirect(302, "..");
  },
};
