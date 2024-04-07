import { setError, superValidate } from "sveltekit-superforms/server";
import { z } from "zod";
import { fail, redirect } from "@sveltejs/kit";
import { defaultExpiry, localDateToUTC } from "$lib/date";

const formSchema = z.object({
  name: z.string(),
  email: z.string().email(),
  locker: z.string(),
  expiry: z.date(),
});

export async function load({ request }) {
  const form = await superValidate(formSchema);
  const url = new URL(request.url);
  form.data.locker = url.searchParams.get("locker") ?? "";
  form.data.expiry = defaultExpiry(new Date());
  return { form };
}

export const actions = {
  default: async ({ request, fetch }) => {
    const form = await superValidate(request, formSchema);
    if (!form.valid) {
      return fail(400, { form });
    }
    const expiry = localDateToUTC(form.data.expiry);
    const resp = await fetch("/admin/api/lockers", {
      method: "POST",
      body: JSON.stringify({
        ...form.data,
        expiry,
      }),
    });
    const body = await resp.json();
    const { err } = body;
    if (err === "locker-taken") {
      return setError(
        form,
        "locker",
        "This locker is taken, please register another"
      );
    }
    throw redirect(302, `./`);
  },
};
