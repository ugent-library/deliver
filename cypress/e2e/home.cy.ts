describe('The home page', () => {
  it('should be able to load the home page anonymously', () => {
    cy.visit('/')

    cy.location('pathname').should('eq', '/')

    cy.contains('header .btn', 'Log in').should('be.visible')
    cy.contains('main .btn', 'Log in').should('be.visible')

    cy.get('.c-sidebar').should('have.length', 1).should('not.have.class', 'c-sidebar--dark-gray')

    cy.get('.c-sub-sidebar').should('not.exist')
    cy.contains('Your deliver spaces').should('not.exist')
  })

  it('should redirect to the login page when clicking the Login buttons', () => {
    cy.visit('/')

    const assertLoginRedirection = href => {
      cy.request(href).then(response => {
        expect(response).to.have.property('isOkStatusCode', true)
        expect(response).to.have.property('redirects').that.is.an('array').that.has.length(1)

        const redirects = response.redirects
          .map(url => url.replace(/^3\d\d\: /, '')) // Redirect entries are in form '3XX: {url}'
          .map(url => new URL(url))

        expect(redirects[0]).to.have.property('hostname', 'test.liblogin.ugent.be')
      })
    }

    cy.get('header .btn:contains("Log in"), main .btn:contains("Log in")')
      .should('have.length', 2)
      .map('href')
      .unique() // No need to check the same URL more than once
      .each(assertLoginRedirection)
  })

  it('should be able to load the homepage as space admin', () => {
    cy.loginAsSpaceAdmin()

    cy.visit('/')

    cy.location('pathname').should('eq', '/spaces')

    cy.contains('header .btn', 'Log in').should('not.exist')
    cy.contains('main .btn', 'Log in').should('not.exist')

    cy.get('.c-sidebar').should('have.length', 1).should('have.class', 'c-sidebar--dark-gray')

    cy.wrap(Cypress.env('DEFAULT_SPACE')).should('not.be.empty')

    cy.get('.c-sub-sidebar').should('be.visible')
    cy.contains('.bc-navbar', 'Your deliver spaces')
      .should('be.visible')
      .next('.c-sub-sidebar__menu')
      .contains('a', Cypress.env('DEFAULT_SPACE'))
      .should('be.visible')
      .click()

    cy.location('pathname').should('eq', `/spaces/${Cypress.env('DEFAULT_SPACE')}`)
  })
})
