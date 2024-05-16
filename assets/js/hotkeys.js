import htmx from "htmx.org";

export default function (rootEl) {
  const hotkeyElements = Array.from(rootEl.querySelectorAll("[data-hotkey]"));
  if (hotkeyElements.length) {
    htmx.on("keyup", (evt) => {
      const el = hotkeyElements.find((el) =>
        el.matches(`[data-hotkey="${evt.key}"]`),
      );

      // Make sure the "hotkey-element" is still in the DOM
      // and the triggering element is not a form field
      if (el && el.isConnected && !isFormField(evt.target)) {
        el.dispatchEvent(new MouseEvent("click"));
      }
    });
  }
}

function isFormField(el) {
  return el.matches(
    "input:not([type=button]):not([type=reset]):not([type=submit]):not([type=image]):not([type=file]), textarea, select",
  );
}
