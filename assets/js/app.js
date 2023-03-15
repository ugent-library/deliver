import BSN from 'bootstrap.native/dist/bootstrap-native-v4'
import bootstrapPopper from './bootstrap_popper.js'
import htmx from 'htmx.org'
import * as Turbo from '@hotwired/turbo'
import toast from './toast.js'
import formSubmit from './form_submit.js'
import formUploadProgress from './form_upload_progress.js'
import modalClose from './modal_close.js'
import clipboard from './clipboard.js'

// configure htmx
htmx.config.defaultFocusScroll = true
htmx.onLoad(function(rootEl) {
    BSN.initCallback(rootEl)
    bootstrapPopper(rootEl)
    toast(rootEl)
    formSubmit(rootEl)
    formUploadProgress(rootEl)
    modalClose(rootEl)
    clipboard(rootEl)

    window.htmx = htmx

    document.body.addEventListener('htmx:configRequest', evt => {
        evt.detail.headers['X-CSRF-Token'] = document.querySelector('meta[name="csrf-token"]').content
    })

    let ws = new WebSocket("ws://" + document.location.host + "/ws")

    Turbo.connectStreamSource(ws);

    document.querySelectorAll('.breadcrumb').forEach(el => {
        let n = 0
        el.addEventListener('click', evt => {
            evt.preventDefault()
            n++
            ws.send(JSON.stringify({
                type: "turbo",
                body: {
                    route: 'home',
                    params: {
                        name: "Matthias " + n
                    }    
                }
            }))
        })
    })
});

// document.body.addEventListener('turbo:load', evt => {
//     BSN.initCallback(rootEl)
//     bootstrapPopper(rootEl)
//     toast(rootEl)
//     formSubmit(rootEl)
//     formUploadProgress(rootEl)
//     modalClose(rootEl)
//     clipboard(rootEl)    
// })
