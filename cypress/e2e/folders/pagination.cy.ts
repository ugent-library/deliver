// https://github.com/ugent-library/deliver/issues/91

import { getRandomText } from "support/util";

describe("Issue #91: [Speed and usability] Add pagination to folder overview", () => {
  before(() => {
    cy.loginAsSpaceAdmin();

    cy.visitSpace();

    // All following tests assume there are exactly 21 folders in this space
    cy.getFolderCount("total").then((total) => {
      if (total < 21) {
        for (let i = 0; i < 21 - total; i++) {
          cy.makeFolder(getRandomText());
        }
      } else if (total > 21) {
        for (let i = 0; i < total - 21; i++) {
          cy.get("#folders tbody tr:first-of-type td:first-of-type a")
            .prop("href")
            .then((href) => {
              cy.visit(`${href}/edit`);

              cy.contains(".btn", "Delete folder").click();
            });
        }
      }
    });

    cy.visitSpace();
    cy.getFolderCount("total").should("eq", 21);
  });

  beforeEach(() => {
    cy.loginAsSpaceAdmin();

    cy.visitSpace();

    cy.get('input[name="q"]').as("q");
    cy.get('select[name="sort"]').as("sort");

    cy.get('.pagination .page-item:has(.page-link[aria-label="Previous"])').as(
      "previous"
    );
    cy.get('.pagination .page-item:has(.page-link[aria-label="Next"])').as(
      "next"
    );

    cy.intercept(`/spaces/${Cypress.env("DEFAULT_SPACE")}/folders*`).as(
      "getFolders"
    );
  });

  it("should display maximum 20 folders per page by default", () => {
    cy.visitSpace();

    cy.getFolderCount("text").should("eq", "Showing 1-20 of 21 folder(s)");

    cy.contains(
      "#folders .card-header .pagination .page-item a.page-link",
      "2"
    ).should("be.visible");

    cy.get("#folders table tbody tr").should("have.length", 20);

    cy.get(".u-scroll-wrapper__body").scrollTo("bottom");

    cy.contains(
      "#folders .card-footer .pagination .page-item a.page-link",
      "2"
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

  it("should display the filtered folder count when search query is used", () => {
    cy.visitSpace();

    cy.get("@q").type("4");
    cy.contains(".btn", "Search").click();

    cy.getFolderCount("total").should("be.lessThan", 21);
  });

  it("should highlight the current page button", () => {
    cy.visitSpace({ qs: { limit: 6 } });

    cy.getFolderCount("text").should("eq", "Showing 1-6 of 21 folder(s)");

    [1, 3, 4, 2].forEach((page) => {
      goToPage(page);

      cy.get(`#folders .pagination .page-item:contains(${page})`)
        .should("have.length", 2)
        .should("have.class", "active");

      cy.get(`#folders .pagination .page-item:not(:contains(${page}))`).should(
        "not.have.class",
        "active"
      );
    });
  });

  it("should be possible to jump to the next and previous page", () => {
    cy.visitSpace({ qs: { limit: 5 } });
    cy.getFolderCount("text").should("eq", "Showing 1-5 of 21 folder(s)");

    goToPage("next");
    cy.getFolderCount("text").should("eq", "Showing 6-10 of 21 folder(s)");
    getActivePage().should("eq", 2);
    cy.param("offset").should("eq", "5");

    goToPage("next");
    cy.getFolderCount("text").should("eq", "Showing 11-15 of 21 folder(s)");
    getActivePage().should("eq", 3);
    cy.param("offset").should("eq", "10");

    goToPage("previous");
    cy.getFolderCount("text").should("eq", "Showing 6-10 of 21 folder(s)");
    getActivePage().should("eq", 2);
    cy.param("offset").should("eq", "5");

    goToPage("previous");
    cy.getFolderCount("text").should("eq", "Showing 1-5 of 21 folder(s)");
    getActivePage().should("eq", 1);
    cy.param("offset").should("not.exist");
  });

  it("should disable the previous page button when on the first page", () => {
    cy.visitSpace({ qs: { limit: 5 } });
    getActivePage().should("eq", 1);
    cy.get("@previous").should("have.class", "disabled");

    // Make sure nothing happens if we click it
    cy.get("@previous").first().click();
    cy.get("@previous").last().click();
    cy.get("@getFolders").should("be.null");

    // Test again after paging
    goToPage("next");
    getActivePage().should("eq", 2);
    cy.get("@previous").should("not.have.class", "disabled");

    goToPage("previous");
    getActivePage().should("eq", 1);
    cy.get("@previous").should("have.class", "disabled");
  });

  it("should disable the next page button when on the last page", () => {
    cy.visitSpace({ qs: { limit: 5, offset: 20 } });
    getActivePage().should("eq", 5);
    cy.get("@next").should("have.class", "disabled");

    // Make sure nothing happens if we click it
    cy.get("@next").first().click();
    cy.get("@next").last().click();
    cy.get("@getFolders").should("be.null");

    // Test again after paging
    goToPage("previous");
    getActivePage().should("eq", 4);
    cy.get("@next").should("not.have.class", "disabled");

    goToPage("next");
    getActivePage().should("eq", 5);
    cy.get("@next").should("have.class", "disabled");
  });

  it("should disable previous and next page buttons when all results fit on one page (without search query)", () => {
    cy.visitSpace({ qs: { limit: 21 } });
    getActivePage().should("eq", 1);

    cy.get("@previous").should("have.class", "disabled");
    cy.get("@next").should("have.class", "disabled");
  });

  it("should disable previous and next page buttons when all results fit on one page (with query)", () => {
    cy.visitSpace();
    getActivePage().should("eq", 1);

    cy.get("@previous").should("have.class", "disabled");
    cy.get("@next").should("not.have.class", "disabled");

    cy.get("@q").type("4");

    cy.wait("@getFolders");

    cy.get("@previous").should("have.class", "disabled");
    cy.get("@next").should("have.class", "disabled");
  });

  it(
    "should display an ellipsis section in the pager when there are a lot of pages"
  );

  it("should keep search query when switching pages", () => {
    cy.visitSpace({ qs: { limit: 3 } });
    cy.param("q").should("be.null");

    cy.get("@q").type("4");
    cy.wait("@getFolders");

    cy.getFolderCount("total").as("previousTotal").should("be.lessThan", 21);

    goToPage(3);
    cy.param("q").should("eq", "4");
    cy.get("@q").should("have.value", "4");
    cy.getFolderCount("total").should("eq", "@previousTotal");

    goToPage("previous");
    cy.param("q").should("eq", "4");
    cy.get("@q").should("have.value", "4");
    cy.getFolderCount("total").should("eq", "@previousTotal");

    goToPage("next");
    cy.param("q").should("eq", "4");
    cy.get("@q").should("have.value", "4");
    cy.getFolderCount("total").should("eq", "@previousTotal");
  });

  it("should keep sorting when switching pages", () => {
    cy.visitSpace({ qs: { limit: 3 } });
    cy.param("sort").should("be.null");

    cy.setFieldByLabel("Sort by", "expires-last");
    cy.wait("@getFolders");

    cy.getFolderCount("total").should("eq", 21);

    goToPage(3);
    cy.param("sort").should("eq", "expires-last");
    cy.get("@sort").should("have.value", "expires-last");
    cy.getFolderCount("total").should("eq", 21);

    goToPage("previous");
    cy.param("sort").should("eq", "expires-last");
    cy.get("@sort").should("have.value", "expires-last");
    cy.getFolderCount("total").should("eq", 21);

    goToPage("next");
    cy.param("sort").should("eq", "expires-last");
    cy.get("@sort").should("have.value", "expires-last");
    cy.getFolderCount("total").should("eq", 21);
  });

  it("should keep the page limit when switching pages", () => {
    cy.visitSpace({ qs: { limit: 6 } });
    cy.param("limit").should("eq", "6");
    cy.getFolderCount("text").should("eq", "Showing 1-6 of 21 folder(s)");

    goToPage(3);
    cy.param("limit").should("eq", "6");
    cy.getFolderCount("text").should("eq", "Showing 13-18 of 21 folder(s)");

    goToPage("previous");
    cy.param("limit").should("eq", "6");
    cy.getFolderCount("text").should("eq", "Showing 7-12 of 21 folder(s)");

    goToPage("next");
    cy.param("limit").should("eq", "6");
    cy.getFolderCount("text").should("eq", "Showing 13-18 of 21 folder(s)");
  });

  it("should reset to first page when changing the search query");

  it("should reset to first page when sorting the results");

  function goToPage(page: number | "previous" | "next") {
    cy.document().then((document): void => {
      document.body.addEventListener(
        "htmx:pushedIntoHistory",
        cy.stub().as("htmx:pushedIntoHistory"),
        { once: true }
      );
    });

    if (typeof page === "string") {
      cy.get(`@${page}`).random().click();
    } else {
      cy.contains(".pagination .page-item", page).click();
    }

    // Make sure the new URL is pushed by HTMX.
    // Awaiting the getFolders API is not sufficient as the new URL is pushed asynchronously and Cypress would take over too soon.
    // This assertion makes sure this happened before proceeding.
    cy.get("@htmx:pushedIntoHistory").should("have.been.calledOnce");
  }

  function getActivePage(): Cypress.Chainable<number> {
    return cy
      .get(".pagination .page-item.active a.page-link")
      .should(($a) => {
        if ($a.length !== 2) {
          expect($a).to.have.length(2);
        }
      })
      .map<JQuery<HTMLAnchorElement>, string>("text")
      .then(Cypress._.uniq)
      .should((texts) => {
        if (texts.length !== 1) {
          expect(texts).to.have.length(
            1,
            "Active page is out of sync in header and pager"
          );
        }
      })
      .then((pageNumbers) => parseInt(pageNumbers[0]));
  }
});
