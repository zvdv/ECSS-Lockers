import express from "express";
// @ts-ignore
import { handler } from "../build/handler";
import "./cron";

const app = express();

app.use(handler);

app.listen(3000, () => {
  console.log("Listening on port 3000");
});
