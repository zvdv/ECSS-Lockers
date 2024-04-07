import { env } from "$env/dynamic/private";
import { Kysely } from "kysely";
import { PlanetScaleDialect } from "kysely-planetscale";
import type { DB } from "../../db/db.d.ts";

export const db = new Kysely<DB>({
  dialect: new PlanetScaleDialect({
    url: env.DATABASE_URL,
  }),
});
