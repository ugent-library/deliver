describe("Folders sorting", () => {
  const randomSuffix = Math.floor(Math.random() * 1000000).toLocaleString(
    undefined,
    { minimumIntegerDigits: 6, useGrouping: false }
  );

  const TEST_FOLDER_NAMES = [
    `CYPRESS-XYZ-${randomSuffix}`,
    `CYPRESS-OPQ-${randomSuffix}`,
    `CYPRESS-LMN-${randomSuffix}`,
    `CYPRESS-FGH-${randomSuffix}`,
    `CYPRESS-ABC-${randomSuffix}`,
  ];

  before(() => {
    cy.loginAsSpaceAdmin();

    cy.visitSpace();

    // Make sure test folders exist
    TEST_FOLDER_NAMES.forEach((folderName) => {
      cy.setFieldByLabel("Folder name", folderName);
      cy.contains(".btn", "Make folder").click();

      cy.visitSpace();

      // Make sure the creation (and expiration dates) are different for each folder
      cy.wait(1000);
    });
  });

  beforeEach(() => {
    cy.loginAsSpaceAdmin();

    cy.intercept(`/spaces/${Cypress.env("DEFAULT_SPACE")}/folders*`).as(
      "filterFolders"
    );
  });

  it("should sort folders by expiration date asc by default", () => {
    cy.visitSpace({ qs: { q: randomSuffix } });

    cy.get("select[name=sort]").should("have.value", "expires-first");

    cy.get("#folders table tbody tr td:first-of-type a")
      .map("textContent")
      .should("eql", TEST_FOLDER_NAMES);
  });

  it("should be possible to sort folders by expiration date asc", () => {
    cy.visitSpace({ qs: { q: randomSuffix, sort: "expires-last" } });

    cy.get("select[name=sort]").should("have.value", "expires-last");

    cy.setFieldByLabel("Sort by", "expires-first");
    cy.wait("@filterFolders");

    cy.get("#folders table tbody tr td:first-of-type a")
      .map("textContent")
      .should("eql", TEST_FOLDER_NAMES);
  });

  it("should be possible to sort folders by expiration date desc", () => {
    cy.visitSpace({ qs: { q: randomSuffix } });

    cy.get("select[name=sort]").should("have.value", "expires-first");

    cy.setFieldByLabel("Sort by", "expires-last");
    cy.wait("@filterFolders");

    cy.get("#folders table tbody tr td:first-of-type a")
      .map("textContent")
      .should("eql", TEST_FOLDER_NAMES.reverse());
  });

  it("should keep the sort choice when searching", () => {
    cy.visitSpace();

    cy.get("select[name=sort]").should("have.value", "expires-first");

    cy.setFieldByLabel("Sort by", "expires-last");

    cy.wait("@filterFolders");

    cy.contains(".btn", "Search").click();

    cy.get("@filterFolders.all").should(
      "have.length",
      1,
      "Search shouldn't have fired an AJAX request"
    );

    cy.url().should("have.param", "sort", "expires-last");
    cy.get("select[name=sort]").should("have.value", "expires-last");
  });

  it("should keep the sort choice when searching using AJAX", () => {
    cy.visitSpace();

    cy.get("select[name=sort]").should("have.value", "expires-first");

    cy.setFieldByLabel("Sort by", "expires-last");

    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("eql", { q: "", sort: "expires-last" });

    cy.get("input[name=q]").type("name").blur();

    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("eql", { q: "name", sort: "expires-last" });

    cy.url().should("have.param", "sort", "expires-last");
    cy.get("select[name=sort]").should("have.value", "expires-last");
  });

  it("should reselect the chosen sort option", () => {
    cy.visitSpace({ qs: { sort: "expires-last" } });

    cy.get("select[name=sort]").should("have.value", "expires-last");
    cy.get('option[value="expires-last"]').should("have.attr", "selected");
  });
});
