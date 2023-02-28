export default function(rootEl) {
    rootEl.querySelectorAll(".modal-close").forEach(function (btn) {
        btn.addEventListener("click", evt => {
            document.getElementById("modals").innerHTML = ""
        });
    });
}
