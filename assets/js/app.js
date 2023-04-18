import formSubmit from './form_submit.js'
import formUploadProgress from './form_upload_progress.js'

window.addEventListener("DOMContentLoaded", (evt) => {
    formSubmit(document)
    formUploadProgress(document)
})

// import '@hotwired/turbo'
// Turbo.session.drive = false

// import { Application } from "@hotwired/stimulus"
// import BaseController from "./controllers/base_controller"
// import ClipboardController from "./controllers/clipboard_controller"
// import ToastController from "./controllers/toast_controller"

// const app = Application.start()
// app.register("base", BaseController)
// app.register("clipboard", ClipboardController)
// app.register("toast", ToastController)

// old above, new below

import htmx from 'htmx.org'
import bsn from 'bootstrap.native/dist/bootstrap-native-v4'
import bsnPopper from './lib/bsn_popper.js'

window.htmx = htmx

htmx.config.defaultFocusScroll = true

htmx.onLoad(function(el) {
    bsn.initCallback(el)
    bsnPopper(el)
})

htmx.on('htmx:config-request', (evt) => {
    evt.detail.headers['X-CSRF-Token'] = document.querySelector('meta[name="csrf-token"]').content
});

htmx.on('htmx:confirm', (evt) => {
    let el = evt.detail.elt

    if (el.dataset.confirm) {
        evt.preventDefault()

        let modalEl = document.getElementById('modal-confirm').content.firstElementChild.cloneNode(true)
        document.body.appendChild(modalEl)

        if (el.dataset.confirmHeader) {
            modalEl.querySelector('.confirm-header').innerHTML = el.dataset.confirmHeader
        }
        if (el.dataset.confirmProceed) {
            modalEl.querySelector('.confirm-proceed').innerHTML = el.dataset.confirmProceed
        }

        modalEl.querySelector('.confirm-proceed').addEventListener('click', () => {
            evt.detail.issueRequest()
        }, false)

        modalEl.addEventListener('hidden.bs.modal', () => {
            modalEl.remove()
        }, false);

        new bsn.Modal(modalEl).show()
    }
});
