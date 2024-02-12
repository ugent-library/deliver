describe("Clean up test folders and files", { redirectionLimit: 1000 }, () => {
  const SELECTOR = 'table.table tr td:first-of-type a:contains("CYPRESS-")';

  it("should clean up all files and folders", () => {
    cy.loginAsSpaceAdmin();

    cy.intercept("POST", "/folders/*").as("deleteFolder");

    deleteAllCypressTestFolders(Cypress.env("DEFAULT_SPACE"));

    cy.get(SELECTOR).should("not.exist");
  });

  function deleteAllCypressTestFolders(space: string) {
    cy.visit(`/spaces/${space}`).then(() => {
      // Using Cypress.$() direct jQuery selector tool here.
      // Using cy.get() the test would fail if none are left.
      const links = Cypress.$<HTMLAnchorElement>(SELECTOR)
        .map((_, link) => Cypress.$(link).attr("href"))
        .get();

      cy.log(`${links.length} test folder(s) found`);

      links.forEach(deleteTestFolder);
    });
  }

  function deleteTestFolder(href: string) {
    cy.visit(`${href}/edit`);

    cy.contains(".btn", "Delete folder").click();

    cy.wait("@deleteFolder");
  }
});
