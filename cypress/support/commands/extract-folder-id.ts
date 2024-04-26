import { logCommand } from "./helpers";

export default function extractFolderId(alias: false | string = "folderId") {
  const log = logCommand("extractFolderId");

  const REGEX = /^\/folders\/(?<folderId>\w{26})$/;

  const chain = cy.location("pathname", { log: false }).then((pathname) => {
    const matches = pathname.match(REGEX);
    if (matches) {
      return matches.groups!["folderId"];
    }

    // Only assert when there is a problem so the command log does not get bloated
    expect(pathname).to.match(
      REGEX,
      "Folder ID cannot be extracted from the URL",
    );
  });

  if (alias) {
    chain.as(alias);

    // @ts-ignore 'alias' is not a keyof LogConfig but does work
    log.set("alias", `@${alias}`);
  }

  chain.finishLog(log, true);
}

declare global {
  namespace Cypress {
    interface Chainable {
      extractFolderId(alias?: false | string): Chainable<string> | never;
    }
  }
}
