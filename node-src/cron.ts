import { SMTPError } from "emailjs";
import cron from "node-cron";
import { db } from "../src/lib/db.js";
import { sendExpiryEmail } from "../src/lib/email.js";

// https://crontab.guru
cron.schedule("0 * * * *", async () => {
  const expired = await db
    .selectFrom("registration")
    .selectAll()
    .where("expiry", "<", new Date())
    .where(({ cmpr, or, ref }) =>
      or([
        cmpr("expiryEmailSent", "is", null),
        cmpr("expiryEmailSent", "<", ref("expiry")),
      ])
    )
    .execute();

  console.log("expired lockers:", expired);
  for (const { user, locker } of expired) {
    try {
      await sendExpiryEmail(user, locker);
      await db
        .updateTable("registration")
        .set({ expiryEmailSent: new Date() })
        .executeTakeFirstOrThrow();
      console.debug("sent expiry email", {
        user,
        locker,
      });
    } catch (e) {
      if (e instanceof SMTPError) {
        console.error("SMTP error sending expiry email", e);
        continue;
      }
      throw e;
    }
  }
});
