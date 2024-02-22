import {
  logCommand,
  updateConsoleProps,
  updateLogMessage,
} from "support/commands/helpers";

type GetFolderCountResult = {
  text: string;
  start?: number;
  end?: number;
  count: number;
  total: number;
};

type X = GetFolderCountResult["end"];

export default function getFolderCount<T extends keyof GetFolderCountResult>(
  key?: T
): () => GetFolderCountResult | GetFolderCountResult[T] {
  const log = logCommand("getFolderCount", { key });

  const regex =
    /Showing ((?<start>\d+)-(?<end>\d+) of )?(?<total>\d+) folder\(s\)/;

  const containsFn = cy.now("contains", regex, {
    log: false,
  }) as () => JQuery<HTMLElement>;

  return () => {
    const text = containsFn().text();
    const { start, end, total } = text.match(regex)!.groups!;

    const result: GetFolderCountResult = {
      text,
      total: parseInt(total),
      count: 0,
    };

    if (start) {
      result.start = parseInt(start);
    }

    if (end) {
      result.end = parseInt(end);
    }

    if (result.start && result.end) {
      result.count = result.end - result.start + 1;
    }

    const queryResult = key ? result[key] : result;

    if (key) {
      log.set("message", `${key}: ${result[key]}`);
      updateConsoleProps(log, (cp) => {
        cp.result = result;
      });
    } else {
      log.set(
        "message",
        start || end ? `${start || "..."}-${end || "..."} of ${total}` : total
      );
    }

    updateConsoleProps(log, (cp) => (cp.yielded = queryResult));

    return queryResult;
  };
}

declare global {
  namespace Cypress {
    interface Chainable {
      getFolderCount(): Chainable<GetFolderCountResult>;

      getFolderCount<T extends keyof GetFolderCountResult>(
        key?: T
      ): Chainable<GetFolderCountResult[T]>;
    }
  }
}
