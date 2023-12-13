import { logCommand } from './helpers'

export default function logout(): void {
  logCommand('logout')

  cy.clearAllCookies({ log: false })
}

declare global {
  namespace Cypress {
    interface Chainable {
      logout(): Chainable<void>
    }
  }
}
