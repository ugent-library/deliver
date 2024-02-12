export default function (rootEl) {
  rootEl.querySelectorAll("[data-select-value]").forEach((input) => {
    input.addEventListener("click", () => {
      input.select();
    });
  });
}
