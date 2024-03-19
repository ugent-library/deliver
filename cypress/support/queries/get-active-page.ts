import {
  logCommand,
  updateConsoleProps,
  updateLogMessage,
} from "support/commands/helpers";

export default function getActivePage(): () => number {
  const log = logCommand("getActivePage");

  const getFn = cy.now("get", ".pagination .page-item.active a.page-link", {
    log: false,
  }) as () => JQuery<HTMLAnchorElement>;

  return () => {
    const $a = getFn();

    if ($a.length !== 2) {
      expect($a).to.have.length(2);
    }

    const pageNumbers = Cypress._.uniq($a.map((_, a) => a.textContent));
    if (pageNumbers.length !== 1) {
      expect(pageNumbers).to.have.length(
        1,
        "Active page is out of sync in header and footer"
      );
    }

    const result = parseInt(pageNumbers[0]);

    updateLogMessage(log, result);
    updateConsoleProps(log, (cp) => {
      cp.yielded = result;
    });

    return result;
  };
}

declare global {
  namespace Cypress {
    interface Chainable {
      getActivePage(): Chainable<number>;
    }
  }
}
