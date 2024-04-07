import { z } from "zod";

const schema = z.array(
  z.union([
    z.object({
      locker: z.string(),
      user: z.string(),
      name: z.string(),
      expiry: z.coerce.date(),
    }),
    z.object({
      locker: z.string(),
      user: z.null(),
      name: z.null(),
      expiry: z.null(),
    }),
  ])
);

export async function load({ fetch }) {
  const resp = await fetch("/admin/api/lockers");
  const data = schema.parse(await resp.json());
  const registration = data.map((data) => {
    let status: "claimed" | "available" | "expired" = "claimed";
    if (data.user === null) {
      status = "available";
      return { ...data, expiry: null, status };
    } else {
      const expiry = data.expiry;
      if (expiry.valueOf() < Date.now()) {
        status = "expired";
      }
      return {
        ...data,
        expiry: new Date(expiry),
        status,
      };
    }
  });

  return { registration };
}
