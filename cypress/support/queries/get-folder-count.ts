import { logCommand, updateConsoleProps } from "support/commands/helpers";

type GetFolderCountResult = {
  text: string;
  start?: number;
  end?: number;
  count: number;
  total: number;
};

export default function getFolderCount<T extends keyof GetFolderCountResult>(
  key?: T
): () => GetFolderCountResult | GetFolderCountResult[T] {
  const log = logCommand("getFolderCount", { key });

  const getFn = cy.now(
    "get",
    ".bc-toolbar:has(.pagination) .bc-toolbar-right .bc-toolbar-item span",
    { log: false }
  ) as () => JQuery<HTMLSpanElement>;

  return () => {
    const texts = Cypress._.uniq(
      getFn().map((_, e) => e.textContent?.trim() || "")
    );
    if (texts.length === 0 || texts?.at(0) === "") {
      throw new Error("Found no folder count messages to parse.");
    }

    if (texts.length > 1) {
      throw new Error(
        "Found multiple non-matching folder counts: \n- " + texts.join("\n - ")
      );
    }

    const text = texts.at(0)!;
    const matches = text.match(
      /Showing ((?<start>\d+)-(?<end>\d+) of )?(?<total>\d+) folder\(s\)/
    );

    if (!matches) {
      throw new Error(`Could not parse folder count message: ${text}`);
    }
    const { start, end, total } = matches.groups!;

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
