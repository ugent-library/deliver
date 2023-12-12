import getRandomText from 'support/util'

const DEFAULT_SPACE = 'test'

describe('Managing files', () => {
  let FOLDER_NAME: string
  let FILE_COUNT: number

  beforeEach(() => {
    FILE_COUNT = 0
    FOLDER_NAME = `CYPRESS-${getRandomText()}`

    cy.loginAsSpaceAdmin()

    cy.visit(`/spaces/${DEFAULT_SPACE}`)

    cy.setFieldByLabel('Folder name', FOLDER_NAME)
    cy.contains('.btn', 'Make folder').click()
  })

  it('should be possible to upload multiple file types ', () => {
    cy.contains('.card-header', 'Available files').should('contain', '0 items')
    cy.get('#files table').should('not.exist')

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.pdf')
    assertFileUpload('test.pdf', {
      sortOrder: 0,
      md5Checksum: 'f4e486fddb1f3d9d438926f053d53c6a',
      mimeType: 'application/pdf',
    })

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.json')
    assertFileUpload('test.json', {
      sortOrder: 0,
      md5Checksum: '58e0494c51d30eb3494f7c9198986bb9',
      mimeType: 'application/json',
    })

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.txt')
    assertFileUpload('test.txt', {
      sortOrder: 2,
      md5Checksum: '7215ee9c7d9dc229d2921a40e899ec5f',
      mimeType: 'text/plain',
    })

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.docx')
    assertFileUpload('test.docx', {
      sortOrder: 0,
      md5Checksum: 'f694ce9bacf8d1d83e4978e43908f4e8',
      mimeType: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    })

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.xlsx')
    assertFileUpload('test.xlsx', {
      sortOrder: 4,
      md5Checksum: '351aaa090424e8614e5d80efed489e33',
      mimeType: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    })

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.png')
    assertFileUpload('test.png', {
      sortOrder: 3,
      md5Checksum: '6e516cbb6a21ec3b78228829bdec6bd9',
      mimeType: 'image/png',
    })
  })

  it('should be possible to upload multiple files simultaneously')

  it('should display a progress file during larger uploads', () => {
    cy.get('#file-upload-progress').as('uploadProgress').should('not.be.visible')

    cy.intercept('POST', '/folders/*/files', req => {
      req.on('response', res => {
        res.setDelay(500)
      })
    }).as('uploadFile')

    cy.get('input[type=file]').selectFile('cypress/fixtures/large.pdf')

    cy.get('@uploadProgress').should('be.visible').contains('.btn', 'Cancel upload').should('be.visible')
    cy.get('@uploadProgress').contains('0%')
    cy.get('@uploadProgress').contains('100%')
    cy.get('@uploadProgress')
      .find('.progress-bar')
      .should('have.length', 1)
      // Don't use the .have.property chainer here as that will dump the whole CSSStyleDeclaration object in the log
      // Don't use the .have.css chainer here as that would give you the width in pixels
      .prop('style')
      .its('width')
      .should('eq', '100%')
    cy.get('@uploadProgress').contains('Processing your file. Hold on, do not refresh the page.')

    cy.wait('@uploadFile')

    assertFileUpload('large.pdf')
  })

  it('should be possible to cancel an upload')

  it('should be possible to consult the public shareable link anonymously and download all files')

  it('should keep the download count for each file')

  it('should be possible to delete files')

  function assertFileUpload(
    fileName: string,
    {
      sortOrder = 0,
      md5Checksum,
      mimeType,
    }: {
      sortOrder?: number
      md5Checksum?: string
      mimeType?: string
    } = {}
  ) {
    FILE_COUNT++
    cy.contains('.card-header', 'Available files').should('contain', `${FILE_COUNT} items`)

    cy.get('#files table tbody tr')
      .should('have.length', FILE_COUNT)
      .eq(sortOrder)
      .as('testFile')
      .should('contain', fileName)
      .find('p[id^="file-"][id$="-downloads"]')
      .should('have.text', '0')

    if (md5Checksum) {
      cy.get('@testFile').should('contain', `md5 checksum: ${md5Checksum}`)
    }

    if (mimeType) {
      cy.get('@testFile').should('contain', mimeType)
    }
  }
})