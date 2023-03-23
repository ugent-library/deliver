import Popper from 'popper.js'

// Wire in popper.js support. This ensure popups stays within the viewport.
// See: https://github.com/thednp/bootstrap.native/issues/211
export default function(el) {
    el.querySelectorAll("div.dropdown > button").forEach(function(button) {
        button.addEventListener("click", function(evt) {
            let menu = button.parentElement.children.item(1);

            if (menu.classList.contains("show")) {
                menu.removeAttribute("x-placement");
                menu.removeAttribute("style");

                let popper = new Popper(button, menu, {
                    modifiers: {
                        preventOverflow: { enabled: true },
                        flip: { enabled: true},
                        hide: { enabled: false}
                    }
                });
            }
        })
    });
}