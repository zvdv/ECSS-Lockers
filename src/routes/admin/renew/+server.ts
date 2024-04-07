import { z } from "zod";
import { defaultExpiry } from "$lib/date";
import { db } from "$lib/db";

const schema = z.object({
  locker: z.string(),
});

export const POST = async ({ request }) => {
  const parsed = schema.safeParse(await request.json());
  if (!parsed.success) {
    return new Response("Missing locker", { status: 400 });
  }
  const { locker } = parsed.data;
  await db
    .updateTable("registration")
    .where("locker", "=", locker)
    .set({ expiry: defaultExpiry() })
    .executeTakeFirstOrThrow();
  return new Response();
};
