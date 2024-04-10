import { base } from "$app/paths";
import { env } from "$env/dynamic/private";
import { SMTPClient } from "emailjs";
import { makeMagicToken } from "./magic";

const client = new SMTPClient({
  host: "smtp.gmail.com",
  port: 465,
  ssl: true,
  user: env.GMAIL_USER,
  password: env.GMAIL_PASSWORD,
});

let prefix = env.ORIGIN;

function makeLink(path: string) {
  return new URL(path, prefix + base);
}

function makeMagicLink(token: string) {
  return makeLink(`magic/${encodeURIComponent(token)}`);
}

export async function sendRegisterEmail(
  email: string,
  locker: string,
  name: string
) {
  const token = await makeMagicToken({
    type: "register",
    user: email,
    locker,
    name,
  });
  const text = `Click here to register locker ${locker} for ${email}: ${makeMagicLink(
    token
  )}`;
  await client.sendAsync({
    from: `Lockers <${env.GMAIL_USER}>`,
    to: email,
    subject: `ECSS Locker Registration: ${locker}`,
    text,
  });
}

export async function sendRegisterLimitEmail(email: string) {
  const text = `You have reached the maximum number of lockers you can register. Please contact the ECSS to request more, or transfer ownership of one of your lockers to another user.
  .`;
  await client.sendAsync({
    from: `Lockers <${env.GMAIL_USER}>`,
    to: email,
    subject: `ECSS Locker Registration Limit Reached`,
    text,
  });
}

export async function sendLoginEmail(email: string) {
  const token = await makeMagicToken({ type: "login", user: email });
  const text = `Click here to login as ${email}: ${makeMagicLink(token)}`;
  await client.sendAsync({
    from: `Lockers <${env.GMAIL_USER}>`,
    to: email,
    subject: `ECSS Locker Login`,
    text,
  });
}

export async function sendAccountDoesNotExistEmail(email: string) {
  const text = `You do not have an account. Please register before trying to log in.`;
  await client.sendAsync({
    from: `Lockers <${env.GMAIL_USER}>`,
    to: email,
    subject: `ECSS Locker Account Does Not Exist`,
    text,
  });
}

export async function sendTransferEmail(email: string, locker: string) {
  const token = await makeMagicToken({ type: "transfer", user: email, locker });
  const text = `You have been sent a locker transfer request for locker ${locker}. Click here to accept: ${makeMagicLink(
    token
  )}`;
  await client.sendAsync({
    from: `Lockers <${env.GMAIL_USER}>`,
    to: email,
    subject: `ECSS Locker Transfer Request: ${locker}`,
    text,
  });
}

export async function sendExpiryEmail(email: string, locker: string) {
  const msg = `Your locker ${locker} has expired. Renew it at ${makeLink(
    `lockers/${locker}`
  )}`;
  await client.sendAsync({
    from: `Lockers <${env.GMAIL_USER}>`,
    to: email,
    subject: `ECSS Locker Expiry: ${locker}`,
    text: msg,
  });
}
