export function getRandomText(includeSubstring?: string): string {
  const randomText =Math.random().toString(36).split('.').at(1)!.toUpperCase()

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
