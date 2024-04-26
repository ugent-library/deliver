// https://github.com/ugent-library/deliver/issues/72

import { getRandomText } from "support/util";

describe("Issue #73: [Speed and usability] Add sort to folder overview on expiry date", () => {
  const randomSuffix = getRandomText();

  const TEST_FOLDER_NAMES = ["XYZ", "OPQ", "LMN", "FGH", "ABC"].map(
    (f) => `${f}-${randomSuffix}`,
  );

  before(() => {
    cy.loginAsSpaceAdmin();

    // Make sure test folders exist
    TEST_FOLDER_NAMES.forEach((folderName) => {
      cy.makeFolder(folderName, { noRedirect: true });

      // Make sure the creation (and expiration dates) are different for each folder
      cy.wait(1000);
    });
  });

  beforeEach(() => {
    cy.loginAsSpaceAdmin();

    cy.intercept(`/spaces/${Cypress.env("DEFAULT_SPACE")}/folders*`).as(
      "filterFolders",
    );

    cy.visitSpace();
    cy.get("select[name=sort]").as("sort");
  });

  it("should have guaranteed order for sort field options", () => {
    cy.visitSpace();

    cy.get("@sort")
      .find("option")
      .map("value")
      .should("eql", ["default", "expires-last"]);
  });

  it("should sort folders by expiration date asc by default", () => {
    cy.visitSpace({ qs: { q: randomSuffix } });
    cy.get("@sort").should("have.value", "default");

    cy.get("#folders table tbody tr td:first-of-type a")
      .map("textContent")
      .should("have.ordered.members", TEST_FOLDER_NAMES);
  });

  it("should be possible to sort folders by expiration date asc", () => {
    cy.visitSpace({ qs: { q: randomSuffix, sort: "expires-last" } });
    cy.get("@sort").should("have.value", "expires-last");

    cy.setFieldByLabel("Sort by", "default");
    cy.wait("@filterFolders");

    cy.get("#folders table tbody tr td:first-of-type a")
      .map("textContent")
      .should("have.ordered.members", TEST_FOLDER_NAMES);
  });

  it("should be possible to sort folders by expiration date desc", () => {
    cy.visitSpace({ qs: { q: randomSuffix } });
    cy.get("@sort").should("have.value", "default");

    cy.setFieldByLabel("Sort by", "expires-last");
    cy.wait("@filterFolders");

    cy.get("#folders table tbody tr td:first-of-type a")
      .map("textContent")
      .should("have.ordered.members", TEST_FOLDER_NAMES.reverse());
  });

  it("should keep the sort choice when searching", () => {
    cy.visitSpace();
    cy.get("@sort").should("have.value", "default");

    cy.setFieldByLabel("Sort by", "expires-last");
    cy.wait("@filterFolders");
    cy.getParams("sort").should("eq", "expires-last");
    cy.get("@sort").should("have.value", "expires-last");

    cy.contains(".btn", "Search").click();
    cy.get("@filterFolders.all").should(
      "have.length",
      1,
      "Search shouldn't have fired an AJAX request",
    );
    cy.getParams("sort").should("eq", "expires-last");
    cy.get("@sort").should("have.value", "expires-last");
  });

  it("should keep the sort choice when searching using AJAX", () => {
    cy.visitSpace();
    cy.get("@sort").should("have.value", "default");

    cy.setFieldByLabel("Sort by", "expires-last");
    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("contain", { sort: "expires-last" });
    cy.getParams("sort").should("eq", "expires-last");

    cy.get("input[name=q]").type("name").blur();
    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("contain", { sort: "expires-last" });

    cy.getParams("sort").should("eq", "expires-last");
    cy.get("@sort").should("have.value", "expires-last");
  });

  it("should reselect the chosen sort option", () => {
    cy.visitSpace({ qs: { sort: "expires-last" } });

    cy.get("@sort").should("have.value", "expires-last");
    cy.get('option[value="expires-last"]').should("have.attr", "selected");
  });

  it("should clear the sort param from the URL when default is selected", () => {
    cy.visitSpace({ qs: { q: "test", sort: "expires-last" } });

    cy.get("@sort").should("have.value", "expires-last");

    cy.setFieldByLabel("Sort by", "default");
    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("contain", { sort: "default" });

    cy.url().should("not.have.param", "sort");
    cy.location("search").should("not.contain", "sort");
  });
});
