import { logCommand } from "./helpers";

type EnsureToastOptions = {
  log?: boolean;
};

export default function ensureToast(
  expectedTitle?: string | RegExp,
  options?: EnsureToastOptions,
): Cypress.Chainable<JQuery<HTMLElement>> {
  if (typeof expectedTitle === "object" && !(expectedTitle instanceof RegExp)) {
    if (options) {
      throw new Error("Invalid arguments provided for command useToast.");
    }

    // Only options were provided
    options = expectedTitle;
    expectedTitle = undefined;
  }

  const log =
    options?.log !== false
      ? logCommand(
          "ensureToast",
          { "Expected title": expectedTitle, options },
          expectedTitle,
        )
      : undefined;

  return cy
    .get(".toast", { log: false })
    .should((toast) => {
      if (toast.length !== 1) {
        expect(toast).to.have.length(1);
      }

      if (expectedTitle) {
        if (typeof expectedTitle === "string") {
          expect(toast).to.contain(expectedTitle);
        } else {
          expect(toast).to.match(expectedTitle);
        }
      }

      return toast;
    })
    .should("be.visible")
    .finishLog(log);
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureToast(
        expectedTitle?: string | RegExp,
        options?: EnsureToastOptions,
      ): Chainable<JQuery<HTMLElement>>;
    }
  }
}
