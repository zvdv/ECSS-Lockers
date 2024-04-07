import { env } from "$env/dynamic/private";
import { SignJWT, jwtVerify } from "jose";
import { error, type Cookies } from "@sveltejs/kit";

const cookieName = "admin";

const secret = new TextEncoder().encode(env.JWT_SECRET);
const algo = "HS256";

export async function login(cookies: Cookies) {
  const token = await new SignJWT({})
    .setProtectedHeader({ alg: algo })
    .sign(secret);
  cookies.set(cookieName, token, { path: "/" });
}

export function logout(cookies: Cookies) {
  cookies.set(cookieName, "", { path: "/", maxAge: 0 });
}

export type AuthResult =
  | {
      authorized: false;
    }
  | {
      authorized: true;
    };

export async function authorize(cookies: Cookies): Promise<AuthResult> {
  const token = cookies.get(cookieName);
  if (!token) {
    return { authorized: false };
  }
  try {
    await jwtVerify(token, secret, { algorithms: [algo] });
    return {
      authorized: true,
    };
  } catch (e) {
    return { authorized: false };
  }
}
export async function mustAuthorize(cookies: Cookies) {
  const auth = await authorize(cookies);
  if (!auth.authorized) {
    throw error(401);
  }
}
