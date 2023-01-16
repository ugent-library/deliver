export default function() {
    document.querySelectorAll("form input[data-upload-progress-target]").forEach(input => {
        input.addEventListener("change", () => {
            let files = Array.from(input.files)

            if (!files.length) {
                return
            }

            let target = document.getElementById(input.dataset.uploadProgressTarget)
            let form = input.closest("form")
            let numUploaded = 0;

            form.style.display = 'none';

            files.forEach((file, i) => {
                let tmpl = document.getElementById("tmpl-upload-progress").content.firstElementChild.cloneNode(true)
                tmpl.querySelector('.upload-name').innerText = file.name
                target.appendChild(tmpl)

                let formData = new FormData()
                formData.append('file', file)
                if (input.dataset.uploadProgressInclude) {
                    input.dataset.uploadProgressInclude.split(",").forEach((name) => {
                        formData.append(name, form.querySelector(`[name="${name}"]`).value)
                    })
                }
    
                let ajax = new XMLHttpRequest();

                ajax.addEventListener('error', e => {
                    console.log(e)
                }, false);

                ajax.upload.addEventListener('progress', e => {
                    let percent = Math.ceil(e.loaded / e.total) * 100
                    console.log(file.name + ": " + percent)
                    tmpl.querySelector('.upload-size').innerText = friendlyBytes(e.loaded)
                    tmpl.querySelector('.upload-percent').innerText = percent
                    let pb = tmpl.querySelector('.progress-bar')
                    pb.style['width'] = `${percent}%`
                    pb.setAttribute('aria-valuenow', percent)
                }, false);

                ajax.addEventListener('loadend', () => {
                    tmpl.style.display = 'none';
                    numUploaded++;
                    if (numUploaded === files.length) {
                        window.location.reload()
                    }
                }, false);

                ajax.open(form.method, form.action);
                ajax.send(formData);
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
    if (val < 10) {
        return val.toFixed(1) + ' ' + unit
    }
    return val.toFixed(0) + ' ' + unit 
}
