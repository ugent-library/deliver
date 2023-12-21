import 'cypress-common'

import './commands'
import './commands/overwrite/should'
import './commands/overwrite/visit'
import './queries'

Cypress.env('DEFAULT_SPACE', 'CYPRESS')
