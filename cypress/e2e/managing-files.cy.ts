import getRandomText from 'support/util'

const DEFAULT_SPACE = 'test'

describe('Managing files', () => {
  let FOLDER_NAME: string

  beforeEach(() => {
    FOLDER_NAME = `CYPRESS-${getRandomText()}`

    cy.loginAsSpaceAdmin()

    cy.visit(`/spaces/${DEFAULT_SPACE}`)

    cy.setFieldByLabel('Folder name', FOLDER_NAME)
    cy.contains('.btn', 'Make folder').click()

    cy.extractFolderId()
  })

  it('should be possible to upload multiple file types', () => {
    cy.contains('.card-header', 'Available files').should('contain', '0 items')
    cy.get('#files table').should('not.exist')

    uploadTestFile('test.pdf', 1, 0, 'f4e486fddb1f3d9d438926f053d53c6a', 'application/pdf')

    uploadTestFile('test.json', 2, 0, '58e0494c51d30eb3494f7c9198986bb9', 'application/json')

    uploadTestFile('test.txt', 3, 2, '7215ee9c7d9dc229d2921a40e899ec5f', 'text/plain')

    uploadTestFile(
      'test.docx',
      4,
      0,
      'f694ce9bacf8d1d83e4978e43908f4e8',
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    )

    uploadTestFile(
      'test.xlsx',
      5,
      4,
      '351aaa090424e8614e5d80efed489e33',
      'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    )

    uploadTestFile('test.png', 6, 3, '6e516cbb6a21ec3b78228829bdec6bd9', 'image/png')
  })

  it('should display a progress file during larger uploads')

  it('should be possible to cancel an upload')

  it('should be possible to consult the public shareable link anonymously and download all files')

  it('should keep the download count for each file')

  it('should be possible to delete files')

  function uploadTestFile(
    fileName: string,
    fileCount: number,
    sortOrder: number,
    md5Checksum: string,
    mimeType: string
  ) {
    cy.get('input[type=file]').selectFile(`cypress/fixtures/${fileName}`)

    cy.contains('.card-header', 'Available files').should('contain', `${fileCount} items`)

    cy.get('#files table tbody tr')
      .should('have.length', fileCount)
      .eq(sortOrder)
      .should('contain', fileName)
      .should('contain', `md5 checksum: ${md5Checksum}`)
      .should('contain', mimeType)
      .find('p[id^="file-"][id$="-downloads"]')
      .should('have.text', '0')
  }
})
