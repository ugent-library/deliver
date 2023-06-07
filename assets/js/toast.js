export default function(rootEl) {
    rootEl.querySelectorAll('.toast').forEach((el) => {
        bootstrap.Toast.getOrCreateInstance(el).show()
    })
}
