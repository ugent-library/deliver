import { logCommand } from "support/commands/helpers";

type GetLabelOptions = {
  log?: boolean;
};

export default function (
  caption: string | RegExp,
  options: GetLabelOptions = { log: true },
) {
  const log =
    options.log !== false && logCommand("getLabel", undefined, caption);

  const getFn = cy.now("get", "label", {
    log: false,
  }) as () => JQuery<HTMLElement>;

  return (): JQuery<HTMLElement> => {
    const $el = getFn().filter((_, el) => {
      const $currentLabel = Cypress.$(el).clone();

      $currentLabel.find(".badge, .visually-hidden").remove();

      const currentLabelText = $currentLabel
        .text()
        .trim()
        .split(/\s/)
        .join(" ");

      if (typeof caption === "string") {
        return currentLabelText === caption;
      }

      return caption.test(currentLabelText);
    });

    if (log && cy.state("current") === this) {
      log
        .set({
          $el,
          consoleProps: () => ({
            Caption: caption,
            Yielded: $el?.length ? $el[0] : "--nothing--",
            Elements: $el != null ? $el.length : 0,
          }),
        })
        .finish();
    }

    return $el;
  };
}

declare global {
  namespace Cypress {
    interface Chainable {
      getLabel(
        caption: string | RegExp,
        options?: GetLabelOptions,
      ): Chainable<JQuery<HTMLLabelElement>>;
    }
  }
}
