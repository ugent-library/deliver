import BSN from 'bootstrap.native/dist/bootstrap-native-v4';

export default function(rootEl) {
    rootEl.querySelectorAll('form input[data-upload-progress-target]').forEach(input => {
        input.addEventListener('change', () => {
            let files = Array.from(input.files)

            if (!files.length) return;

            let target = document.getElementById(input.dataset.uploadProgressTarget)
            let form = input.closest('form')
            let csrfToken = document.querySelector('meta[name=csrf_token]').content

            files.forEach((file, i) => {
                let tmpl = document.getElementById('tmpl-upload-progress').content.firstElementChild.cloneNode(true)
                tmpl.querySelector('.upload-name').innerText = file.name
                target.appendChild(tmpl)

                let hideBtnCancelUpload = function(){
                  let b = tmpl.querySelector('.btn-cancel-upload')
                  if (b == null) return
                  b.parentElement.removeChild(b)
                }
                let showBtnRemoveUpload = function(){
                  let b = tmpl.querySelector('.btn-remove-upload')
                  if (b == null) return
                  b.classList.remove('d-none')
                }
                tmpl.querySelector('.btn-remove-upload').addEventListener('click', function(){
                  let i = this.closest('.list-group-item')
                  i.parentElement.removeChild(i)
                })
                let showMessage = function(msg, level) {
                  let m = tmpl.querySelector('.upload-msg')
                  let cl = 'text-muted'
                  if (level == 'info') {
                    cl = 'text-muted'
                  } else if(level = 'error') {
                    cl = 'text-danger'
                  }
                  m.classList.remove('text-muted')
                  m.classList.add(cl)
                  m.innerText = msg
                }
                let hlPgBar = function(pgBar, level) {
                  let cl = 'bg-info'
                  if (level == 'warning') {
                    cl = 'bg-warning'
                  } else if(level == 'error') {
                    cl = 'bg-danger'
                  }
                  pgBar.classList.remove('bg-info')
                  pgBar.classList.add(cl)
                }

                // prevent file upload when above max file size
                let maxFileSize = input.dataset.uploadMaxFileSize
                if (!isNaN(maxFileSize) && file.size > maxFileSize) {

                  hideBtnCancelUpload()
                  showBtnRemoveUpload()
                  showMessage(input.dataset.uploadMsgFileTooLarge, 'error')
                  hlPgBar(tmpl.querySelector('.progress-bar'), 'error')
                  return

                }

                // send headers along with request
                let headers = [
                  ['X-CSRF-Token', csrfToken],
                  // weird, but makes sure that middleware does not try to read _method from form
                  ['X-HTTP-Method-Override', 'POST'],
                  //"Failed to execute 'setRequestHeader' on 'XMLHttpRequest': String contains non ISO-8859-1 code point"
                  ['X-Upload-Filename', encodeURIComponent(file.name)],
                  //refused by browser
                  //['Content-Length', file.size],
                  ['Content-Type', file.type]
                ]

                let req = new XMLHttpRequest();

                req.addEventListener('abort', e => {

                  showMessage(input.dataset.uploadMsgFileAborted, 'error')
                  hlPgBar(tmpl.querySelector('.progress-bar'), 'error')
                  hideBtnCancelUpload()
                  showBtnRemoveUpload()

                }, false);

                req.upload.addEventListener('progress', e => {
                    if (!e.lengthComputable) return;

                    let percent = Math.floor(e.loaded / e.total * 100)
                    tmpl.querySelector('.upload-size').innerText = friendlyBytes(e.loaded)
                    tmpl.querySelector('.upload-percent').innerText = percent
                    let pb = tmpl.querySelector('.progress-bar')
                    pb.style['width'] = `${percent}%`
                    pb.setAttribute('aria-valuenow', percent)

                    if (e.loaded == e.total) {

                      hideBtnCancelUpload()
                      showMessage(input.dataset.uploadMsgFileProcessing, 'info')

                    } else {

                      showMessage(input.dataset.uploadMsgFileUploading, 'info')

                    }

                }, false);

                req.addEventListener('readystatechange', evt => {

                  if (req.readyState !== 4) return

                  hideBtnCancelUpload()

                  // file created
                  if (req.status == 200 || req.status == 201) {

                    tmpl.parentElement.removeChild(tmpl)
                    let filesBody = document.getElementById('files-body')
                    filesBody.innerHTML = req.response
                    //trigger htmx and bootstrap on newly added elements
                    htmx.process(filesBody)
                    //htmx.process does not trigger htmx.onload, so repeat here
                    BSN.initCallback()

                  }
                  /*
                   * file too large. Unfortunately this cannot be detected
                   * anymore at server as the error is wrapped inside others.
                   */
                  else if(req.status == 413) {

                    showBtnRemoveUpload()
                    showMessage(input.dataset.uploadMsgFileTooLarge, 'error')
                    hlPgBar(tmpl.querySelector('.progress-bar'), 'error')

                  }
                  // directory has been removed in the meantime
                  else if(req.status == 404) {

                    showBtnRemoveUpload()
                    showMessage(input.dataset.uploadMsgDirNotFound, 'error')
                    hlPgBar(tmpl.querySelector('.progress-bar'), 'error')

                  }
                  // undetermined errors
                  else {

                    showBtnRemoveUpload()
                    showMessage(input.dataset.uploadMsgUnexpected, 'error')
                    hlPgBar(tmpl.querySelector('.progress-bar'), 'error')

                  }

                })

                req.open(form.method, form.action);
                for(let i = 0; i < headers.length; i++) {
                  req.setRequestHeader(headers[i][0], headers[i][1])
                }
                tmpl.querySelector('.btn-cancel-upload').addEventListener('click', function(evt) {
                  evt.preventDefault()
                  req.abort()
                })
                req.send(file);

            })

            // important to retrigger "change" when someone enters the same file again
            input.value = "";
        })
    });
}

const byteUnits = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB']

function friendlyBytes(n) {
    if (n < 10) {
        return n + ' B'
    }
    let e = Math.floor(Math.log(n) / Math.log(1000))
    let unit = byteUnits[e]
    let val =  Math.floor(n / Math.pow(1000, e)*10+0.5) / 10
    if (val < 10 && !Number.isInteger(val)) {
        return val.toFixed(1) + ' ' + unit
    }
    return val.toFixed(0) + ' ' + unit
}
