export function mapAlias(
  thisObject: Record<string, unknown>,
  argument
): unknown {
  if (typeof argument === "string" && argument.startsWith("@")) {
    const alias = argument.slice(1);
    if (Object.keys(thisObject).includes(alias)) {
      return thisObject[alias];
    }

    throw new Error(`Alias '${argument}' was not found in test scope.`);
  }

  return argument;
}
