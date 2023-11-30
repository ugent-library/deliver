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
        expect(response.isOkStatusCode).to.be.true
        expect(response.redirects).is.an('array').that.has.length(1)

        const redirects = response.redirects
          .map(url => url.replace(/^3\d\d\: /, '')) // Redirect entries are in form '3XX: {url}'
          .map(url => new URL(url))

        expect(redirects[0]).to.have.property('hostname', 'test.liblogin.ugent.be')
      })
    }

    cy.contains('header .btn', 'Log in').invoke('attr', 'href').then(assertLoginRedirection)

    cy.contains('main .btn', 'Log in').invoke('attr', 'href').then(assertLoginRedirection)
  })

  it('should be able to load the homepage as space admin', () => {
    cy.loginAsSpaceAdmin()

    cy.visit('/')

    cy.location('pathname').should('eq', '/spaces')

    cy.contains('header .btn', 'Log in').should('not.exist')
    cy.contains('main .btn', 'Log in').should('not.exist')

    cy.get('.c-sidebar').should('have.length', 1).should('have.class', 'c-sidebar--dark-gray')

    cy.get('.c-sub-sidebar').should('be.visible')
    cy.contains('Your deliver spaces').should('be.visible')
  })
})
