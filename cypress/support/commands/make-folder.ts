import { logCommand } from "./helpers";

const NO_LOG = { log: false };

export default function makeFolder(
  name: string,
  space: string = Cypress.env("DEFAULT_SPACE")
): void {
  logCommand("makeFolder", { name, space }, name);

  cy.location("pathname", { log: false }).then((pathname) => {
    if (pathname !== `/spaces/${space}`) {
      cy.visitSpace(space, NO_LOG);
    }
  });

  cy.setFieldByLabel("Folder name", name, NO_LOG);
  cy.contains(".btn", "Make folder", NO_LOG).click(NO_LOG);

  cy.ensureToast("Folder created successfully", NO_LOG).closeToast(NO_LOG);
}

declare global {
  namespace Cypress {
    interface Chainable {
      makeFolder(name: string, space?: string): Chainable<void>;
    }
  }
}
