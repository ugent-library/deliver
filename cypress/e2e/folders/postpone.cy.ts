// https://github.com/ugent-library/deliver/issues/87

import * as dayjs from "dayjs";

import { getRandomText } from "support/util";

describe("Issue #87: Postpone button (extend folder expiration date by one month)", () => {
  let FOLDER_NAME: string;

  beforeEach(() => {
    FOLDER_NAME = getRandomText();

    cy.loginAsSpaceAdmin();

    cy.makeFolder(FOLDER_NAME);

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

    testPostpone(dayjs().add(30, "days"));

    // You can postpone multiple times
    testPostpone(dayjs().add(60, "days"));

    testPostpone(dayjs().add(90, "days"));

    // Check expiration date in folder list
    cy.visitSpace({ qs: { q: FOLDER_NAME } });

    cy.get("#folders tbody tr:first-of-type td")
      .eq(2)
      .should("contain.text", dayjs().add(120, "days").format("YYYY-MM-DD"));
  });

  it("should not trigger the expiration logic when the modal is cancelled", () => {
    cy.loginAsSpaceAdmin();

    cy.visit("@adminUrl");

    cy.ensureNoModal();

    cy.contains(".btn", "Postpone expiration").should("be.visible").click();

    cy.ensureModal(
      new RegExp(
        `Postpone the expiration date of\\s*${FOLDER_NAME} by one month`
      )
    ).closeModal("No, cancel");

    cy.get("@postponeExpiration").should("be.null");

    cy.ensureNoModal();

    cy.url().should("eq", "@adminUrl");
  });

  function testPostpone(initDate: dayjs.Dayjs) {
    cy.contains("expires on")
      .should("be.visible")
      .invoke("text")
      .then((expiresOn: string) => {
        const expiresOnDate = dayjs(getExpiresOnDate(expiresOn)).startOf("day"); // Ignore time portion

        const calculatedExpirationDate = initDate.startOf("day");

        expect(expiresOnDate.valueOf()).eq(
          calculatedExpirationDate.valueOf(),
          `Expected ${calculatedExpirationDate.toISOString()} to be equal to ${expiresOnDate.toISOString()}.`
        );

        cy.wrap(expiresOnDate.format("YYYY-MM-DD")).as("expirationDate");
        cy.wrap(expiresOnDate.add(30, "days").format("YYYY-MM-DD")).as(
          "nextExpirationDate"
        );
      });

    cy.ensureNoModal();

    cy.contains(".btn", "Postpone expiration").should("be.visible").click();

    cy.ensureModal(
      new RegExp(
        `Postpone the expiration date of\\s*${FOLDER_NAME} by one month`
      )
    )
      .within(function (this: Record<string, unknown>) {
        cy.contains(`Current expiration date: ${this.expirationDate}`).should(
          "be.visible"
        );
        // Since tests run instantly, this will be the same date
        cy.contains(
          `Expiration date after postponing: ${this.nextExpirationDate}`
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

    cy.get<string>("@nextExpirationDate").then((nextExpirationDate) => {
      cy.ensureToast(
        `New expiration date for ${FOLDER_NAME}: ${nextExpirationDate}`
      );
    });

    cy.wait(3500);

    cy.ensureNoToast({ timeout: 0 });

    cy.reload();

    cy.contains("expires on")
      .should("be.visible")
      .invoke("text")
      .should("contain", "@nextExpirationDate");
  }

  function getExpiresOnDate(expiresOn: string) {
    const matches = expiresOn.match(/^expires on (?<date>.*)$/);

    if (!matches) {
      throw new Error('Could not extract date from "expires on" text.');
    }

    return matches.groups!["date"];
  }
});
