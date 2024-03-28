export function getRandomText(includeSubstring?: string): string {
  const randomText = crypto
    .randomUUID()
    .replace(/-/g, "")
    .toUpperCase()
    .substring(0, 10);

  if (includeSubstring) {
    return (
      randomText.slice(0, 5) +
      " " +
      includeSubstring +
      " " +
      randomText.slice(5)
    );
  }

  return randomText;
}

export function mapAlias(
  thisObject: Record<string, unknown>,
  argument: unknown
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
