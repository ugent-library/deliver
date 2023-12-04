import getRandomText from 'support/util'

const DEFAULT_SPACE = 'test'

describe('Managing folders', () => {
  beforeEach(() => {
    cy.loginAsSpaceAdmin()
  })

  it('should be possible to create a new folder', () => {
    const FOLDER_NAME = `CYPRESS-${getRandomText()}`

    cy.visit(`/spaces/${DEFAULT_SPACE}`)

    cy.contains('a', FOLDER_NAME).should('not.exist')

    cy.setFieldByLabel('Folder name', FOLDER_NAME)

    cy.contains('.btn', 'Make folder').click()

    cy.location('pathname').should('match', /\/folders\/\w{26}/)

    cy.ensureToast('Folder created succesfully') // TODO: update when typo fix deployed
      .closeToast()

    cy.ensureNoToast()

    cy.get('.bc-toolbar-title').should('contain.text', DEFAULT_SPACE).should('contain.text', FOLDER_NAME)

    cy.visit(`/spaces/${DEFAULT_SPACE}`)

    cy.contains('a', FOLDER_NAME).should('exist')
  })

  it('should return an error if a folder name is already in use within the same space')

  it('should be possible to edit a folder name')

  it('should be possible to delete a folder')
})
