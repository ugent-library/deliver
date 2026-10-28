export default function loginAsSpaceAdmin(): void {
  cy.login(
    Cypress.expose("SPACE_ADMIN_USER_NAME"),
    Cypress.expose("SPACE_ADMIN_USER_PASSWORD"),
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      loginAsSpaceAdmin(): Chainable<void>;
    }
  }
}
