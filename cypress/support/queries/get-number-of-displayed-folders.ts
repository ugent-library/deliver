import {
  logCommand,
  updateConsoleProps,
  updateLogMessage,
} from "support/commands/helpers";

export default function getNumberOfDisplayedFolders(): () => number {
  const log = logCommand("getNumberOfDisplayedFolders");

  const regex = /Showing (?<count>\d+) of \d+ folders/;

  const containsFn = cy.now("contains", regex, {
    log: false,
  }) as () => JQuery<HTMLElement>;

  return () => {
    const total = parseInt(containsFn().text().match(regex).groups["count"]);

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
      getNumberOfDisplayedFolders(): Chainable<number>;
    }
  }
}
