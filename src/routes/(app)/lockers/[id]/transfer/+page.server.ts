import { superValidate } from "sveltekit-superforms/server";
import { z } from "zod";
import { error, redirect } from "@sveltejs/kit";
import { mustAuthorize } from "$lib/auth.server";
import { db } from "$lib/db";
import { sendTransferEmail } from "$lib/email";

const formSchema = z.object({
  email: z.string().email(),
});

export async function load() {
  const form = await superValidate(formSchema);
  return { form };
}

export const actions = {
  default: async ({ params, cookies, request }) => {
    const { user } = await mustAuthorize(cookies);
    const locker = params.id;

    const form = await superValidate(request, formSchema);
    if (!form.valid) {
      return { form };
    }
    const { email } = form.data;
    const { count } = await db
      .selectFrom("registration")
      .select([db.fn.countAll<string>().as("count")])
      .where("user", "=", user)
      .where("locker", "=", locker)
      .executeTakeFirstOrThrow();

    if (+count === 0) {
      return error(403, { message: "You do not own this locker" });
    }

    sendTransferEmail(email, locker);
    throw redirect(302, "..");
  },
};
