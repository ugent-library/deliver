describe("Folders sorting", () => {
  beforeEach(() => {
    cy.loginAsSpaceAdmin();

    cy.intercept(`/spaces/${Cypress.env("DEFAULT_SPACE")}/folders*`).as(
      "filterFolders"
    );
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

    // TODO: check actual sorting
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

    // TODO: check actual sorting
  });

  it("should reselect the chosen sort option", () => {
    cy.visitSpace({ qs: { sort: "expires-last" } });

    cy.get("select[name=sort]").should("have.value", "expires-last");
    cy.get('option[value="expires-last"]').should("have.attr", "selected");

    // TODO: check actual sorting
  });
});
