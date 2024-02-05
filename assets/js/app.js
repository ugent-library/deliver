import htmx from 'htmx.org'
import * as bs from 'bootstrap'
import toast from './toast.js'
import formSubmit from './form_submit.js'
import formUploadProgress from './form_upload_progress.js'
import clipboard from './clipboard.js'
import selectValue from './select_value.js'

window.htmx = htmx

// load htmx extensions
require('htmx.org/dist/ext/ws.js')

htmx.config.defaultFocusScroll = true

htmx.onLoad(function (el) {
  toast(el)
  formSubmit(el)
  formUploadProgress(el)
  clipboard(el)
  selectValue(el)
})

htmx.on('htmx:config-request', evt => {
  evt.detail.headers['X-CSRF-Token'] = document.querySelector('meta[name="csrf-token"]').content
})

htmx.on('htmx:confirm', evt => {
  let el = evt.detail.elt

  if (el.dataset.confirm) {
    evt.preventDefault()

    let modalEl = document.getElementById('modal-confirm').content.firstElementChild.cloneNode(true)

    document.body.append(modalEl)

    if (el.dataset.confirmHeader) {
      modalEl.querySelector('.confirm-header').innerHTML = el.dataset.confirmHeader
    }
    if (el.dataset.confirmContent) {
      modalEl.querySelector('.confirm-content').innerHTML = el.dataset.confirmContent
    }
    if (el.dataset.confirmProceed) {
      modalEl.querySelector('.confirm-proceed').innerHTML = el.dataset.confirmProceed
    }

    modalEl.querySelector('.confirm-proceed').addEventListener(
      'click',
      () => {
        evt.detail.issueRequest()
      },
      false
    )

    modalEl.addEventListener(
      'hidden.bs.modal',
      () => {
        modalEl.remove()
      },
      false
    )

    new bs.Modal(modalEl).show()
  }
})
