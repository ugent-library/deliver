import toast from './toast.js'
import formSubmit from './form_submit.js'
import formUploadProgress from './form_upload_progress.js'

window.addEventListener("DOMContentLoaded", (evt) => {
    toast(document)
    formSubmit(document)
    formUploadProgress(document)
})

// old above, new below

import '@hotwired/turbo'
import { Application } from "@hotwired/stimulus"
import BaseController from "./controllers/base_controller"
import ClipboardController from "./controllers/clipboard_controller"

const app = Application.start()
app.register("base", BaseController)
app.register("clipboard", ClipboardController)
