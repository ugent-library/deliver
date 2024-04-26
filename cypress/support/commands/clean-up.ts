import { logCommand, updateConsoleProps, updateLogMessage } from "./helpers";

const SELECTOR = "table.table tr td:first-of-type a";

const NO_LOG = { log: false };

export default function cleanUp(space?: string) {
  const log = logCommand("cleanUp", { space }, space);

  deleteAllTestFolders(log, space || Cypress.env("DEFAULT_SPACE")).then(
    (count) => {
      log.set("message", count.toString());
    },
  );

  cy.get(SELECTOR, NO_LOG).should("not.exist");
}

function deleteAllTestFolders(log: Cypress.Log, space: string) {
  return cy.visitSpace(space, { qs: { limit: 1000 }, ...NO_LOG }).then(() => {
    // Using Cypress.$() direct jQuery selector tool here.
    // Using cy.get() the test would fail if none are left.
    const anchors = Cypress.$<HTMLAnchorElement>(SELECTOR).get();

    updateConsoleProps(
      log,
      (cp) => (cp.links = anchors.map((a) => a.textContent)),
    );

    const links = anchors.map((a) => a.getAttribute("href"));

    return cy
      .wrap(links, NO_LOG)
      .each((href) => {
        cy.visit(`${href}/edit`, NO_LOG);

        cy.contains(".btn", "Delete folder", NO_LOG)
          .submitForm(undefined, NO_LOG)
          .then(() => updateLogMessage(log, ".", ""));
      })
      .then(() => links.length);
  });
}

declare global {
  namespace Cypress {
    interface Chainable {
      cleanUp(space?: string): Chainable<void>;
    }
  }
}
