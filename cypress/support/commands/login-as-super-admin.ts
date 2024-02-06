export default function loginAsSuperAdmin(): void {
  cy.login(Cypress.env('SUPER_ADMIN_USER_NAME'), Cypress.env('SUPER_ADMIN_USER_PASSWORD'))
}

declare global {
  namespace Cypress {
    interface Chainable {
      loginAsSuperAdmin(): Chainable<void>
    }
  }
}
