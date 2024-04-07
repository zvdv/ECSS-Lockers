import { env } from "$env/dynamic/private";
import { jwtVerify, SignJWT } from "jose";
import { z } from "zod";

const magicDataSchema = z.union([
  z.object({
    type: z.literal("register"),
    user: z.string(),
    locker: z.string(),
    name: z.string(),
  }),
  z.object({
    type: z.literal("login"),
    user: z.string(),
  }),
  z.object({
    type: z.literal("transfer"),
    user: z.string(),
    locker: z.string(),
  }),
]);

export type MagicData = z.infer<typeof magicDataSchema>;

const secret = new TextEncoder().encode(env.JWT_SECRET);
const algo = "HS256";

export async function makeMagicToken(data: MagicData): Promise<string> {
  return await new SignJWT(data)
    .setExpirationTime("1d")
    .setProtectedHeader({ alg: algo })
    .setIssuedAt()
    .sign(secret);
}

export async function parseMagicToken(token: string): Promise<MagicData> {
  return magicDataSchema.parse(
    (await jwtVerify(token, secret, { algorithms: [algo] })).payload
  );
}
