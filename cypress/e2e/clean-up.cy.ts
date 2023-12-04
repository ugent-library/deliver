const DEFAULT_SPACE = 'test'

describe('Clean up test folders and files', () => {
  it('should clean up all files and folders', () => {
    cy.loginAsSpaceAdmin()

    // TODO Remove when issue #99 is resolved
    Cypress.on('uncaught:exception', (err, runnable) => {
      // returning false here prevents Cypress from failing the test
      return false
    })

    deleteAllCypressTestFolders(DEFAULT_SPACE)
  })

  function deleteAllCypressTestFolders(space: string) {
    cy.visit(`/spaces/${space}`).then(() => {
      // Using Cypress.$() direct jQuery selector tool.
      // Using cy.get() the test would fail if none are left.
      const links = Cypress.$<HTMLAnchorElement>('table.table tr td:first-of-type a:contains("CYPRESS-")')

      cy.log(`${links.length} test folder(s) found`)

      links.each(deleteTestFolder)
    })
  }

  function deleteTestFolder(_index: number, link: HTMLAnchorElement) {
    cy.visit(link.href)

    cy.contains('.btn', 'Edit').click()

    cy.contains('.btn', 'Delete folder').click()
  }
})
