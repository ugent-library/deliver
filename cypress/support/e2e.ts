import 'cypress-common'

import './commands'
import './commands/overwrite/should'
import './commands/overwrite/visit'
import './queries'

before(() => {
  const DEFAULT_SPACE = Cypress.env('DEFAULT_SPACE')

  function createDefaultSpace() {
    cy.contains('label', 'Space name').type(DEFAULT_SPACE)
    cy.contains('label', 'Space admins').type(Cypress.env('SPACE_ADMIN_USER_NAME'))
    cy.contains('.btn', 'Make Space').click()
  }

  cy.loginAsSuperAdmin()

  cy.visit('/')

  cy.location('pathname').then(pathname => {

    // When no spaces exist yet, the superadmin is enforced to create one
    if (pathname === '/new-space') {
      cy.ensureToast('Create an initial space to get started').closeToast()
      cy.ensureNoToast()

      createDefaultSpace()
    }

    // When some spaces already exist, we need to check if the DEFAULT_SPACE has to be created first 
    if (pathname === '/spaces') {
      cy.then(() => {
        const defaultSpace = Cypress.$(`.c-sub-sidebar .c-sidebar__label:contains("${DEFAULT_SPACE}")`)

        if (defaultSpace.length === 0) {
          cy.visit('/new-space')

          createDefaultSpace()
        }
      })
    }
  })

  cy.logout()
})
