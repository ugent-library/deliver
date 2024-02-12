describe("Folder searching", () => {
  const TEST_FOLDER_NAMES = [
    "Personal documents",
    "School work",
    "Work documents",
    "School projects",
    "Financial records",
  ] as const;

  type TestFolderNames = (typeof TEST_FOLDER_NAMES)[number];

  before(() => {
    cy.loginAsSpaceAdmin();

    // Make sure test folders exist
    TEST_FOLDER_NAMES.forEach((folderName) => {
      cy.visitSpace({
        qs: { q: folderName },
      });

      cy.getNumberOfDisplayedFolders().then((count) => {
        if (count === 0) {
          cy.setFieldByLabel("Folder name", "CYPRESS-" + folderName);
          cy.contains(".btn", "Make folder").click();
        }
      });
    });
  });

  beforeEach(() => {
    cy.loginAsSpaceAdmin();

    cy.visitSpace();

    cy.getNumberOfDisplayedFolders().should("be.at.least", 5);
    cy.getTotalNumberOfFolders()
      .should("be.at.least", 5)
      .as("totalBeforeFiltering");

    cy.intercept(`/spaces/${Cypress.env("DEFAULT_SPACE")}/folders*`).as(
      "filterFolders"
    );
  });

  it("should filter when clicking the search button", () => {
    cy.visitSpace();

    cy.get("input[name=q]").type("School", { delay: 0 });
    cy.contains(".btn", "Search").click();

    cy.get("@filterFolders").should("be.null");

    cy.url().should("have.param", "q", "School");

    assertFilteredFolders(["School work", "School projects"]);
  });

  it("should filter case insensitively", () => {
    cy.visitSpace();

    cy.get("input[name=q]").type("wORk", { delay: 0 });
    cy.contains(".btn", "Search").click();

    cy.url().should("have.param", "q", "wORk");

    assertFilteredFolders(["School work", "Work documents"]);
  });

  it("should filter when hitting the ENTER key in the search field", () => {
    cy.get("input[name=q]").type("documents{enter}", { delay: 0 });

    cy.get("@filterFolders").should("be.null");

    cy.url().should("have.param", "q", "documents");

    assertFilteredFolders(["Personal documents", "Work documents"]);
  });

  it("should filter when you blur the search field", () => {
    cy.get("input[name=q]").type("School", { delay: 0 }).blur();

    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("eql", { q: "School" });

    // TODO cy.url().should("have.param", "q", "School");

    assertFilteredFolders(["School work", "School projects"]);
  });

  it("should filter automatically upon typing", () => {
    cy.get("input[name=q]").type("records", { delay: 0 });

    cy.get("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("eql", { q: "records" });

    // TODO cy.url().should("have.param", "q", "records");

    assertFilteredFolders(["Financial records"]);
  });

  it("should debounce filtering when still typing", () => {
    cy.get("input[name=q]").as("q").type("Fin", { delay: 100 });
    cy.wait(300);
    cy.get("@filterFolders").should("be.null");

    cy.get("@q").type("anci", { delay: 100 });
    cy.wait(300);
    cy.get("@filterFolders").should("be.null");

    cy.get("@q").type("al rec", { delay: 100 });
    cy.wait(600);
    cy.get("@filterFolders", { timeout: 0 })
      .should("not.be.null")
      .should("have.nested.property", "request.query")
      .should("eql", { q: "Financial rec" });

    // TODO cy.url().should("have.param", "q", "Financial rec");

    assertFilteredFolders(["Financial records"]);
  });

  it("should retain the current AJAX search when reloading the page", () => {
    cy.get("input[name=q]").type("School", { delay: 0 }).blur();

    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("eql", { q: "School" });

    assertFilteredFolders(["School work", "School projects"]);

    cy.reload();

    assertFilteredFolders(["School work", "School projects"]);
  });

  describe("Folder results table", () => {
    before(() => {
      // Add 1 uploaded file to validate folders table
      cy.visitSpace({
        qs: { q: "School work" },
      });

      cy.getNumberOfDisplayedFolders().should("eq", 1);

      cy.contains("#folders a", "School work").click();

      cy.get("#folder-number-of-files")
        .invoke("text")
        .then((filesText) => {
          if (filesText === "0 items") {
            cy.intercept("POST", "/folders/*/files").as("uploadFile");

            cy.get("input[type=file]").selectFile("cypress/fixtures/test.pdf", {
              action: "select",
            });

            cy.wait("@uploadFile");
          }
        });
    });

    it("should load correctly during full page load", () => {
      cy.visitSpace({ qs: { q: "School work" } });

      cy.get("@filterFolders").should("be.null");

      assertFileDetailsLoadedCorrectly();
    });

    it("should load correctly during AJAX filtering", () => {
      cy.visitSpace();

      cy.get("input[name=q]").type("School work");

      cy.wait(600);

      cy.get("@filterFolders", { timeout: 0 }).should("not.be.null");

      assertFileDetailsLoadedCorrectly();
    });

    function assertFileDetailsLoadedCorrectly() {
      cy.getNumberOfDisplayedFolders().should("eq", 1);

      cy.get("#folders table tbody tr td")
        .eq(3)
        .should("contain", "1 file")
        .should("contain", "130 B");
    }
  });

  function assertFilteredFolders(filteredFolders: TestFolderNames[]) {
    cy.get("#folders").scrollIntoView();

    cy.getNumberOfDisplayedFolders().should("eq", filteredFolders.length);
    cy.getTotalNumberOfFolders().should("eq", "@totalBeforeFiltering");

    TEST_FOLDER_NAMES.forEach((folderName) => {
      const assertion = filteredFolders.includes(folderName)
        ? "be.visible"
        : "not.exist";

      cy.contains("#folders table td a", folderName).should(assertion);
    });
  }
});
