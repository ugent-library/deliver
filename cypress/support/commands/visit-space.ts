import { logCommand } from "./helpers";

export default function visitSpace(
  space?: string
): Cypress.Chainable<Cypress.AUTWindow> {
  space ||= Cypress.env("DEFAULT_SPACE");

  logCommand("visitSpace", { space }, space);

  return cy.visit(`/spaces/${space}`, { log: false });
}

declare global {
  namespace Cypress {
    interface Chainable {
      visitSpace(space?: string): Chainable<AUTWindow>;
    }
  }
}
