export default function () {
    document.querySelectorAll('[data-dismiss="toast"]').forEach((el) => {
        el.Toast.show()
    })
}