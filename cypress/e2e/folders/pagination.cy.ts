// https://github.com/ugent-library/deliver/issues/91

import { getRandomText } from "support/util";

describe("Issue #91: [Speed and usability] Add pagination to folder overview", () => {
  // All following tests assume there are exactly 21 folders in this space
  const NUMBER_OF_TEST_FOLDERS = 21;

  before(() => {
    cy.loginAsSpaceAdmin();

    cy.visitSpace({ qs: { limit: NUMBER_OF_TEST_FOLDERS } });
    cy.getFolderCount("total").then((total) => {
      if (total === NUMBER_OF_TEST_FOLDERS) {
        const folderNames = Cypress.$(
          "#folders table tbody tr td:first-of-type a",
        )
          .get()
          .map((a) => a.textContent!);

        if (
          folderNames.filter((name) => name.includes("CYPRESS")).length === 7
        ) {
          // No need to recreate the test folders
          return;
        }
      }

      cy.cleanUp();

      cy.visitSpace();
      cy.getFolderCount("total").should("eq", 0);

      Cypress._.shuffle([
        ...Cypress._.times(7, () => getRandomText("CYPRESS")),
        ...Cypress._.times(NUMBER_OF_TEST_FOLDERS - 7, () => getRandomText()),
      ]).forEach((name) => cy.makeFolder(name, { noRedirect: true }));
    });

    cy.visitSpace();
    cy.getFolderCount("total").should("eq", NUMBER_OF_TEST_FOLDERS);
  });

  beforeEach(() => {
    cy.loginAsSpaceAdmin();

    cy.visitSpace();

    cy.get('input[name="q"]').as("q");
    cy.get('select[name="sort"]').as("sort");

    cy.get('.pagination .page-item:has(.page-link[aria-label="Previous"])').as(
      "previous",
    );
    cy.get('.pagination .page-item:has(.page-link[aria-label="Next"])').as(
      "next",
    );

    cy.intercept(`/spaces/${Cypress.env("DEFAULT_SPACE")}/folders*`).as(
      "getFolders",
    );
  });

  it("should display maximum 20 folders per page by default", () => {
    cy.visitSpace();

    cy.getFolderCount("text").should("eq", "Showing 1-20 of 21 folder(s)");
    getNumberOfPages().should("eq", 2);

    cy.contains(
      "#folders .card-header .pagination .page-item a.page-link",
      "2",
    ).should("be.visible");

    cy.get("#folders table tbody tr").should("have.length", 20);

    cy.get(".u-scroll-wrapper__body").scrollTo("bottom");

    cy.contains(
      "#folders .card-footer .pagination .page-item a.page-link",
      "2",
    ).should("be.visible");
  });

  it("should display the start and end result count on the current page", () => {
    cy.visitSpace({ qs: { limit: 10 } });
    cy.getFolderCount("text").should("eq", "Showing 1-10 of 21 folder(s)");

    cy.visitSpace({ qs: { limit: 10, offset: 10 } });
    cy.getFolderCount("text").should("eq", "Showing 11-20 of 21 folder(s)");

    cy.visitSpace({ qs: { limit: 10, offset: 20 } });
    cy.getFolderCount("text").should("eq", "Showing 21-21 of 21 folder(s)");
  });

  it("should remove offset from the URL when default (0)", () => {
    cy.visitSpace({ qs: { offset: 20 } });
    cy.getParams("offset").should("eq", "20");

    goToPage(1);
    cy.getParams("offset").should("not.exist");
  });

  it("should remove limit from the URL when default (20)", () => {
    cy.visitSpace({ qs: { limit: 20 } });
    cy.getParams("limit").should("eq", "20");

    goToPage(2);
    cy.getParams("limit").should("not.exist");
  });

  it("should display the filtered folder count when search query is used", () => {
    cy.visitSpace();

    cy.get("@q").type("4");
    cy.contains(".btn", "Search").click();

    cy.getFolderCount("total").should("be.lessThan", 21);
    getNumberOfPages().should("eq", 1);
  });

  it("should highlight the current page button", () => {
    cy.visitSpace({ qs: { limit: 6 } });

    cy.getFolderCount("text").should("eq", "Showing 1-6 of 21 folder(s)");

    [1, 4, 3, 2].forEach((page) => {
      goToPage(page);

      cy.get(`#folders .pagination .page-item:contains(${page})`)
        .should("have.length", 2)
        .should("have.class", "active");

      cy.get(`#folders .pagination .page-item:not(:contains(${page}))`).should(
        "not.have.class",
        "active",
      );
    });
  });

  it("should always display page 1 pagination button, even without results", () => {
    cy.visitSpace();
    cy.getFolderCount("text").should("eq", "Showing 1-20 of 21 folder(s)");

    cy.get(".card-header .pagination .page-item").should("have.length", 4);
    cy.get(".card-footer .pagination .page-item").should("have.length", 4);

    cy.get("@q").type("grmbl", { delay: 0 });
    cy.wait("@getFolders");
    cy.getFolderCount("text").should("eq", "Showing 0 folder(s)");

    getNumberOfPages().should("eq", 1);

    cy.get("@previous").should("have.class", "disabled");
    cy.get("@next").should("have.class", "disabled");

    cy.get(".card-header .pagination .page-item")
      .should("have.length", 3)
      .eq(1)
      .should("have.text", 1)
      .should("have.class", "active");

    // Footer pagination is not rendered with no results
    cy.get(".card-footer .pagination").should("not.exist");
  });

  it("should be possible to jump to the next and previous page", () => {
    cy.visitSpace({ qs: { limit: 5 } });
    cy.getFolderCount("text").should("eq", "Showing 1-5 of 21 folder(s)");

    goToPage("next");
    cy.getFolderCount("text").should("eq", "Showing 6-10 of 21 folder(s)");
    cy.getActivePage().should("eq", 2);
    cy.getParams().should("eql", { limit: "5", offset: "5" });

    goToPage("next");
    cy.getFolderCount("text").should("eq", "Showing 11-15 of 21 folder(s)");
    cy.getActivePage().should("eq", 3);
    cy.getParams().should("eql", { limit: "5", offset: "10" });

    goToPage("previous");
    cy.getFolderCount("text").should("eq", "Showing 6-10 of 21 folder(s)");
    cy.getActivePage().should("eq", 2);
    cy.getParams().should("eql", { limit: "5", offset: "5" });

    goToPage("previous");
    cy.getFolderCount("text").should("eq", "Showing 1-5 of 21 folder(s)");
    cy.getActivePage().should("eq", 1);
    cy.getParams().should("eql", { limit: "5" });
  });

  it("should disable the previous page button when on the first page", () => {
    cy.visitSpace({ qs: { limit: 5 } });
    cy.getActivePage().should("eq", 1);
    cy.get("@previous").should("have.class", "disabled");

    // Make sure nothing happens if we click it
    cy.get("@previous").first().click();
    cy.get("@previous").last().click();
    cy.get("@getFolders").should("be.null");

    // Test again after paging
    goToPage("next");
    cy.getActivePage().should("eq", 2);
    cy.get("@previous").should("not.have.class", "disabled");

    goToPage("previous");
    cy.getActivePage().should("eq", 1);
    cy.get("@previous").should("have.class", "disabled");
  });

  it("should disable the next page button when on the last page", () => {
    cy.visitSpace({ qs: { limit: 5, offset: 20 } });
    cy.getActivePage().should("eq", 5);
    cy.get("@next").should("have.class", "disabled");

    // Make sure nothing happens if we click it
    cy.get("@next").first().click();
    cy.get("@next").last().click();
    cy.get("@getFolders").should("be.null");

    // Test again after paging
    goToPage("previous");
    cy.getActivePage().should("eq", 4);
    cy.get("@next").should("not.have.class", "disabled");

    goToPage("next");
    cy.getActivePage().should("eq", 5);
    cy.get("@next").should("have.class", "disabled");
  });

  it("should disable previous and next page buttons when all results fit on one page (without search query)", () => {
    cy.visitSpace({ qs: { limit: 21 } });
    cy.getActivePage().should("eq", 1);
    getNumberOfPages().should("eq", 1);

    cy.get("@previous").should("have.class", "disabled");
    cy.get("@next").should("have.class", "disabled");
  });

  it("should disable previous and next page buttons when all results fit on one page (with search query)", () => {
    cy.visitSpace();
    cy.getActivePage().should("eq", 1);
    getNumberOfPages().should("eq", 2);
    cy.get("@previous").should("have.class", "disabled");
    cy.get("@next").should("not.have.class", "disabled");

    cy.get("@q").type("4");
    cy.wait("@getFolders");
    cy.getActivePage().should("eq", 1);
    getNumberOfPages().should("eq", 1);
    cy.get("@previous").should("have.class", "disabled");
    cy.get("@next").should("have.class", "disabled");
  });

  it("should display an ellipsis section in the pager when there are a lot of pages", () => {
    cy.visitSpace({ qs: { limit: 2 } });
    getNumberOfPages().should("eq", 11);
    cy.getActivePage().should("eq", 1);
    getPagerButtons().should("eql", ["<", 1, 2, "...", 11, ">"]);

    goToPage("next");
    cy.getActivePage().should("eq", 2);
    getPagerButtons().should("eql", ["<", 1, 2, 3, "...", 11, ">"]);

    goToPage("next");
    cy.getActivePage().should("eq", 3);
    getPagerButtons().should("eql", ["<", 1, 2, 3, 4, "...", 11, ">"]);

    goToPage("next");
    cy.getActivePage().should("eq", 4);
    getPagerButtons().should("eql", ["<", 1, 2, 3, 4, 5, "...", 11, ">"]);

    goToPage("next");
    cy.getActivePage().should("eq", 5);
    getPagerButtons().should("eql", ["<", 1, "...", 4, 5, 6, "...", 11, ">"]);

    goToPage("next");
    cy.getActivePage().should("eq", 6);
    getPagerButtons().should("eql", ["<", 1, "...", 5, 6, 7, "...", 11, ">"]);

    goToPage("next");
    cy.getActivePage().should("eq", 7);
    getPagerButtons().should("eql", ["<", 1, "...", 6, 7, 8, "...", 11, ">"]);

    goToPage("next");
    cy.getActivePage().should("eq", 8);
    getPagerButtons().should("eql", ["<", 1, "...", 7, 8, 9, 10, 11, ">"]);

    goToPage("next");
    cy.getActivePage().should("eq", 9);
    getPagerButtons().should("eql", ["<", 1, "...", 8, 9, 10, 11, ">"]);

    goToPage("next");
    cy.getActivePage().should("eq", 10);
    getPagerButtons().should("eql", ["<", 1, "...", 9, 10, 11, ">"]);

    goToPage("next");
    cy.getActivePage().should("eq", 11);
    getPagerButtons().should("eql", ["<", 1, "...", 10, 11, ">"]);
  });

  it("should not be possible to click ellipsis buttons", () => {
    cy.visitSpace({ qs: { limit: 2, offset: 10 } });
    getNumberOfPages().should("eq", 11);
    cy.getActivePage().should("eq", 6);

    cy.get(".pagination .page-item:has(.if-more)")
      .should("have.length", 4) // 2 in top pager + 2 in bottom pager
      .should("have.class", "disabled")
      .find("a")
      .should("not.exist");
  });

  it("should keep search query when switching pages", () => {
    cy.visitSpace({ qs: { limit: 3 } });
    cy.getParams("q").should("not.exist");
    cy.getFolderCount("text").should("eq", "Showing 1-3 of 21 folder(s)");

    cy.get("@q").type("CYPRESS");
    cy.wait("@getFolders");
    cy.getFolderCount("text").should("eq", "Showing 1-3 of 7 folder(s)");

    goToPage(3);
    cy.getParams("q").should("eq", "CYPRESS");
    cy.get("@q").should("have.value", "CYPRESS");
    cy.getFolderCount("text").should("eq", "Showing 7-7 of 7 folder(s)");

    goToPage("previous");
    cy.getParams("q").should("eq", "CYPRESS");
    cy.get("@q").should("have.value", "CYPRESS");
    cy.getFolderCount("text").should("eq", "Showing 4-6 of 7 folder(s)");

    goToPage("next");
    cy.getParams("q").should("eq", "CYPRESS");
    cy.get("@q").should("have.value", "CYPRESS");
    cy.getFolderCount("text").should("eq", "Showing 7-7 of 7 folder(s)");
  });

  it("should keep sorting when switching pages", () => {
    cy.visitSpace({ qs: { limit: 3 } });
    cy.getParams("sort").should("not.exist");

    cy.setFieldByLabel("Sort by", "expires-last");
    cy.wait("@getFolders");

    cy.getFolderCount("total").should("eq", 21);

    goToPage(2);
    cy.getParams("sort").should("eq", "expires-last");
    cy.get("@sort").should("have.value", "expires-last");
    cy.getFolderCount("total").should("eq", 21);

    goToPage(3);
    cy.getParams("sort").should("eq", "expires-last");
    cy.get("@sort").should("have.value", "expires-last");
    cy.getFolderCount("total").should("eq", 21);

    goToPage("previous");
    cy.getParams("sort").should("eq", "expires-last");
    cy.get("@sort").should("have.value", "expires-last");
    cy.getFolderCount("total").should("eq", 21);

    goToPage("next");
    cy.getParams("sort").should("eq", "expires-last");
    cy.get("@sort").should("have.value", "expires-last");
    cy.getFolderCount("total").should("eq", 21);
  });

  it("should keep the page limit when switching pages", () => {
    cy.visitSpace({ qs: { limit: 6 } });
    cy.getParams("limit").should("eq", "6");
    cy.getFolderCount("text").should("eq", "Showing 1-6 of 21 folder(s)");

    goToPage(2);
    cy.getParams("limit").should("eq", "6");
    cy.getFolderCount("text").should("eq", "Showing 7-12 of 21 folder(s)");

    goToPage(3);
    cy.getParams("limit").should("eq", "6");
    cy.getFolderCount("text").should("eq", "Showing 13-18 of 21 folder(s)");

    goToPage("previous");
    cy.getParams("limit").should("eq", "6");
    cy.getFolderCount("text").should("eq", "Showing 7-12 of 21 folder(s)");

    goToPage("next");
    cy.getParams("limit").should("eq", "6");
    cy.getFolderCount("text").should("eq", "Showing 13-18 of 21 folder(s)");
  });

  it("should reset to first page when changing the search query", () => {
    cy.visitSpace({ qs: { offset: 20 } });
    cy.getFolderCount("text").should("eq", "Showing 21-21 of 21 folder(s)");
    cy.getActivePage().should("eq", 2);
    getPagerButtons().should("eql", ["<", 1, 2, ">"]);

    cy.get("@q").type("cypress");
    cy.wait("@getFolders");
    cy.getParams("offset").should("not.exist");
    cy.getFolderCount("text").should("eq", "Showing 1-7 of 7 folder(s)");
    cy.getActivePage().should("eq", 1);
    getPagerButtons().should("eql", ["<", 1, ">"]);
  });

  it("should reset to first page when clearing the search query", () => {
    cy.visitSpace({ qs: { q: "cypress", limit: 3, offset: 3 } });
    cy.getFolderCount("text").should("eq", "Showing 4-6 of 7 folder(s)");
    cy.getActivePage().should("eq", 2);
    getPagerButtons().should("eql", ["<", 1, 2, 3, ">"]);

    cy.get("@q").clear();
    cy.wait("@getFolders");
    cy.getParams("offset").should("not.exist");
    cy.getFolderCount("text").should("eq", "Showing 1-3 of 21 folder(s)");
    cy.getActivePage().should("eq", 1);
    getPagerButtons().should("eql", ["<", 1, 2, "...", 7, ">"]);
  });

  it("should reset to first page when sorting the results", () => {
    cy.visitSpace();
    cy.getFolderCount("text").should("eq", "Showing 1-20 of 21 folder(s)");
    cy.getActivePage().should("eq", 1);

    goToPage(2);
    cy.getParams("offset").should("eq", "20");
    cy.getParams("sort").should("not.exist");
    cy.getActivePage().should("eq", 2);
    cy.getFolderCount("text").should("eq", "Showing 21-21 of 21 folder(s)");

    cy.setFieldByLabel("Sort by", "expires-last");
    cy.wait("@getFolders");
    cy.getParams("offset").should("not.exist");
    cy.getActivePage().should("eq", 1);
    cy.getFolderCount("text").should("eq", "Showing 1-20 of 21 folder(s)");

    goToPage(2);
    cy.getParams("offset").should("eq", "20");
    cy.getParams("sort").should("eq", "expires-last");
    cy.getActivePage().should("eq", 2);
    cy.getFolderCount("text").should("eq", "Showing 21-21 of 21 folder(s)");

    cy.setFieldByLabel("Sort by", "default");
    cy.wait("@getFolders");
    cy.getParams("offset").should("not.exist");
    cy.getParams("sort").should("not.exist");
    cy.getActivePage().should("eq", 1);
    cy.getFolderCount("text").should("eq", "Showing 1-20 of 21 folder(s)");
  });

  it("should work with HTMX disabled", () => {
    cy.disableHtmx();

    cy.visitSpace({ qs: { limit: 2 } });
    getNumberOfPages().should("eq", 11);
    cy.getActivePage().should("eq", 1);
    getPagerButtons().should("eql", ["<", 1, 2, "...", 11, ">"]);

    goToPage(2);
    cy.get("@getFolders").should("be.null");
    cy.getParams().should("eql", { limit: "2", offset: "2" });
    cy.getActivePage().should("eq", 2);
    getPagerButtons().should("eql", ["<", 1, 2, 3, "...", 11, ">"]);

    goToPage("next");
    cy.get("@getFolders").should("be.null");
    cy.getParams().should("eql", { limit: "2", offset: "4" });
    cy.getActivePage().should("eq", 3);
    getPagerButtons().should("eql", ["<", 1, 2, 3, 4, "...", 11, ">"]);

    goToPage(4);
    cy.get("@getFolders").should("be.null");
    cy.getParams().should("eql", { limit: "2", offset: "6" });
    cy.getActivePage().should("eq", 4);
    getPagerButtons().should("eql", ["<", 1, 2, 3, 4, 5, "...", 11, ">"]);

    goToPage(5);
    cy.get("@getFolders").should("be.null");
    cy.getParams().should("eql", { limit: "2", offset: "8" });
    cy.getActivePage().should("eq", 5);
    getPagerButtons().should("eql", ["<", 1, "...", 4, 5, 6, "...", 11, ">"]);

    goToPage("previous");
    cy.get("@getFolders").should("be.null");
    cy.getParams().should("eql", { limit: "2", offset: "6" });
    cy.getActivePage().should("eq", 4);
    getPagerButtons().should("eql", ["<", 1, 2, 3, 4, 5, "...", 11, ">"]);

    goToPage(11);
    cy.get("@getFolders").should("be.null");
    cy.getParams().should("eql", { limit: "2", offset: "20" });
    cy.getActivePage().should("eq", 11);
    getPagerButtons().should("eql", ["<", 1, "...", 10, 11, ">"]);

    goToPage("previous");
    cy.get("@getFolders").should("be.null");
    cy.getParams().should("eql", { limit: "2", offset: "18" });
    cy.getActivePage().should("eq", 10);
    getPagerButtons().should("eql", ["<", 1, "...", 9, 10, 11, ">"]);
  });

  it("should not crash with zero or negative limit", () => {
    cy.visitSpace({ qs: { limit: 0 } });
    cy.getFolderCount("text").should("eq", "Showing 1-20 of 21 folder(s)");

    cy.visitSpace({ qs: { limit: -1 } });
    cy.getFolderCount("text").should("eq", "Showing 1-20 of 21 folder(s)");
  });

  it("should not crash with zero or negative offset", () => {
    cy.visitSpace({ qs: { offset: 0 } });
    cy.getFolderCount("text").should("eq", "Showing 1-20 of 21 folder(s)");

    cy.visitSpace({ qs: { offset: -1 } });
    cy.getFolderCount("text").should("eq", "Showing 1-20 of 21 folder(s)");
  });

  function goToPage(page: number | "previous" | "next") {
    if (typeof page === "string") {
      cy.get(`@${page}`).random().click();
    } else {
      cy.contains(".pagination .page-item", page).click();
    }
  }

  function getNumberOfPages() {
    return cy
      .get(
        '.card-header .pagination .page-item:has(.page-link[aria-label="Next"])',
      )
      .prev()
      .invoke("text")
      .then(parseInt);
  }

  function getPagerButtons() {
    return cy
      .get(".card-header .pagination .page-item")
      .map((item: HTMLElement) => {
        const $item = Cypress.$(item);

        if ($item.is(":has(.if-chevron-left)")) {
          return "<";
        }

        if ($item.is(":has(.if-chevron-right)")) {
          return ">";
        }

        if ($item.is(":has(.if-more)")) {
          return "...";
        }

        return parseInt($item.text());
      });
  }
});
