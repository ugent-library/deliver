import { getRandomText } from "support/util";

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
          cy.makeFolder(folderName);
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

    cy.get("input[name=q]").as("q");
  });

  it("should filter when clicking the search button", () => {
    cy.visitSpace();

    cy.get("@q").type("School", { delay: 0 });
    cy.contains(".btn", "Search").click();

    cy.url().should("have.param", "q", "School");
    cy.get("@q").should("have.value", "School");

    assertFilteredFolders(["School work", "School projects"]);
  });

  it("should filter case insensitively", () => {
    cy.visitSpace();

    cy.get("@q").type("wORk", { delay: 0 });
    cy.contains(".btn", "Search").click();

    cy.url().should("have.param", "q", "wORk");
    cy.get("@q").should("have.value", "wORk");

    assertFilteredFolders(["School work", "Work documents"]);
  });

  it("should filter when hitting the ENTER key in the search field", () => {
    cy.get("@q").type("documents{enter}", { delay: 0 });

    cy.url().should("have.param", "q", "documents");
    cy.get("@q").should("have.value", "documents");

    assertFilteredFolders(["Personal documents", "Work documents"]);
  });

  it("should filter when you blur the search field", () => {
    cy.get("@q").type("School", { delay: 0 }).blur();

    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("contain", { q: "School" });

    cy.url().should("have.param", "q", "School");
    cy.get("@q").should("have.value", "School");

    assertFilteredFolders(["School work", "School projects"]);
  });

  it("should filter automatically upon typing", () => {
    cy.get("@q").type("records", { delay: 0 });

    cy.get("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("contain", { q: "records" });

    cy.url().should("have.param", "q", "records");

    assertFilteredFolders(["Financial records"]);
  });

  it("should debounce filtering when still typing", () => {
    cy.get("@q").type("Fin", { delay: 100 });
    cy.wait(300);
    cy.get("@filterFolders").should("be.null");

    cy.get("@q").type("amci", { delay: 100 });
    cy.wait(300);
    cy.get("@filterFolders").should("be.null");

    // Oops typo!
    cy.get("@q").type("{backspace}{backspace}{backspace}", { delay: 100 });
    cy.wait(300);
    cy.get("@filterFolders").should("be.null");

    cy.get("@q").type("nci", { delay: 100 });
    cy.wait(300);
    cy.get("@filterFolders").should("be.null");

    cy.get("@q").type("al rec", { delay: 100 });
    cy.wait(600);
    cy.get("@filterFolders", { timeout: 0 })
      .should("not.be.null")
      .should("have.nested.property", "request.query")
      .should("contain", { q: "Financial rec" });

    cy.url().should("have.param", "q", "Financial rec");
    cy.get("@q").should("have.value", "Financial rec");

    assertFilteredFolders(["Financial records"]);
  });

  it("should retain the current AJAX search when reloading the page", () => {
    cy.getNumberOfDisplayedFolders().as("folderCountBeforeSearch");

    cy.get("@q").type("School", { delay: 0 }).blur();

    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("contain", { q: "School" });

    cy.url().should("have.param", "q", "School");

    assertFilteredFolders(["School work", "School projects"]);

    cy.reload();

    assertFilteredFolders(["School work", "School projects"]);

    cy.get("@q").clear();

    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("contain", { q: "" });

    cy.url().should("not.have.param", "q");
    cy.location("search").should("not.contain", "q=");
    cy.get("@q").should("have.value", "");

    cy.get("@folderCountBeforeSearch").should("eq", "@folderCountBeforeSearch");
  });

  it("should be possible to clear the search field", () => {
    cy.visitSpace({ qs: { q: "School work" } });

    cy.get("@q").should("have.value", "School work");

    // Clear by setting the value property instead of fake typing with cy.clear()
    cy.get<HTMLInputElement>("@q").then((q) => (q.get(0).value = ""));
    // Trigger the search event, which is equivalent to clicking the clear "✖️" button in the search input
    cy.get("@q").trigger("search");

    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("contain", { q: "" });

    cy.url().should("not.have.param", "q");
    cy.location("search").should("not.contain", "q=");
    cy.get("@q").should("have.value", "");
  });

  it("should be protected against SQL injection", () => {
    cy.get("@q").type("' OR 1=1 --").blur();

    cy.wait("@filterFolders")
      .should("have.nested.property", "request.query")
      .should("contain", { q: "' OR 1=1 --" });

    cy.getNumberOfDisplayedFolders().should("eq", 0);
  });

  it("should display a default message if nothing was filtered", () => {
    cy.get("@q").type(getRandomText());
    cy.contains(".btn", "Search").click();

    cy.get("#folders .c-blank-slate")
      .should("be.visible")
      .and("contain.text", "No folders to display.")
      .and("contain.text", "Refine your search or add a new folder.");
  });

  describe("Folder results table", () => {
    before(() => {
      cy.loginAsSpaceAdmin();

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

      cy.get("@q").type("School work");

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

    cy.get("#folders .c-blank-slate").should("not.exist");

    cy.getNumberOfDisplayedFolders().should("eq", filteredFolders.length);
    cy.getTotalNumberOfFolders().should("eq", "@totalBeforeFiltering");
    cy.get("#folders tbody tr").should("have.length", filteredFolders.length);

    TEST_FOLDER_NAMES.forEach((folderName) => {
      const assertion = filteredFolders.includes(folderName)
        ? "be.visible"
        : "not.exist";

      cy.contains("#folders table td a", folderName).should(assertion);
    });
  }
});
