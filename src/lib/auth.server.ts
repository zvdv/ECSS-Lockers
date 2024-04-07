import { env } from "$env/dynamic/private";
import { SignJWT, jwtVerify } from "jose";
import { z } from "zod";
import { error, type Cookies } from "@sveltejs/kit";

const cookieName = "auth";

const secret = new TextEncoder().encode(env.JWT_SECRET);
const algo = "HS256";

const AuthDataSchema = z.object({
  user: z.string(),
});
export type AuthData = z.infer<typeof AuthDataSchema>;

async function makeToken(data: AuthData) {
  return await new SignJWT(data)
    .setExpirationTime("1d")
    .setProtectedHeader({ alg: algo })
    .setIssuedAt()
    .sign(secret);
}

export function loginCookie(data: AuthData) {
  return `${cookieName}=${makeToken(data)}; Path=/`;
}

export async function login(data: AuthData, cookies: Cookies) {
  cookies.set(cookieName, await makeToken(data), { path: "/" });
}

export function logout(cookies: Cookies) {
  cookies.set(cookieName, "", { path: "/", maxAge: 0 });
}

export type AuthResult =
  | {
      authorized: false;
    }
  | ({
      authorized: true;
    } & AuthData);

export async function authorize(cookies: Cookies): Promise<AuthResult> {
  const token = cookies.get(cookieName);
  if (!token) {
    return { authorized: false };
  }
  try {
    const jsonData = (await jwtVerify(token, secret, { algorithms: [algo] }))
      .payload;
    return {
      ...AuthDataSchema.parse(jsonData),
      authorized: true,
    };
  } catch (e) {
    return { authorized: false };
  }
}
export async function mustAuthorize(cookies: Cookies): Promise<AuthData> {
  const auth = await authorize(cookies);
  if (!auth.authorized) {
    throw error(401);
  }
  return auth;
}
