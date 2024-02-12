import * as bs from "bootstrap";

export default function (rootEl) {
  rootEl.querySelectorAll('.toast [data-bs-dismiss="toast"]').forEach((btn) => {
    let el = btn.closest(".toast");
    bs.Toast.getOrCreateInstance(el).show();
  });
}
