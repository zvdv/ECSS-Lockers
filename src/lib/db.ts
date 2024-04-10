import { env } from "$env/dynamic/private";
//import { Kysely } from "kysely";
//import { PlanetScaleDialect } from "kysely-planetscale";
//import type { DB } from "../../db/db.d.ts";

/* export const db = new Kysely<DB>({
  dialect: new PlanetScaleDialect({
    url: env.DATABASE_URL,
  }),
}); */

import { createClient } from "@libsql/client";

console.log(env.DATABASE_URL)

export const db = createClient({
  url: env.DATABASE_URL,
  authToken: env.DB_AUTH,
});
