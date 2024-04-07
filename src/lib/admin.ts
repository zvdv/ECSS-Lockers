import { env } from "$env/dynamic/private";

let password: string | undefined = env.ADMIN_PASSWORD;

export function getAdminPassword() {
  if (!password) {
    throw Error("ADMIN_PASSWORD undefined");
  }
  return password;
}

export function setAdminPassword(newPassword: string) {
  password = newPassword;
}
