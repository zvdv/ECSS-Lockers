import { db } from '$lib/db';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	// TODO move to API route and try fetch?
	const registration = await db
		.selectFrom('locker')
		.leftJoin('registration', 'locker.id', 'registration.locker')
		.select(['name', 'user', 'locker.id as locker', 'expiry'])
		.orderBy('locker')
		.execute();
	return { registration };
};
