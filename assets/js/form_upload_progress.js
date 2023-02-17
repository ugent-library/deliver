export default function() {
    document.querySelectorAll('form input[data-upload-progress-target]').forEach(input => {
        /*
         *  TODO: after full upload provide state to show that server is still processing content,
         *  otherwise you will end up with with a progress bar showing 100% for a long time
         *
         *  TODO: check traefik request timeout when proxying content
         *
         *  TODO: inspect at server if the request body is really streamed into the S3 storage,
         *  because there is some lag after upload, which suggests that the upload to the S3 storage
         *  still needs to start at that point. Maybe because golang middlewares CSRF and ProxyMethodOverride
         *  have started reading the request body, and therefore the request is blocked in those
         *  middlewares before arriving at the destination request handler? That would mean that
         *  the upload ends while we are still in the middleware, and that the request handler
         *  gets a "temporary file" that needs to be uploaded to S3
         *
         *  TODO: even if at server side the middlewars CSRF and ProxyMethodOverride use only the headers,
         *  the request handler cannot return a response before the full request body is sent.
         *  http.MaxBytesReader solves this by closing the connection, but can we do this here too?
         *
         *  TODO: prematurely stop upload on client side
         *
         *  TODO: on server ParseForm is executed, which reads file bigger than 32MB into tmp files.
         *        so .. why not upload file in body only, and pass content-type and file-name via the headers?
         * */
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

                let maxSize = input.dataset.uploadMaxSize

                if (!isNaN(maxSize) && file.size > maxSize) {

                  hideBtnCancelUpload()
                  showMessage(input.dataset.uploadMsgFileTooLarge, 'error')
                  hlPgBar(tmpl.querySelector('.progress-bar'), 'error')
                  return

                }

                let formData = new FormData()
                formData.append('file', file)

                let headers = [
                  ['X-CSRF-Token', csrfToken],
                  // weird, but makes sure that middleware does not try to read _method from form
                  ['X-HTTP-Method-Override', 'POST']
                ]

                let req = new XMLHttpRequest();

                req.addEventListener('abort', e => {
                  showMessage(input.dataset.uploadMsgFileAborted, 'error')
                  hlPgBar(tmpl.querySelector('.progress-bar'), 'error')
                  hideBtnCancelUpload()
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
                      pb.classList.add('progress-bar-striped')

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
                    //trigger htmx on newly added elements
                    htmx.process(filesBody)

                  }
                  // file too large
                  else if(req.status == 413) {

                    showMessage(input.dataset.uploadMsgFileTooLarge, 'error')
                    hlPgBar(tmpl.querySelector('.progress-bar'), 'error')

                  }
                  else if(req.status == 404) {

                    showMessage(input.dataset.uploadMsgDirNotFound, 'error')
                    hlPgBar(tmpl.querySelector('.progress-bar'), 'error')

                  }
                  // undetermined errors
                  else {

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
                req.send(formData);

            })
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
