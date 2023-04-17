import formSubmit from './form_submit.js'
import formUploadProgress from './form_upload_progress.js'

window.addEventListener("DOMContentLoaded", (evt) => {
    formSubmit(document)
    formUploadProgress(document)
})

import '@hotwired/turbo'
Turbo.session.drive = false

// old above, new below

import 'htmx.org'
import { Application } from "@hotwired/stimulus"
import BaseController from "./controllers/base_controller"
import ClipboardController from "./controllers/clipboard_controller"
import ToastController from "./controllers/toast_controller"

const app = Application.start()
app.register("base", BaseController)
app.register("clipboard", ClipboardController)
app.register("toast", ToastController)
