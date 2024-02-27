describe("The home page", () => {
  it("should be able to load the home page anonymously", () => {
    cy.visit("/");

    cy.location("pathname").should("eq", "/");

    cy.contains("header .btn", "Log in").should("be.visible");
    cy.contains("main .btn", "Log in").should("be.visible");

    cy.get(".c-sidebar")
      .should("have.length", 1)
      .should("not.have.class", "c-sidebar--dark-gray");

    cy.get(".c-sub-sidebar").should("not.exist");
    cy.contains("Your deliver spaces").should("not.exist");
  });

  it("should redirect to the login page when clicking the Login buttons", () => {
    cy.visit("/");

    const assertLoginRedirection = (href: string) => {
      cy.request(href).then((response) => {
        expect(response).to.have.property("isOkStatusCode", true);
        expect(response).to.have.property("redirects").that.is.an("array").that
          .is.not.empty;

        const redirect = new URL(
          response.redirects!.at(-1)!.replace(/^3\d\d: /, "")
        ); // Redirect entries are in form '3XX: {url}'

        expect(redirect).to.have.property("origin", Cypress.env("OIDC_ORIGIN"));
      });
    };

    cy.get('header .btn:contains("Log in"), main .btn:contains("Log in")')
      .should("have.length", 2)
      .map("href")
      .unique() // No need to check the same URL more than once
      .each(assertLoginRedirection);
  });

  it("should be able to load the first space as space admin", () => {
    cy.loginAsSpaceAdmin();

    cy.request("/").then((response) => {
      const spacesUrl = new URL("/spaces", Cypress.config("baseUrl")!);
      expect(response.redirects!.at(-2)).to.eq(`303: ${spacesUrl}`);
      expect(response.redirects!.at(-1)).to.match(
        new RegExp(`^303: ${spacesUrl}/[\\w\\d]+$`)
      );
    });

    cy.visit("/");

    cy.location("pathname").should("match", /^\/spaces\/[\w\d]+$/);

    cy.contains("header .btn", "Log in").should("not.exist");
    cy.contains("main .btn", "Log in").should("not.exist");

    cy.get(".c-sidebar")
      .should("have.length", 1)
      .should("have.class", "c-sidebar--dark-gray");

    const DEFAULT_SPACE = Cypress.env("DEFAULT_SPACE");

    cy.wrap(DEFAULT_SPACE).should(
      "not.be.empty",
      "A default space has not been configured."
    );

    cy.get(".c-sub-sidebar").should("be.visible");
    cy.contains(".bc-navbar", "Your deliver spaces")
      .should("be.visible")
      .next(".c-sub-sidebar__menu")
      .contains("a", DEFAULT_SPACE)
      .should("be.visible")
      .click();

    cy.location("pathname").should("eq", `/spaces/${DEFAULT_SPACE}`);
  });
});
