import { logCommand } from "./helpers";

type SubmitFormOptions = {
  log?: boolean;
};

export default function submitForm<T = unknown>(
  subject: JQuery<HTMLElement>,
  fields: Record<string, string>,
  options: SubmitFormOptions,
): Cypress.Chainable<Cypress.Response<T>> {
  if (subject.length !== 1) {
    throw new Error("Expected exactly one selected FORM element.");
  }

  if (subject.prop("tagName") !== "FORM") {
    subject = subject.closest("form");
    if (subject.length !== 1) {
      throw new Error(`Expected FORM element, got ${subject.prop("tagName")}`);
    }
  }

  const form = subject as JQuery<HTMLFormElement>;
  const method = form.prop("method").toUpperCase();
  const action = form.attr("action");

  const inputs = Object.fromEntries(
    form
      .find("input, select, textarea")
      .get()
      .map((i) => [i.name, i.value]),
  );

  const body = { ...inputs, ...fields };

  const log =
    options?.log !== false
      ? logCommand(
          "submitForm",
          {
            subject,
            fields,
            "All fields": body,
            method,
            action,
          },
          `${method} ${action}`,
        )
      : undefined;

  return cy
    .request<T>({
      method,
      url: action,
      body,
      form: true,
      followRedirect: false,
      log: false,
    })
    .finishLog(log);
}

declare global {
  namespace Cypress {
    interface Chainable {
      submitForm<T = unknown>(
        fields?: Record<string, string>,
        options?: SubmitFormOptions,
      ): Chainable<Response<T>>;
    }
  }
}
