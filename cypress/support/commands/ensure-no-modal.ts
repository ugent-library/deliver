import { logCommand } from "./helpers";

type EnsureNoModalOptions = {
  log?: boolean;
};

export default function ensureNoModal(
  options: EnsureNoModalOptions = { log: true }
): void {
  if (options.log === true) {
    logCommand("ensureNoModal");
  }

  cy.get(".modal > *", { log: false })
    .should("have.length", 0)
    .then(() => {
      // Check before asserting to keep out of command log if ok
      if (Cypress.$(".modal, modal-dialog, .modal-content").length > 0) {
        cy.get(".modal").should("not.exist");
        cy.get(".modal-dialog").should("not.exist");
      }
    });
}

declare global {
  namespace Cypress {
    interface Chainable {
      ensureNoModal(options?: EnsureNoModalOptions): Chainable<void>;
    }
  }
}
