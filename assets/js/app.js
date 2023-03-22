import BSN from 'bootstrap.native/dist/bootstrap-native-v4'
import bootstrapPopper from './bootstrap_popper.js'
import * as Turbo from '@hotwired/turbo'
import toast from './toast.js'
import formSubmit from './form_submit.js'
import formUploadProgress from './form_upload_progress.js'
import modalClose from './modal_close.js'
import clipboard from './clipboard.js'

window.addEventListener("DOMContentLoaded", (evt) => {
    BSN.initCallback(document)
    bootstrapPopper(document)
    toast(document)
    formSubmit(document)
    formUploadProgress(document)
    modalClose(document)
    clipboard(document)
})