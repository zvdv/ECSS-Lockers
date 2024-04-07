import { z } from "zod";
import { json, type RequestHandler } from "@sveltejs/kit";
import { db } from "$lib/db";

export const GET: RequestHandler = async () => {
  const result = await db
    .selectFrom("locker")
    .leftJoin("registration", "locker.id", "registration.locker")
    .select(["name", "user", "locker.id as locker", "expiry"])
    .orderBy("locker")
    .execute();
  return json(result);
};

export const POST: RequestHandler = async ({ request }) => {
  const schema = z.object({
    locker: z.string(),
    name: z.string(),
    email: z.string().email(),
    expiry: z.coerce.date(),
  });
  const parsed = schema.safeParse(await request.json());
  if (!parsed.success) {
    return json({ err: "invalid" }, { status: 400 });
  }
  const { email, locker, name, expiry } = parsed.data;

  const result = await db.transaction().execute(async (trx) => {
    await trx
      .insertInto("user")
      .ignore()
      .columns(["email"])
      .values({ email })
      .execute();

    const q = trx
      .insertInto("registration")
      .ignore()
      .columns(["user", "locker", "name", "expiry"])
      .values({ user: email, locker, name, expiry });
    const result = await q.executeTakeFirstOrThrow();
    if (result.numInsertedOrUpdatedRows === 0n) {
      return "locker-taken";
    }
    return "ok";
  });
  if (result === "locker-taken") {
    return json({ err: "locker-taken" }, { status: 400 });
  }
  return json({});
};

export const DELETE: RequestHandler = async ({ request }) => {
  const schema = z.array(z.string());
  const parsed = schema.safeParse(await request.json());
  if (!parsed.success) {
    return json({ err: "invalid" }, { status: 400 });
  }
  const lockers = parsed.data;
  const result = await db
    .deleteFrom("registration")
    .where("locker", "in", lockers)
    .executeTakeFirstOrThrow();
  if (result.numDeletedRows === 0n) {
    return json({ err: "not-found" }, { status: 404 });
  }
  return json({});
};
