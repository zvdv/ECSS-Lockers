import { createClient } from "@libsql/client";
import { env } from "$env/dynamic/private";

export const db = createClient({
  url: env.DATABASE_URL,
  authToken: env.DB_AUTH,
});

function range(start, end) {
  return [...Array(1 + end - start).keys()].map((n) => start + n);
}

const lockers = [range(1, 200)]
  .flat()
  .map((x) => `ELW ${x.toString().padStart(3, "0")}`);

console.log(
  (
    await db
      .insertInto("locker")
      .ignore()
      .values(
        lockers.map((x) => ({
          id: x,
        }))
      )
      .executeTakeFirstOrThrow()
  ).numInsertedOrUpdatedRows
);
