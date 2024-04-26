import { logCommand } from "./helpers";

const NO_LOG = { log: false };

type MakeFolderOptions = {
  space?: string;
  noRedirect?: boolean;
};

export default function makeFolder(
  name: string,
  { space, noRedirect }: MakeFolderOptions = {},
): void {
  space ||= Cypress.env("DEFAULT_SPACE");

  logCommand("makeFolder", { name, space, "No redirect": noRedirect }, name);

  cy.location("pathname", { log: false }).then((pathname) => {
    if (pathname !== `/spaces/${space}`) {
      cy.visitSpace(space, NO_LOG);
    }
  });

  cy.getLabel("Folder name", NO_LOG)
    .submitForm({ name }, NO_LOG)
    .then((response) => {
      if (!noRedirect && response.redirectedToUrl) {
        cy.visit(response.redirectedToUrl, NO_LOG);
      }
    });
}

declare global {
  namespace Cypress {
    interface Chainable {
      makeFolder(name: string, options?: MakeFolderOptions): Chainable<void>;
    }
  }
}
