import { env } from "$env/dynamic/private";
import { Kysely, MysqlDialect } from "kysely";
import { createPool } from "mysql2";
import type { DB } from "../../db/db.d.ts";

const dialect = new MysqlDialect({
    pool: createPool({
        database: env.MYSQL_DATABASE,
        host: env.MYSQL_HOST,
        user: env.MYSQL_USER,
        password: env.MYSQL_PASSWORD,
        port: parseInt(env.MYSQL_PORT),
        connectionLimit: 10,
    }),
})

export const db = new Kysely<DB>({ dialect });
