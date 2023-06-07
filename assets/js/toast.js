export default function(rootEl) {
    rootEl.querySelectorAll('.toast [data-bs-dismiss="toast"]').forEach((btn) => {
        let el = btn.closest('.toast')
        bootstrap.Toast.getOrCreateInstance(el).show()
    })
}
