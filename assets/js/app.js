import {Turbo} from '@hotwired/turbo'
import toast from './toast.js'
import formSubmit from './form_submit.js'
import formUploadProgress from './form_upload_progress.js'
import modalClose from './modal_close.js'
import clipboard from './clipboard.js'

window.addEventListener("DOMContentLoaded", (evt) => {
    toast(document)
    formSubmit(document)
    formUploadProgress(document)
    modalClose(document)
    clipboard(document)
})

// old above, new below

import { Application } from "@hotwired/stimulus"
import BaseController from "./controllers/base_controller"

const app = Application.start()
app.register("base", BaseController)
