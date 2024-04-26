import { logCommand } from "./helpers";

export default function setClipboardText(
  text: string,
): Cypress.Chainable<void> {
  logCommand("setClipboardText", { text }, `"${text}"`);

  return cy
    .window({ log: false })
    .then((win) => win.navigator.clipboard.writeText(text));
}

declare global {
  namespace Cypress {
    interface Chainable {
      setClipboardText(text: string): Chainable<void>;
    }
  }
}
