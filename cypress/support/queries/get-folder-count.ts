import {
  logCommand,
  updateConsoleProps,
  updateLogMessage,
} from "support/commands/helpers";

type GetFolderCountResult = {
  start?: number;
  end?: number;
  count: number;
  total: number;
};

export default function getFolderCount(): () => GetFolderCountResult {
  const log = logCommand("getFolderCount");

  const regex =
    /Showing ((?<start>\d+)-(?<end>\d+) of )?(?<total>\d+) folder\(s\)/;

  const containsFn = cy.now("contains", regex, {
    log: false,
  }) as () => JQuery<HTMLElement>;

  return () => {
    const { start, end, total } = containsFn().text().match(regex).groups;

    const result: GetFolderCountResult = { total: parseInt(total), count: 0 };

    if (start) {
      result.start = parseInt(start);
    }

    if (end) {
      result.end = parseInt(end);
    }

    if (result.start && result.end) {
      result.count = result.end - result.start + 1;
    }

    log.set(
      "message",
      start || end ? `${start || "..."}-${end || "..."} of ${total}` : total
    );

    updateConsoleProps(log, (cp) => (cp.yielded = result));

    return result;
  };
}

declare global {
  namespace Cypress {
    interface Chainable {
      getFolderCount(): Chainable<GetFolderCountResult>;
    }
  }
}
