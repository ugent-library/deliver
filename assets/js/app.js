import BSN from "bootstrap.native/dist/bootstrap-native-v4";
import bootstrapPopper from './bootstrap_popper.js'
import formChangeSubmit from './form_change_submit.js'
import formConfirm from './form_confirm.js'

// initialize everyting
document.addEventListener('DOMContentLoaded', function () {
    BSN.initCallback()
    bootstrapPopper()
    formChangeSubmit()
    formConfirm()
});
