export default function() {
    document.querySelectorAll("form.change-submit").forEach((form) => {
        form.addEventListener("change", () => {
            let btn = form.querySelector("button[type='submit']")
            if (btn !== null) {
                btn.disabled = true
                let loadingText = btn.dataset.loading || "Loading..."
                btn.innerHTML = '<span class="spinner-border" role="status" aria-hidden="true"></span> ' + 
                    loadingText
            }
            form.submit()
        });
    });
}