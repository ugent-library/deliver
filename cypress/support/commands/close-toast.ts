import { logCommand } from "./helpers";

const NO_LOG = { log: false };

type CloseToastOptions = {
  log?: boolean;
};

export default function closeToast(
  subject: JQuery<HTMLElement>,
  options?: CloseToastOptions
): Cypress.Chainable<JQuery<HTMLElement>> {
  if (options?.log !== false) {
    logCommand("closeToast", { subject, options });
  }

  if (!subject.is(".toast")) {
    throw new Error("Command subject is not a toast.");
  }

  return cy.wrap(subject, NO_LOG).within(NO_LOG, () => {
    cy.get(".btn-close", NO_LOG).click(NO_LOG);
  });
}

declare global {
  namespace Cypress {
    interface Chainable<Subject> {
      closeToast(options?: CloseToastOptions): Chainable<Subject>;
    }
  }
}
