import { logCommand, updateConsoleProps } from "./helpers";

export default function getFolderShareUrl(
  subject: string | null,
  folderIdOrAlias: string,
  folderName?: string
): Cypress.Chainable<string> {
  if (subject) {
    // Handle child mode
    folderName = folderIdOrAlias;
    folderIdOrAlias = subject;
  }

  const folderIdAlias: string | null = folderIdOrAlias.startsWith("@")
    ? folderIdOrAlias
    : null;
  const folderId: string = !!folderIdAlias
    ? this[folderIdAlias.slice(1)]
    : folderIdOrAlias;

  const consoleProps = { "Folder ID": folderId, "Folder name": folderName };
  const log = logCommand("getFolderShareUrl", consoleProps);

  if (subject) {
    log.set("type", "child");
  }

  if (folderIdAlias) {
    updateConsoleProps(log, (cp) => {
      cp["Folder ID alias"] = folderIdAlias;
    });
  }

  folderName = folderName.replace(/[^a-zA-Z0-9]+/g, "-"); // Normalize non-alphanumeric characters

  const url = new URL(
    `/share/${folderIdOrAlias}:${folderName}`,
    Cypress.config().baseUrl
  );

  return cy.wrap(url.toString(), { log: false }).finishLog(log);
}

declare global {
  namespace Cypress {
    interface Chainable {
      getFolderShareUrl(
        folderIdOrAlias: string,
        folderName: string
      ): Chainable<string>;

      getFolderShareUrl(folderName: string): Chainable<string>;
    }
  }
}
