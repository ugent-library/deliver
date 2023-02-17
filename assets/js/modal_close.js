import htmx from 'htmx.org';

export default function() {
    let modalClose = function(evt) {
      document.getElementById("modals").innerHTML = ""
    }
    htmx.onLoad(function(el) {
        el.querySelectorAll(".modal-close").forEach(function (btn) {
            btn.addEventListener("click", modalClose);
        });
    });
}
