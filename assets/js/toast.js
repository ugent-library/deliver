export default function (rootEl) {
    rootEl.querySelectorAll('[data-dismiss="toast"]').forEach((el) => {
        el.Toast.show()
    })
}