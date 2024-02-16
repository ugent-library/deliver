import { getRandomText } from "support/util";

describe("Folders sorting", () => {
  const randomSuffix = getRandomText();

  const TEST_FOLDER_NAMES = ["XYZ", "OPQ", "LMN", "FGH", "ABC"].map(
    (f) => `${f}-${randomSuffix}`
  );

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

  it("should have guaranteed order for sort field options", () => {
    cy.visitSpace();

    cy.get('select[name="sort"] option')
      .map("value")
      .should("eql", ["default", "expires-last"]);
  });

  it("should sort folders by expiration date asc by default", () => {
    cy.visitSpace({ qs: { q: randomSuffix } });

    cy.get("select[name=sort]").should("have.value", "default");

    cy.get("#folders table tbody tr td:first-of-type a")
      .map("textContent")
      .should("eql", TEST_FOLDER_NAMES);
  });

  it("should be possible to sort folders by expiration date asc", () => {
    cy.visitSpace({ qs: { q: randomSuffix, sort: "expires-last" } });

    cy.get("select[name=sort]").should("have.value", "expires-last");

    cy.setFieldByLabel("Sort by", "default");
    cy.wait("@filterFolders");

    cy.get("#folders table tbody tr td:first-of-type a")
      .map("textContent")
      .should("eql", TEST_FOLDER_NAMES);
  });

  it("should be possible to sort folders by expiration date desc", () => {
    cy.visitSpace({ qs: { q: randomSuffix } });

    cy.get("select[name=sort]").should("have.value", "default");

    cy.setFieldByLabel("Sort by", "expires-last");
    cy.wait("@filterFolders");

    cy.get("#folders table tbody tr td:first-of-type a")
      .map("textContent")
      .should("eql", TEST_FOLDER_NAMES.reverse());
  });

  it("should keep the sort choice when searching", () => {
    cy.visitSpace();

    cy.get("select[name=sort]").should("have.value", "default");

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

    cy.get("select[name=sort]").should("have.value", "default");

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

  it("should clear the sort param from the URL when default is selected", () => {
    cy.visitSpace({ qs: { q: "test", sort: "expires-last" } });

    cy.get("select[name=sort]").should("have.value", "expires-last");

    cy.setFieldByLabel("Sort by", "default");

    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("contain", { q: "test", sort: "default" });

    cy.url().should("not.have.param", "sort");
    cy.location("search").should("not.contain", "sort");
  });
});
