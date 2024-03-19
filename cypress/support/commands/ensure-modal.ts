import { logCommand } from "./helpers";

const NO_LOG = { log: false };

type EnsureModalOptions = {
  log?: boolean;
};

export default function ensureModal(
  expectedTitle: string | RegExp,
  options?: EnsureModalOptions
): Cypress.Chainable<JQuery<HTMLElement>> {
  const log =
    options?.log !== false
      ? logCommand(
          "ensureModal",
          { "Expected title": expectedTitle },
          expectedTitle
        )
      : undefined;

  cy.get(".modal", NO_LOG).then((modalBackdrop) => {
    // Only assert with Chai if it is failing to not bloat the command log
    if (!modalBackdrop.get(0).classList.contains("show")) {
      // Assertion "be.visible" doesn't work here because it is behind the dialog
      cy.wrap(modalBackdrop, NO_LOG).should("have.class", "show");
    }
  });

  // Wait 300ms for the dialog slide animation to complete before continuing
  cy.wait(300, NO_LOG);

  return cy
    .get(".modal-dialog", NO_LOG)
    .should("be.visible")
    .within(NO_LOG, () => {
      if (expectedTitle instanceof RegExp) {
        cy.get(".modal-body", NO_LOG)
          .invoke(NO_LOG, "text")
          .should("match", expectedTitle);
      } else {
        cy.get(".modal-body", NO_LOG).should("have.text", expectedTitle);
      }
    })
    .finishLog(log);
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureModal(
        expectedTitle: string | RegExp,
        options?: EnsureModalOptions
      ): Chainable<JQuery<HTMLElement>>;
    }
  }
}
