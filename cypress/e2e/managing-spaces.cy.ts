describe('Managing spaces', () => {
  beforeEach(() => {
    cy.loginAsSuperAdmin()
  })

  it('should return an error if a new space name is already in use', () => {
    const SPACE_NAME = Cypress.env('DEFAULT_SPACE')

    cy.visit('/')

    cy.contains('Make a new space').click()

    cy.location('pathname').should('eq', '/new-space')

    cy.get('#space-name').should('not.have.class', 'is-invalid')
    cy.get('#space-name-invalid').should('not.exist')

    cy.setFieldByLabel('Space name', SPACE_NAME)
    cy.setFieldByLabel(
      'Space admins',
      Cypress.env('SUPER_ADMIN_USER_NAME') + ',' + Cypress.env('SPACE_ADMIN_USER_NAME')
    )
    cy.contains('.btn', 'Make Space').click()

    cy.get('#space-name').should('have.class', 'is-invalid')
    cy.get('#space-name-invalid').should('be.visible').and('have.text', 'name must be unique')
  })
})
