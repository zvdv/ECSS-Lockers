import { superForm as defaultSuperForm } from "sveltekit-superforms/client";
import type { AnyZodObject } from "zod";

type Message = {
  msg: string;
  type?: "info" | "success" | "warning" | "error";
};

export function superForm<T extends AnyZodObject>(
  ...params: Parameters<typeof defaultSuperForm<T>>
) {
  return defaultSuperForm<T, Message>(params[0], {
    delayMs: 100,
    taintedMessage: null,
    clearOnSubmit: "none",
    onResult: ({ result }) => {
      // superform doesn't handle CSRF errors properly
      // forward to actual error handler so we can test in dev
      const coercedResult = result as any;
      if (
        coercedResult.message ==
        "Cross-site POST form submissions are forbidden"
      ) {
        coercedResult.error = { ...coercedResult };
        result.type = "error";
        result.status = 403;
      }
    },
    onError: async ({ result, message }) => {
      if (
        result.status === 403 &&
        result.error.message ===
          "Cross-site POST form submissions are forbidden"
      ) {
        const origin = await (await fetch("/api/origin")).text();
        message.set({
          type: "error",
          msg: `For security reasons, please use ${origin} instead.`,
        });
      }
    },
    ...params[1],
  });
}
