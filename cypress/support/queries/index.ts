import getLabel from './get-label'
import getTotalNumberOfFolders from './get-total-number-of-folders'

Cypress.Commands.addQuery('getLabel', getLabel)
Cypress.Commands.addQuery('getTotalNumberOfFolders', getTotalNumberOfFolders)
