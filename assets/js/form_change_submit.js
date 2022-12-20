export default function() {
    document.querySelectorAll("form.change-submit").forEach(function(el) {
        el.addEventListener("change", function(evt) {
            let btn = el.querySelector("button[type='submit']")
            if (btn !== null) {
                btn.disabled = true
                let loadingText = btn.dataset.loading || "Loading..."
                btn.innerHTML = '<span class="spinner-border" role="status" aria-hidden="true"></span> ' + 
                    loadingText
            }
            el.submit()
        });
    });
}