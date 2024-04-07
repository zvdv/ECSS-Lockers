import { env } from "$env/dynamic/private";
import { text } from "@sveltejs/kit";

export async function GET() {
  return text(env.ORIGIN);
}
