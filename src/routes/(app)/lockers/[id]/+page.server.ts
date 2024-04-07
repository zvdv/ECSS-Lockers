import { redirect, error } from "@sveltejs/kit";
import { mustAuthorize } from "$lib/auth.server";
import { defaultExpiry } from "$lib/date";
import { db } from "$lib/db";

export type LockerStatus = "expired" | "available" | "claimed";
export type Locker = {
  name: string;
  locker: string;
  status: LockerStatus;
};

export async function load({ params, cookies }) {
  const { user } = await mustAuthorize(cookies);
  const locker = params.id;

  const result = await db
    .selectFrom("registration")
    .select(["name", "expiry"])
    .where("user", "=", user)
    .where("locker", "=", locker)
    .executeTakeFirst();
  if (result === undefined) {
    throw error(404, { message: "you haven't registered this locker" });
  }
  const { name, expiry } = result;

  let status: LockerStatus = "claimed";
  const now = Date.now();
  if (expiry.valueOf() < now) {
    status = "expired";
  }

  return { name, locker, status };
}

export const actions = {
  renew: async ({ params, cookies }) => {
    const { user } = await mustAuthorize(cookies);
    const locker = params.id;
    db.updateTable("registration")
      .set({ expiry: defaultExpiry() })
      .where("user", "=", user) // SECURITY
      .where("locker", "=", locker)
      .execute();
    throw redirect(302, `./${locker}`);
  },
};
