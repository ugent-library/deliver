import { logCommand, updateConsoleProps } from "./helpers";

const NO_LOG = { log: true };

type CloseModalOptions = {
  log?: boolean;
};

export default function closeModal(
  subject: unknown,
  save: boolean | string | RegExp = false,
  options: CloseModalOptions = { log: true }
): void {
  const dismissButtonText =
    typeof save === "boolean" ? (save ? "Save" : "Cancel") : save;

  let log: Cypress.Log | undefined;
  if (options.log === true) {
    log = logCommand(
      "closeModal",
      {
        subject: subject ? (subject as JQuery<HTMLElement>).get(0) : null,
        "Dismiss button text": dismissButtonText,
      },
      dismissButtonText
    );
    log.set("type", !subject ? "parent" : "child");
  }

  const doCloseModal = () => {
    cy.contains(".modal-footer .btn", dismissButtonText, NO_LOG)
      .then(($el) => {
        if (log) {
          log.set("$el", $el);
          updateConsoleProps(log, (cp) => {
            cp["Button element"] = $el.get(0);
          });
        }

        return $el;
      })
      .click(NO_LOG);
  };

  if (subject) {
    cy.wrap(subject, NO_LOG).within(NO_LOG, doCloseModal);
  } else {
    doCloseModal();
  }
}

declare global {
  namespace Cypress {
    interface Chainable {
      closeModal(
        save: boolean | string | RegExp,
        options?: CloseModalOptions
      ): Chainable<void>;
    }
  }
}
