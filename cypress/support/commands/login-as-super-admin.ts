export default function loginAsSuperAdmin(): void {
  cy.login(
    Cypress.expose("SUPER_ADMIN_USER_NAME"),
    Cypress.expose("SUPER_ADMIN_USER_PASSWORD"),
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      loginAsSuperAdmin(): Chainable<void>;
    }
  }
}
