export default function() {
    el.querySelectorAll(".modal-close").forEach(function (btn) {
        btn.addEventListener("click", evt => {
            document.getElementById("modals").innerHTML = ""
        });
    });
}
