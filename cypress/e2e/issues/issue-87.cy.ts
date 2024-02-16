// https://github.com/ugent-library/deliver/issues/87

import * as dayjs from "dayjs";

import { getRandomText } from "support/util";

describe("Issue #87: Postpone button (extend folder expiration date by one month)", () => {
  let FOLDER_NAME: string;

  beforeEach(() => {
    FOLDER_NAME = getRandomText();

    cy.loginAsSpaceAdmin();

    cy.visitSpace();

    cy.setFieldByLabel("Folder name", FOLDER_NAME);
    cy.contains(".btn", "Make folder").click();

    cy.ensureToast().closeToast();

    cy.url().as("adminUrl");

    cy.extractFolderId().then((folderId) => {
      cy.intercept("PUT", `/folders/${folderId}/postpone`).as(
        "postponeExpiration"
      );
    });
  });

  it("should display a postpone button that opens a modal dialog to postpone the expiration", () => {
    cy.loginAsSpaceAdmin();

    cy.visit("@adminUrl");

    cy.contains("expires on")
      .should("be.visible")
      .invoke("text")
      .then((expiresOn: string) => {
        const expiresOnDate = dayjs(
          expiresOn.match(/^expires on (?<date>.*)$/).groups["date"]
        );

        // Allow a 2 minute margin to account for computer time glitches
        const lbound = dayjs()
          .second(0)
          .millisecond(0)
          .add(31, "days")
          .subtract(1, "minute");
        const ubound = lbound.clone().add(2, "minutes");
        expect(expiresOnDate.unix()).is.within(
          lbound.unix(),
          ubound.unix(),
          `Expiration date "${expiresOnDate.toDate()}" is not within "${lbound.toDate()}" and "${ubound.toDate()}"`
        );

        cy.wrap(expiresOnDate.format("YYYY-MM-DD")).as("expirationDate");
      });

    cy.ensureNoModal();

    cy.contains(".btn", "Postpone expiration").should("be.visible").click();

    cy.ensureModal(
      new RegExp(
        `Postpone the expiration date of\s*${FOLDER_NAME} by one month`
      )
    )
      .within(function () {
        cy.contains(`Current expiration date: ${this.expirationDate}`).should(
          "be.visible"
        );
        // Since tests run instantly, this will be the same date
        cy.contains(
          `Expiration date after postponing: ${this.expirationDate}`
        ).should("be.visible");
      })
      .closeModal("Postpone");

    cy.wait("@postponeExpiration").should(
      "have.nested.property",
      "response.statusCode",
      200
    );

    cy.url().should("eq", "@adminUrl");

    cy.ensureNoModal();

    cy.get<string>("@expirationDate").then((expirationDate) => {
      cy.ensureToast(
        `New expiration date for ${FOLDER_NAME}: ${expirationDate}`
      );
    });

    cy.wait(3500);

    cy.ensureNoToast({ timeout: 0 });
  });

  it("should not trigger the expiration logic when the modal is cancelled", () => {
    cy.loginAsSpaceAdmin();

    cy.visit("@adminUrl");

    cy.ensureNoModal();

    cy.contains(".btn", "Postpone expiration").should("be.visible").click();

    cy.ensureModal(
      new RegExp(
        `Postpone the expiration date of\s*${FOLDER_NAME} by one month`
      )
    ).closeModal("No, cancel");

    cy.get("@postponeExpiration").should("be.null");

    cy.ensureNoModal();

    cy.url().should("eq", "@adminUrl");
  });
});
