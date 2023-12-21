// Parent commands
import login from './login'
import loginAsSpaceAdmin from './login-as-space-admin'
import ensureToast from './ensure-toast'
import ensureNoToast from './ensure-no-toast'
import setFieldByLabel from './set-field-by-label'
import getClipboardText from './get-clipboard-text'
import setClipboardText from './set-clipboard-text'
import extractFolderId from './extract-folder-id'
import logout from './logout'
import ensureModal from './ensure-modal'
import ensureNoModal from './ensure-no-modal'

// Child commands
import finishLog from './finish-log'
import closeToast from './close-toast'
import setField from './set-field'

// Dual commands
import getFolderShareUrl from './get-folder-share-url'
import closeModal from './close-modal'

// Parent commands
Cypress.Commands.addAll({
  login,

  loginAsSpaceAdmin,

  ensureToast,

  ensureNoToast,

  setFieldByLabel,

  getClipboardText,

  setClipboardText,

  extractFolderId,

  logout,

  ensureModal,

  ensureNoModal,
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

    closeModal,
  }
)
