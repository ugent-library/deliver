import getRandomText from 'support/util'

describe('Managing folders', () => {
  beforeEach(() => {
    cy.loginAsSpaceAdmin()
  })

  it('should be possible to create a new folder', () => {
    const FOLDER_NAME = `CYPRESS-${getRandomText()}`

    cy.visitSpace()

    cy.contains('a', FOLDER_NAME).should('not.exist')

    cy.setFieldByLabel('Folder name', FOLDER_NAME)
    cy.contains('.btn', 'Make folder').click()

    cy.extractFolderId().getFolderShareUrl(FOLDER_NAME).as('shareUrl')

    cy.ensureToast('Folder created successfully').closeToast()
    cy.ensureNoToast()

    cy.get('.bc-toolbar-title').should('contain.text', Cypress.env('DEFAULT_SPACE')).should('contain.text', FOLDER_NAME)

    cy.get('.btn:contains("Copy public shareable link")')
      .as('copyButton')
      .next('input')
      .should('have.value', '@shareUrl')
    cy.getClipboardText().should('not.eq', '@shareUrl')
    cy.get('@copyButton').click().should('contain.text', 'Copied')
    cy.getClipboardText().should('eq', '@shareUrl')

    // Original text resets after 1.5s
    cy.wait(1500)
    cy.get('@copyButton').should('contain.text', 'Copy public shareable link')

    cy.visitSpace()

    cy.contains('tr', FOLDER_NAME)
      .should('exist')
      .find('td')
      .as('folderRow')
      .eq(3)
      .should('contain.text', '0 files')
      .should('contain.text', '0 B')
      .should('contain.text', '0 downloads')

    cy.get('@folderRow').contains('.btn', 'Copy link').as('copyButton').next('input').should('have.value', '@shareUrl')
    cy.setClipboardText('')
    cy.get('@copyButton').click().should('contain.text', 'Copied')
    cy.getClipboardText().should('eq', '@shareUrl')

    // Original text resets after 1.5s
    cy.wait(1500)
    cy.get('@copyButton').should('contain.text', 'Copy link')
  })

  it('should return an error if a new folder name is already in use within the same space', () => {
    const FOLDER_NAME = `CYPRESS-${getRandomText()}`

    cy.visitSpace()

    cy.setFieldByLabel('Folder name', FOLDER_NAME)
    cy.contains('.btn', 'Make folder').click()

    cy.location('pathname').should('match', /\/folders\/\w{26}/)

    cy.visitSpace()

    cy.get('#folder-name').should('not.have.class', 'is-invalid')
    cy.get('#folder-name-invalid').should('not.exist')

    cy.setFieldByLabel('Folder name', FOLDER_NAME)
    cy.contains('.btn', 'Make folder').click()

    cy.get('#folder-name').should('have.class', 'is-invalid')
    cy.get('#folder-name-invalid').should('be.visible').and('have.text', 'name must be unique')
  })

  it('should be possible to edit a folder name', () => {
    const FOLDER_NAME1 = `CYPRESS-FOLDER_NAME_1-${getRandomText()}`
    const FOLDER_NAME2 = `CYPRESS-FOLDER_NAME_2-${getRandomText()}`

    cy.visitSpace()

    cy.setFieldByLabel('Folder name', FOLDER_NAME1)
    cy.contains('.btn', 'Make folder').click()

    cy.extractFolderId('previousFolderId').getFolderShareUrl(FOLDER_NAME1).as('previousShareUrl')
    cy.location('pathname').as('previousPathname')

    cy.get('.bc-toolbar-title').should('contain.text', FOLDER_NAME1)
    cy.contains('Copy public shareable link').next('input').should('have.value', '@previousShareUrl')

    cy.contains('.btn', 'Edit').click()

    cy.setFieldByLabel('Folder name', FOLDER_NAME2)
    cy.contains('.btn', 'Save changes').click()

    cy.extractFolderId(false)
      .should('eq', '@previousFolderId', 'Folder ID must not change after edit')
      .getFolderShareUrl(FOLDER_NAME2)
      .as('newShareUrl')
      .should('not.eq', '@previousShareUrl', 'Share URL should change after edit')

    cy.location('pathname').should('eq', '@previousPathname', 'Pathname should not change after edit')
    cy.get('.bc-toolbar-title').should('contain.text', FOLDER_NAME2, 'Folder title should change after edit')
    cy.contains('Copy public shareable link').next('input').should('have.value', '@newShareUrl')
  })

  it('should return an error if an updated folder name is already in use within the same space', () => {
    const FOLDER_NAME1 = `CYPRESS-FOLDER_NAME_1-${getRandomText()}`
    const FOLDER_NAME2 = `CYPRESS-FOLDER_NAME_2-${getRandomText()}`

    cy.visitSpace()
    cy.setFieldByLabel('Folder name', FOLDER_NAME1)
    cy.contains('.btn', 'Make folder').click()
    cy.location('pathname').as('previousPathname')

    cy.visitSpace()
    cy.setFieldByLabel('Folder name', FOLDER_NAME2)
    cy.contains('.btn', 'Make folder').click()

    cy.contains('.btn', 'Edit').click() // Editing folder 2

    cy.get('#folder-name').should('not.have.class', 'is-invalid')
    cy.get('#folder-name-invalid').should('not.exist')

    cy.setFieldByLabel('Folder name', FOLDER_NAME1)
    cy.contains('.btn', 'Save changes').click()

    cy.get('#folder-name').should('have.class', 'is-invalid')
    cy.get('#folder-name-invalid').should('be.visible').and('have.text', 'name must be unique')

    cy.visit('@previousPathname')

    cy.contains('.btn', 'Edit').click() // Editing folder 1

    cy.get('#folder-name').should('not.have.class', 'is-invalid')
    cy.get('#folder-name-invalid').should('not.exist')

    cy.setFieldByLabel('Folder name', FOLDER_NAME2)
    cy.contains('.btn', 'Save changes').click()

    cy.get('#folder-name').should('have.class', 'is-invalid')
    cy.get('#folder-name-invalid').should('be.visible').and('have.text', 'name must be unique')

    cy.setFieldByLabel('Folder name', FOLDER_NAME1)
    cy.contains('.btn', 'Save changes').click()

    cy.location('pathname').should('eq', '@previousPathname')
  })

  it('should be possible to delete a folder', () => {
    const FOLDER_NAME = `CYPRESS-${getRandomText()}`

    cy.visitSpace()

    cy.contains('a', FOLDER_NAME).should('not.exist')

    cy.setFieldByLabel('Folder name', FOLDER_NAME)
    cy.contains('.btn', 'Make folder').click()

    cy.ensureToast('Folder created successfully').closeToast()
    cy.ensureNoToast()

    cy.visitSpace()

    cy.contains('a', FOLDER_NAME).should('exist')

    cy.contains('a', FOLDER_NAME).click()

    cy.contains('.btn', 'Edit').click()

    // TODO Remove when issue #99 is resolved
    Cypress.on('uncaught:exception', () => {
      // returning false here prevents Cypress from failing the test
      return false
    })

    cy.contains('.btn', 'Delete folder').click()

    cy.location('pathname').should('eq', `/spaces/${Cypress.env('DEFAULT_SPACE')}`)

    cy.ensureToast('Folder deleted successfully').closeToast()
    cy.ensureNoToast()

    cy.contains('a', FOLDER_NAME).should('not.exist')
  })
})
