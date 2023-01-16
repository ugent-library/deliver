export default function() {
    document.querySelectorAll("form [data-confirm]").forEach((btn) => {
        btn.addEventListener("click", (evt) => {
            evt.preventDefault()
            let tmpl = document.getElementById("tmpl-confirm").content.firstElementChild.cloneNode(true)
            tmpl.querySelector('.confirm-text').innerText = btn.dataset.confirm
            tmpl.querySelector('.confirm-proceed').addEventListener("click", () => {
                btn.closest("form").submit()
            })
            document.getElementById("modals").replaceChildren(tmpl)
        });
    });
}