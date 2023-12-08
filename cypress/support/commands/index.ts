// Parent commands
import login from './login'
import loginAsSpaceAdmin from './login-as-space-admin'
import ensureToast from './ensure-toast'
import ensureNoToast from './ensure-no-toast'
import setFieldByLabel from './set-field-by-label'
import getClipboardText from './get-clipboard-text'
import extractFolderId from './extract-folder-id'

// Child commands
import finishLog from './finish-log'
import closeToast from './close-toast'
import setField from './set-field'

// Dual commands
import getFolderShareUrl from './get-folder-share-url'

// Parent commands
Cypress.Commands.addAll({
  login,

  loginAsSpaceAdmin,

  ensureToast,

  ensureNoToast,

  setFieldByLabel,

  getClipboardText,

  extractFolderId,
})

// Child commands
Cypress.Commands.addAll(
  { prevSubject: true },
  {
    finishLog,

    closeToast,

    setField,
  }
)

// Dual commands
Cypress.Commands.addAll(
  {
    prevSubject: 'optional',
  },
  {
    getFolderShareUrl,
  }
)
