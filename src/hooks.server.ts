import { base } from "$app/paths";
import { text, type Handle } from "@sveltejs/kit";
import { authorize as userAuthorize } from "$lib/auth.server";
import { authorize as adminAuthorize } from "./routes/admin/auth.server";

export const handle: Handle = async ({ event, resolve }) => {
  if (
    event.url.pathname.startsWith(`${base}/admin`) &&
    !(await adminAuthorize(event.cookies)).authorized
  ) {
    return text("Please contact the Director of IT to gain access.", {
      status: 401,
    });
  } else if (
    event.url.pathname.startsWith(`${base}/lockers`) &&
    !(await userAuthorize(event.cookies)).authorized
  ) {
    return new Response("You're not logged in", {
      status: 301,
      headers: { location: `${base}/login` },
    });
  }
  return resolve(event);
};
