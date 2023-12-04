// Parent commands
import login from './login'
import loginAsSpaceAdmin from './login-as-space-admin'
import ensureToast from './ensure-toast'
import ensureNoToast from './ensure-no-toast'
import setFieldByLabel from './set-field-by-label'

// Child commands
import finishLog from './finish-log'
import closeToast from './close-toast'
import setField from './set-field'

// Dual commands

// Parent commands
Cypress.Commands.addAll({
  login,

  loginAsSpaceAdmin,

  ensureToast,

  ensureNoToast,

  setFieldByLabel,
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
// Cypress.Commands.addAll(
//   {
//     prevSubject: 'optional',
//   },
//   {
//   }
// )
