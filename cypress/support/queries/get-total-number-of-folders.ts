import {
  logCommand,
  updateConsoleProps,
  updateLogMessage,
} from "support/commands/helpers";

export default function getTotalNumberOfFolders(): () => number {
  const log = logCommand("getTotalNumberOfFolders");

  const regex = /Showing \d+ of (?<total>\d+) folders/;

  const containsFn = cy.now("contains", regex, {
    log: false,
  }) as () => JQuery<HTMLElement>;

  return () => {
    const total = parseInt(containsFn().text().match(regex).groups["total"]);

    if (log.get("message") != total) {
      updateLogMessage(log, total);
    }

    updateConsoleProps(log, (cp) => (cp.yielded = total));

    return total;
  };
}

declare global {
  namespace Cypress {
    interface Chainable {
      getTotalNumberOfFolders(): Chainable<number>;
    }
  }
}
