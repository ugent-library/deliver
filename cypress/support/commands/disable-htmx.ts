import { logCommand } from "./helpers";

export default function disableHtmx() {
  logCommand("disableHtmx");

  cy.on("window:before:load", (window) => {
    window.document.documentElement.setAttribute("hx-disable", "");
  });
}

declare global {
  namespace Cypress {
    interface Chainable {
      disableHtmx(): Chainable<void>;
    }
  }
}
