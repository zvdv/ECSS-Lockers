import { redirect } from "@sveltejs/kit";
import { logout, mustAuthorize } from "$lib/auth.server";
import { db } from "$lib/db";

export async function load({ cookies }) {
  const { user } = await mustAuthorize(cookies);

  const result = await db
    .selectFrom("registration")
    .select(["name", "locker"])
    .where("user", "=", user)
    .execute();
  return { locker: result };
}

export const actions = {
  logout: ({ cookies }) => {
    logout(cookies);
    throw redirect(302, "../login");
  },
};
