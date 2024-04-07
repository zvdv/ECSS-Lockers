import { setError, superValidate } from "sveltekit-superforms/server";
import { z } from "zod";
import { fail, redirect } from "@sveltejs/kit";

const formSchema = z.object({
  locker: z.string(),
});

export async function load({ request }) {
  const form = await superValidate(formSchema);
  const url = new URL(request.url);
  form.data.locker = url.searchParams.get("locker") ?? "";
  return { form };
}

export const actions = {
  default: async ({ request, fetch }) => {
    const form = await superValidate(request, formSchema);
    if (!form.valid) {
      return fail(400, { form });
    }
    const { locker } = form.data;
    const resp = await fetch("/admin/api/lockers", {
      method: "DELETE",
      body: JSON.stringify([locker]),
    });
    if (resp.status === 404) {
      return setError(
        form,
        "locker",
        "This locker isn't registered, try another"
      );
    }
    throw redirect(302, `./`);
  },
};
