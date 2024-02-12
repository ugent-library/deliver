export default function getRandomText() {
  return crypto.randomUUID().replace(/-/g, "").toUpperCase().substring(0, 10);
}
