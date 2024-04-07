export async function load({ url }) {
  const msg = url.searchParams.get("msg") ?? "";
  return { msg };
}
