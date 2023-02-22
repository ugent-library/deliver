import BSN from "bootstrap.native/dist/bootstrap-native-v4";
import bootstrapPopper from './bootstrap_popper.js'
import toast from './toast.js'
import formSubmit from './form_submit.js'
import formUploadProgress from './form_upload_progress.js'
import clipboard from './clipboard.js'

// initialize everyting
document.addEventListener('DOMContentLoaded', function () {
    BSN.initCallback()
    bootstrapPopper()
    toast()
    formSubmit()
    formUploadProgress()
    clipboard()
});
