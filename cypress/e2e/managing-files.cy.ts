import getRandomText from 'support/util'

const DEFAULT_SPACE = 'test'

describe('Managing files', () => {
  let FOLDER_NAME: string
  let FILE_COUNT: number

  beforeEach(() => {
    cy.task('clearDownloads')

    FILE_COUNT = 0
    FOLDER_NAME = `CYPRESS-${getRandomText()}`

    cy.loginAsSpaceAdmin()

    cy.visit(`/spaces/${DEFAULT_SPACE}`)

    cy.setFieldByLabel('Folder name', FOLDER_NAME)
    cy.contains('.btn', 'Make folder').click()

    cy.ensureToast().closeToast()

    cy.url().as('adminUrl')
  })

  it('should be possible to upload multiple file types ', () => {
    cy.contains('.card-header', 'Available files').should('contain', '0 items')
    cy.get('#files table').should('not.exist')

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.pdf', { action: 'select' })
    assertFileUpload('test.pdf', {
      sortOrder: 0,
      md5Checksum: 'f4e486fddb1f3d9d438926f053d53c6a',
      mimeType: 'application/pdf',
    })

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.json', { action: 'drag-drop' })
    assertFileUpload('test.json', {
      sortOrder: 0,
      md5Checksum: '58e0494c51d30eb3494f7c9198986bb9',
      mimeType: 'application/json',
    })

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.txt', { action: 'select' })
    assertFileUpload('test.txt', {
      sortOrder: 2,
      md5Checksum: '7215ee9c7d9dc229d2921a40e899ec5f',
      mimeType: 'text/plain',
    })

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.docx', { action: 'drag-drop' })
    assertFileUpload('test.docx', {
      sortOrder: 0,
      md5Checksum: 'f694ce9bacf8d1d83e4978e43908f4e8',
      mimeType: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    })

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.xlsx', { action: 'select' })
    assertFileUpload('test.xlsx', {
      sortOrder: 4,
      md5Checksum: '351aaa090424e8614e5d80efed489e33',
      mimeType: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    })

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.png', { action: 'drag-drop' })
    assertFileUpload('test.png', {
      sortOrder: 3,
      md5Checksum: '6e516cbb6a21ec3b78228829bdec6bd9',
      mimeType: 'image/png',
    })
  })

  it('should be possible to upload multiple files simultaneously', () => {
    cy.contains('.card-header', 'Available files').should('contain', '0 items')
    cy.get('#files table').should('not.exist')
    cy.get('#file-upload-progress').as('uploadProgress').should('not.be.visible')

    cy.intercept('POST', '/folders/*/files', req => {
      req.on('response', res => {
        // Cause an artificial delay in the response so all assertions have time to succeed during uploading
        res.setDelay(1000)
      })
    }).as('uploadFiles')

    cy.get('input[type=file]').selectFile([
      generateLargeFile('large.txt', 1),
      'cypress/fixtures/test.pdf',
      'cypress/fixtures/test.txt',
      'cypress/fixtures/test.json',
    ])

    cy.get('@uploadProgress')
      .should('be.visible')
      .get('.btn:contains("Cancel upload")')
      .should('have.length', 4)
      .and('be.visible')
    cy.get('@uploadProgress')
      .should('contain', 'large.txt')
      .and('contain', 'test.pdf')
      .and('contain', 'test.txt')
      .and('contain', 'test.json')

    cy.wait('@uploadFiles')
    cy.wait('@uploadFiles')
    cy.wait('@uploadFiles')
    cy.wait('@uploadFiles')

    cy.contains('.card-header', 'Available files').should('contain', '4 items')
    cy.get('#files table')
      .should('be.visible')
      .and('contain', 'large.txt')
      .and('contain', 'test.pdf')
      .and('contain', 'test.txt')
      .and('contain', 'test.json')
  })

  it('should display a progress file during uploads', () => {
    cy.get('#file-upload-progress').as('uploadProgress').should('not.be.visible')

    cy.intercept('POST', '/folders/*/files').as('uploadFile')

    const fileName = 'large-file.pdf'
    cy.get('input[type=file]').selectFile(generateLargeFile(fileName, 1))

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

    assertFileUpload(fileName)
  })

  it('should be possible to cancel an upload', () => {
    const fileName = 'very-large-file.txt'
    cy.get('input[type=file]').selectFile(generateLargeFile(fileName, 10))

    cy.wait(50)

    cy.get('#file-upload-progress').as('uploadProgress').should('be.visible').contains('.btn', 'Cancel upload').click()

    cy.get('@uploadProgress').contains('.btn', 'Cancel upload').should('not.exist')

    cy.get('@uploadProgress').contains('File upload aborted by you').should('be.visible')

    cy.get('@uploadProgress').contains('.btn', 'Remove').click()

    cy.get('@uploadProgress').should('not.be.visible')

    cy.contains(fileName).should('not.exist')

    // Wait some time to make sure the file is not still being processed on the server
    cy.wait(2000)

    cy.reload()

    cy.contains(fileName).should('not.exist')
    cy.contains('.card-header', 'Available files').should('contain', '0 items')
    cy.get('#files table tbody tr').should('have.length', 0)
  })

  it('should be possible to consult the public shareable link anonymously, download all files and keep download counts', () => {
    cy.extractFolderId().getFolderShareUrl(FOLDER_NAME).as('shareUrl')

    cy.intercept('POST', '/folders/*/files').as('uploadFile')

    cy.get('input[type=file]').selectFile('cypress/fixtures/test.txt')
    cy.get('input[type=file]').selectFile('cypress/fixtures/test.pdf')
    cy.get('input[type=file]').selectFile('cypress/fixtures/test.json')

    cy.wait('@uploadFile')

    assertNumberOfDownloads('test.txt', 0)
    assertNumberOfDownloads('test.pdf', 0)
    assertNumberOfDownloads('test.json', 0)

    assertTotalNumberOfDownloads(0)

    cy.logout()

    cy.visit('@shareUrl')

    cy.contains(`Library delivery from ${DEFAULT_SPACE}: ${FOLDER_NAME}`)
      .should('be.visible')
      .then(function () {
        cy.contains(this.shareUrl).should('be.visible')
        cy.contains(this.adminUrl).should('not.exist')

        assertNumberOfDownloads('test.txt', 0)
        assertNumberOfDownloads('test.pdf', 0)
        assertNumberOfDownloads('test.json', 0)

        const zipFileName = `cypress/downloads/${this.folderId}-${FOLDER_NAME}.zip`
        cy.readFile(zipFileName).should('not.exist')
        cy.contains('.btn', 'Download all files').click()
        cy.readFile(zipFileName)

        cy.reload()
        assertNumberOfDownloads('test.txt', 1)
        assertNumberOfDownloads('test.pdf', 1)
        assertNumberOfDownloads('test.json', 1)

        // Now test each individual download using the file name link in the first column
        cy.readFile('cypress/downloads/test.txt').should('not.exist')
        cy.contains('test.txt').click()
        cy.readFile('cypress/downloads/test.txt')

        cy.reload()
        assertNumberOfDownloads('test.txt', 2)
        assertNumberOfDownloads('test.pdf', 1)
        assertNumberOfDownloads('test.json', 1)

        cy.readFile('cypress/downloads/test.pdf').should('not.exist')
        cy.contains('test.pdf').click()
        cy.readFile('cypress/downloads/test.pdf')

        cy.reload()
        assertNumberOfDownloads('test.txt', 2)
        assertNumberOfDownloads('test.pdf', 2)
        assertNumberOfDownloads('test.json', 1)

        cy.readFile('cypress/downloads/test.json').should('not.exist')
        cy.contains('test.json').click()
        cy.readFile('cypress/downloads/test.json')

        cy.reload()
        assertNumberOfDownloads('test.txt', 2)
        assertNumberOfDownloads('test.pdf', 2)
        assertNumberOfDownloads('test.json', 2)

        cy.task('clearDownloads')

        // Now test the same using the "Download" links in the last column
        cy.readFile('cypress/downloads/test.txt').should('not.exist')
        cy.contains('table tbody tr', 'test.tx').find('td').last().contains('a', 'Download').click()
        cy.readFile('cypress/downloads/test.txt')

        cy.reload()
        assertNumberOfDownloads('test.txt', 3)
        assertNumberOfDownloads('test.pdf', 2)
        assertNumberOfDownloads('test.json', 2)

        cy.readFile('cypress/downloads/test.pdf').should('not.exist')
        cy.contains('table tbody tr', 'test.pdf').find('td').last().contains('a', 'Download').click()
        cy.readFile('cypress/downloads/test.pdf')

        cy.reload()
        assertNumberOfDownloads('test.txt', 3)
        assertNumberOfDownloads('test.pdf', 3)
        assertNumberOfDownloads('test.json', 2)

        cy.readFile('cypress/downloads/test.json').should('not.exist')
        cy.contains('table tbody tr', 'test.json').find('td').last().contains('a', 'Download').click()
        cy.readFile('cypress/downloads/test.json')

        cy.reload()
        assertNumberOfDownloads('test.txt', 3)
        assertNumberOfDownloads('test.pdf', 3)
        assertNumberOfDownloads('test.json', 3)
      })

    cy.task('clearDownloads')

    cy.loginAsSpaceAdmin()

    assertTotalNumberOfDownloads(9)

    assertNumberOfDownloads('test.txt', 3)
    assertNumberOfDownloads('test.pdf', 3)
    assertNumberOfDownloads('test.json', 3)

    cy.readFile('cypress/downloads/test.txt').should('not.exist')
    cy.contains('test.txt').click()
    cy.readFile('cypress/downloads/test.txt')

    cy.reload()
    assertNumberOfDownloads('test.txt', 4)
    assertNumberOfDownloads('test.pdf', 3)
    assertNumberOfDownloads('test.json', 3)

    assertTotalNumberOfDownloads(10)

    cy.readFile('cypress/downloads/test.pdf').should('not.exist')
    cy.contains('test.pdf').click()
    cy.readFile('cypress/downloads/test.pdf')

    cy.reload()
    assertNumberOfDownloads('test.txt', 4)
    assertNumberOfDownloads('test.pdf', 4)
    assertNumberOfDownloads('test.json', 3)

    assertTotalNumberOfDownloads(11)

    cy.readFile('cypress/downloads/test.json').should('not.exist')
    cy.contains('test.json').click()
    cy.readFile('cypress/downloads/test.json')

    cy.reload()
    assertNumberOfDownloads('test.txt', 4)
    assertNumberOfDownloads('test.pdf', 4)
    assertNumberOfDownloads('test.json', 4)

    assertTotalNumberOfDownloads(12)
  })

  it('should be possible to delete files', () => {
    cy.contains('.card-header', 'Available files').should('contain', '0 items')

    cy.intercept('POST', '/folders/*/files').as('uploadFile')

    cy.get('input[type=file]').selectFile([
      'cypress/fixtures/test.json',
      'cypress/fixtures/test.pdf',
      'cypress/fixtures/test.txt',
    ])

    cy.wait('@uploadFile')
    cy.wait('@uploadFile')
    cy.wait('@uploadFile')

    assertFolderFileCount(3)

    assertFileDelete(3, 'test.pdf')

    assertFolderFileCount(2)

    assertFileDelete(2, 'test.txt')

    assertFolderFileCount(1)

    assertFileDelete(1, 'test.json')

    assertFolderFileCount(0)
  })

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

    assertFolderFileCount(FILE_COUNT)
  }

  function generateLargeFile(fileName: string, fileSizeInMegaByte: number, mimeType?: string) {
    const largeString = 'a'.repeat(fileSizeInMegaByte * 1024 * 1024) // 5MB
    const buffer = Cypress.Buffer.from(largeString)

    const file: Cypress.FileReferenceObject = {
      fileName,
      contents: buffer,
      mimeType,
    }

    return file
  }

  function assertNumberOfDownloads(fileName: string, expectedNumberOfDownloads: number) {
    cy.contains('table tbody tr', fileName).find('td').eq(3).should('have.text', expectedNumberOfDownloads)
  }

  function assertTotalNumberOfDownloads(expectedNumberOfDownloads: number) {
    cy.visit(`/spaces/${DEFAULT_SPACE}`)

    cy.contains('table tbody tr', FOLDER_NAME)
      .find('td')
      .eq(3)
      .should('contain', `${expectedNumberOfDownloads} downloads`)

    cy.visit('@adminUrl')
  }

  function assertFileDelete(numberOfAvailableFilesAtStart: number, fileToDelete: string) {
    cy.contains('.card-header', 'Available files').should('contain', `${numberOfAvailableFilesAtStart} items`)

    cy.ensureNoModal()

    cy.contains('#files table tr', fileToDelete).contains('.btn', 'Delete').click()

    cy.ensureModal(new RegExp(`Are you sure you want to delete the file.*${fileToDelete}\?`)).closeModal(
      'Yes, delete this file'
    )

    cy.ensureNoModal()

    cy.contains('.card-header', 'Available files').should('contain', `${numberOfAvailableFilesAtStart - 1} items`)
    cy.contains('#files table tr', fileToDelete).should('not.exist')
  }

  function assertFolderFileCount(FILE_COUNT: number) {
    cy.visit(`/spaces/${DEFAULT_SPACE}`)

    cy.contains('table tr', FOLDER_NAME).find('td').eq(3).should('contain', `${FILE_COUNT} files`)

    cy.visit('@adminUrl')
  }
})
