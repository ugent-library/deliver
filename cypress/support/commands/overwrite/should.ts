import { mapAlias } from "../../util";

Cypress.Commands.overwrite<"should", "optional">(
  "should",
  function (
    this: Record<string, unknown>,
    originalFn: CallableFunction,
    subject: unknown,
    chainer: unknown,
    ...args: unknown[]
  ) {
    args = args.map((a) => mapAlias(this, a));

    return originalFn(subject, chainer, ...args);
  },
);
