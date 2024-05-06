/*
 * this script is deprecated.
 * the seeding action is the `schema.sql` file.
 * I'm keeping this for reference however, once the
 * new site is hosted (and if I remember to) this file will
 * be deleted.
 *
 * - hal
 * */
import { Kysely } from "kysely";
import { PlanetScaleDialect } from "kysely-planetscale";

console.warn(`THIS SCRIPT IS DEPRECATED!!!
THE SEEDING ACTION CAN NOW BE DONE WITHIN THE \`./SCHEMA.SQL\` FILE.
SEE \`README.md\``);

export const db = new Kysely({
	dialect: new PlanetScaleDialect({
		url: process.env.DATABASE_URL,
	}),
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
