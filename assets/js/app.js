import BSN from 'bootstrap.native/dist/bootstrap-native-v4';
import bootstrapPopper from './bootstrap_popper.js'
import htmx from 'htmx.org'
import toast from './toast.js'
import formSubmit from './form_submit.js'
import formUploadProgress from './form_upload_progress.js'
import modalClose from './modal_close.js'

// configure htmx
htmx.config.defaultFocusScroll = true
htmx.onLoad(function(rootEl) {
    BSN.initCallback(rootEl)
    bootstrapPopper(rootEl)
    toast(rootEl)
    formSubmit(rootEl)
    formUploadProgress(rootEl)
    modalClose(rootEl)
});

window.htmx = htmx

document.body.addEventListener('htmx:configRequest', (evt) => {
    evt.detail.headers['X-CSRF-Token'] = document.querySelector('meta[name="csrf_token"]').content
})
