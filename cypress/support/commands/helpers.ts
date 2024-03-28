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
    consoleProps: () => ({ props: consoleProps }),
  });
}

export function updateLogMessage(
  log: Cypress.Log | undefined,
  append: unknown,
  separator: string = ", "
) {
  if (!log) return;

  const message: unknown[] = (log.get("message") as string)
    .split(separator)
    .filter(Boolean);

  message.push(append);

  log.set("message", message.join(separator));
}

export function updateConsoleProps(
  log: Cypress.Log | undefined,
  callback: (consoleProps: Cypress.ObjectLike) => void
) {
  if (!log) return;

  const consoleProps = log.get("consoleProps")();

  callback(consoleProps.props);

  log.set({ consoleProps: () => consoleProps });
}
