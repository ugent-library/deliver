export function logCommand(
  name,
  consoleProps = {},
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
    consoleProps: () => ({ props: consoleProps }),
  });
}

export function updateLogMessage(log: Cypress.Log, append: unknown) {
  if (!log) return;

  const message = log.get("message").split(", ").filter(Boolean);

  message.push(append);

  log.set("message", message.join(", "));
}

export function updateConsoleProps(
  log: Cypress.Log,
  callback: (ObjectLike) => void
) {
  if (!log) return;

  const consoleProps = log.get("consoleProps")();

  callback(consoleProps.props);

  log.set({ consoleProps: () => consoleProps });
}
