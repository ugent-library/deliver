// Parent commands
import login from './login'
import loginAsSpaceAdmin from './login-as-space-admin'

// Child commands
import finishLog from './finish-log'

// Dual commands

// Parent commands
Cypress.Commands.addAll({
  login,

  loginAsSpaceAdmin,
})

// Child commands
Cypress.Commands.addAll(
  { prevSubject: true },
  {
    finishLog,
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
