import { logCommand } from "./helpers";

export default function visitSpace(
  space?: string,
  options?: Partial<Cypress.VisitOptions>,
): Cypress.Chainable<Cypress.AUTWindow> {
  if (typeof space !== "string") {
    options = space;
    space = Cypress.env("DEFAULT_SPACE");
  }

  if (options?.log !== false) {
    logCommand("visitSpace", { space, options }, space);
  }

  return cy.visit(`/spaces/${space}`, { ...options, log: false });
}

declare global {
  namespace Cypress {
    interface Chainable {
      visitSpace(options?: string): Chainable<AUTWindow>;

      visitSpace(
        space?: string,
        options?: Partial<VisitOptions>,
      ): Chainable<AUTWindow>;
      visitSpace(options?: Partial<VisitOptions>): Chainable<AUTWindow>;
    }
  }
}
