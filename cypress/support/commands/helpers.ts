export function logCommand(
  name: string,
  consoleProps: Cypress.ObjectLike = {},
  message: unknown = "",
  $el = undefined
) {
  return Cypress.log({
    $el,
    name,
    displayName: name
      .replace(/([A-Z])/g, " $1")
      .trim()
      .toUpperCase(),
    message,
    consoleProps: () => consoleProps,
  });
}

export function updateLogMessage(
  log: Cypress.Log | undefined,
  append: unknown
) {
  if (!log) return;

  const message = log.get("message").split(", ").filter(Boolean);

  message.push(append);

  log.set("message", message.join(", "));
}

export function updateConsoleProps(
  log: Cypress.Log | undefined,
  callback: (consoleProps: Cypress.ObjectLike) => void
) {
  if (!log) return;

  const consoleProps = log.get("consoleProps")();

  callback(consoleProps);

  log.set({ consoleProps: () => consoleProps });
}
