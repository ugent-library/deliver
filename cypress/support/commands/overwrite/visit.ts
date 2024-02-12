import { mapAlias } from "./util";

Cypress.Commands.overwrite<"visit">(
  "visit",
  function (
    this: Record<string, unknown>,
    originalFn: CallableFunction,
    subject: unknown,
    ...args: unknown[]
  ) {
    subject = mapAlias(this, subject);

    return originalFn(subject, ...args);
  }
);
