import BSN from "bootstrap.native/dist/bootstrap-native-v4";
import bootstrapPopper from './bootstrap_popper.js'
import formChangeSubmit from './form_change_submit.js'
import formUploadProgress from './form_upload_progress.js'
import formConfirm from './form_confirm.js'
import toast from './toast.js'

// initialize everyting
document.addEventListener('DOMContentLoaded', function () {
    BSN.initCallback()
    bootstrapPopper()
    formChangeSubmit()
    formUploadProgress()
    formConfirm()
    toast()
});
