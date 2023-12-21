import { logCommand } from './helpers'

export default function getClipboardText(): Cypress.Chainable<string> {
  const log = logCommand('getClipboardText')

  return cy
    .window({ log: false })
    .then(win => win.navigator.clipboard.readText())
    .finishLog(log)
}

declare global {
  namespace Cypress {
    interface Chainable {
      getClipboardText(): Chainable<string>
    }
  }
}
