export default function (rootEl) {
  const requests = [];
  let onBeforeUnloadListenerAdded = false;

  rootEl
    .querySelectorAll("form input[data-upload-progress-target]")
    .forEach((input) => {
      input.addEventListener("change", () => {
        const files = Array.from(input.files);

        if (!files.length) return;

        const {
          uploadProgressTarget,
          uploadMaxFileSize,
          uploadMsgFileTooLarge,
          uploadMsgFileAborted,
          uploadMsgFileUploading,
          uploadMsgFileProcessing,
          uploadMsgDirNotFound,
          uploadMsgUnexpected,
        } = input.dataset;

        const target = document.getElementById(uploadProgressTarget);
        const form = input.closest("form");

        files.forEach((file) => {
          const tmpl = new FileTemplate(target, file.name);

          // prevent file upload when above max file size
          if (!isNaN(uploadMaxFileSize) && file.size > uploadMaxFileSize) {
            tmpl.showRemoveUploadButton();
            tmpl.showMessage(uploadMsgFileTooLarge, "error");
            return;
          }

          if (!onBeforeUnloadListenerAdded) {
            window.addEventListener("beforeunload", (evt) => {
              // Cancelled/aborted requests have readyState UNSENT (0)
              if (
                requests.some(
                  (r) =>
                    r.readyState != XMLHttpRequest.UNSENT &&
                    r.readyState != XMLHttpRequest.DONE,
                )
              ) {
                evt.preventDefault();
                evt.returnValue =
                  "One or more file uploads are still in progress. Are you sure you want to leave this page? Your upload(s) will be cancelled.";
              }
            });

            onBeforeUnloadListenerAdded = true;
          }

          const req = new XMLHttpRequest();
          requests.push(req);

          req.addEventListener("abort", () => {
            tmpl.showRemoveUploadButton();
            tmpl.showMessage(uploadMsgFileAborted, "error");
          });

          req.upload.addEventListener("progress", (e) => {
            if (!e.lengthComputable) return;

            tmpl.uploadSize = friendlyBytes(e.loaded);
            tmpl.uploadPercentage = Math.floor((e.loaded / e.total) * 100);

            if (e.loaded == e.total) {
              tmpl.showMessage(uploadMsgFileProcessing, "info");
            } else {
              tmpl.showMessage(uploadMsgFileUploading, "info");
            }
          });

          req.addEventListener("readystatechange", () => {
            if (req.readyState !== XMLHttpRequest.DONE) return;

            switch (req.status) {
              case 200:
              case 201:
                // file created
                tmpl.destroy();
                htmx.trigger("body", "refresh-files");
                break;

              case 413:
                // File too large. Unfortunately this cannot be detected
                // anymore at server as the error is wrapped inside others.
                tmpl.showRemoveUploadButton();
                tmpl.showMessage(uploadMsgFileTooLarge, "error");
                break;

              case 404:
                // directory has been removed in the meantime
                tmpl.showRemoveUploadButton();
                tmpl.showMessage(uploadMsgDirNotFound, "error");
                break;

              default:
                // undetermined errors
                tmpl.showRemoveUploadButton();
                tmpl.showMessage(uploadMsgUnexpected, "error");
            }
          });

          req.open(form.method, form.action);

          const headers = getHttpHeaders(file);
          for (const key in headers) {
            req.setRequestHeader(key, headers[key]);
          }

          tmpl.onCancel((evt) => {
            evt.preventDefault();
            req.abort();
          });

          req.send(file);
        });

        // important to retrigger "change" when someone enters the same file again
        input.value = "";
      });
    });
}

class FileTemplate {
  #template = null;
  #qsCache = new Map();

  constructor(target, fileName) {
    this.#template = document
      .getElementById("tmpl-upload-progress")
      .content.firstElementChild.cloneNode(true);
    target.appendChild(this.#template);

    this.#qs(".upload-name").innerText = fileName;

    this.#qs(".btn-remove-upload").addEventListener(
      "click",
      // Make sure the event handler can reach the class instance via this
      this.destroy.bind(this),
    );
  }

  set uploadSize(size) {
    this.#qs(".upload-size").innerText = size;
  }

  set uploadPercentage(percentage) {
    this.#qs(".upload-percent").innerText = percentage;

    const pb = this.#qs(".progress-bar");
    pb.style["width"] = `${percentage}%`;
    pb.setAttribute("aria-valuenow", percentage);
  }

  showMessage(msg, level) {
    const uploadMessage = this.#qs(".upload-msg");
    const progressBar = this.#qs(".progress-bar");

    uploadMessage.innerText = msg;

    uploadMessage.classList.remove("text-muted");
    uploadMessage.classList.remove("text-danger");
    progressBar.classList.remove("bg-info");
    progressBar.classList.remove("bg-danger");

    if (level == "error") {
      uploadMessage.classList.add("text-danger");
      progressBar.classList.add("bg-danger");
    } else {
      uploadMessage.classList.add("text-muted");
      progressBar.classList.add("bg-info");
    }
  }

  onCancel(callback) {
    this.#qs(".btn-cancel-upload").addEventListener("click", callback);
  }

  showRemoveUploadButton() {
    this.#remove(".btn-cancel-upload");

    const removeButton = this.#qs(".btn-remove-upload");
    if (removeButton) {
      removeButton.classList.remove("d-none");
    }
  }

  destroy() {
    if (this.#template && this.#template.parentElement) {
      this.#template.parentElement.removeChild(this.#template);
    }

    this.#template = null;
  }

  #qs(selector) {
    if (!this.#qsCache.has(selector)) {
      this.#qsCache.set(selector, this.#template.querySelector(selector));
    }

    return this.#qsCache.get(selector);
  }

  #remove(selectorOrElement) {
    if (typeof selectorOrElement == "string") {
      selectorOrElement = this.#qs(selectorOrElement);
    }

    if (selectorOrElement && selectorOrElement.parentElement) {
      selectorOrElement.parentElement.removeChild(selectorOrElement);
    }
  }
}

function getHttpHeaders(file) {
  return {
    "X-CSRF-Token": getCSRFToken(),
    // weird, but makes sure that middleware does not try to read _method from form
    "X-HTTP-Method-Override": "POST",
    //"Failed to execute 'setRequestHeader' on 'XMLHttpRequest': String contains non ISO-8859-1 code point"
    "X-Upload-Filename": encodeURIComponent(file.name),
    // refused by browser
    // "Content-Length": file.size,
    "Content-Type": file.type,
  };
}

function getCSRFToken() {
  return document.querySelector("meta[name=csrf-token]").content;
}

function friendlyBytes(n) {
  const byteUnits = ["B", "KB", "MB", "GB", "TB", "PB", "EB"];

  if (n < 10) {
    return n + " B";
  }

  const e = Math.floor(Math.log(n) / Math.log(1000));
  const unit = byteUnits[e];
  const val = Math.floor((n / Math.pow(1000, e)) * 10 + 0.5) / 10;
  if (val < 10 && !Number.isInteger(val)) {
    return val.toFixed(1) + " " + unit;
  }

  return val.toFixed(0) + " " + unit;
}
