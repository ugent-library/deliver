import { logCommand } from "./helpers";

const NO_LOG = { log: false };

export default function login(username: string, password: string): void {
  // WARNING: Whenever you change the code of the session setup, Cypress will throw an error:
  //   This session already exists. You may not create a new session with a previously used identifier.
  //   If you want to create a new session with a different setup function, please call cy.session() with
  //   a unique identifier other than...
  //
  // Temporarily uncomment the following line to clear the sessions if this happens
  // Cypress.session.clearAllSavedSessions()

  logCommand("login", { username }, username);

  cy.session(
    username,
    () => {
      cy.request({ ...NO_LOG, url: "/login", followRedirect: false }).then(
        (response) => {
          if (typeof response.headers.location !== "string") {
            throw new Error("No location header found in response");
          }

          const formUrl = new URL(response.headers.location);
          const formData = [...formUrl.searchParams];
          formData.push(["username", username]);

          if (password) {
            formData.push(["password", password]);
          }

          cy.request({
            ...NO_LOG,
            method: "POST",
            url: formUrl.origin + formUrl.pathname,
            body: new URLSearchParams(formData).toString(),
            headers: {
              "Content-Type": "application/x-www-form-urlencoded",
            },
          });
        },
      );
    },
    {
      cacheAcrossSpecs: true,
    },
  );
}

declare global {
  namespace Cypress {
    interface Chainable {
      login(username: string, password: string): Chainable<void>;
    }
  }
}
