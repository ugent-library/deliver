import { getRandomText } from "support/util";

describe("Managing folders", () => {
  const NBSP = String.fromCharCode(160);

  beforeEach(() => {
    cy.loginAsSpaceAdmin();
  });

  it("should be possible to create a new folder", () => {
    const FOLDER_NAME = getRandomText();

    cy.visitSpace({ qs: { limit: 1000 } });

    cy.contains("a", FOLDER_NAME).should("not.exist");

    cy.getTotalNumberOfFolders().as("totalNumberOfFolders", { type: "static" });

    cy.setFieldByLabel("Folder name", FOLDER_NAME);
    cy.contains(".btn", "Make folder").click();

    cy.extractFolderId().getFolderShareUrl(FOLDER_NAME).as("shareUrl");

    cy.ensureToast("Folder created successfully").closeToast();
    cy.ensureNoToast();

    cy.get(".bc-toolbar-title")
      .should("contain.text", Cypress.env("DEFAULT_SPACE"))
      .should("contain.text", FOLDER_NAME);

    cy.get('.btn:contains("Copy public shareable link")')
      .as("copyButton")
      .next("input")
      .should("have.value", "@shareUrl");
    cy.getClipboardText().should("not.eq", "@shareUrl");
    cy.get("@copyButton").click().should("contain.text", "Copied");
    cy.getClipboardText().should("eq", "@shareUrl");

    // Original text resets after 1.5s
    cy.wait(1500);
    cy.get("@copyButton").should("contain.text", "Copy public shareable link");

    cy.visitSpace({ qs: { limit: 1000 } });

    cy.get<number>("@totalNumberOfFolders").then((totalNumberOfFolders) => {
      cy.getTotalNumberOfFolders().should("eq", totalNumberOfFolders + 1);
    });

    cy.contains("tr", FOLDER_NAME)
      .should("exist")
      .find("td")
      .as("folderRow")
      .eq(3)
      .should("contain.text", "0 files")
      .should("contain.text", "0 B")
      .should("contain.text", "0 downloads");

    cy.get("@folderRow")
      .contains(".btn", "Copy link")
      .as("copyButton")
      .next("input")
      .should("have.value", "@shareUrl");
    cy.setClipboardText("");
    cy.get("@copyButton").click().should("contain.text", "Copied");
    cy.getClipboardText().should("eq", "@shareUrl");

    // Original text resets after 1.5s
    cy.wait(1500);
    cy.get("@copyButton").should("contain.text", "Copy link");
  });

  it("should trim folder names", () => {
    const FOLDER_NAME = getRandomText();

    cy.visitSpace();

    cy.setFieldByLabel("Folder name", ` \t  ${FOLDER_NAME}  ${NBSP} `);
    cy.contains(".btn", "Make folder").click();

    cy.get(".bc-toolbar-title")
      .should("contain.text", Cypress.env("DEFAULT_SPACE"))
      .should("contain.text", FOLDER_NAME);

    cy.contains(".btn", "Edit").click();

    cy.get("input#folder-name").should("have.value", FOLDER_NAME);
  });

  it("should return an error if a new folder name is empty", () => {
    cy.visitSpace();

    cy.getTotalNumberOfFolders().as("totalNumberOfFolders");

    cy.get("#folder-name").should("not.have.class", "is-invalid");
    cy.get("#folder-name-invalid").should("not.exist");

    cy.setFieldByLabel("Folder name", " ");
    cy.contains(".btn", "Make folder").click();

    cy.get("#folder-name").should("have.class", "is-invalid");
    cy.get("#folder-name-invalid")
      .should("be.visible")
      .and("have.text", "name cannot be empty");

    cy.location("pathname").should(
      "eq",
      `/spaces/${Cypress.env("DEFAULT_SPACE")}/folders`
    );

    cy.getTotalNumberOfFolders().should("eq", "@totalNumberOfFolders");
  });

  it("should return an error if a new folder name is already in use within the same space", () => {
    const FOLDER_NAME = getRandomText();

    cy.makeFolder(FOLDER_NAME);

    cy.location("pathname").should("match", /\/folders\/\w{26}/);

    cy.visitSpace();

    cy.get("#folder-name").should("not.have.class", "is-invalid");
    cy.get("#folder-name-invalid").should("not.exist");

    cy.setFieldByLabel("Folder name", FOLDER_NAME);
    cy.contains(".btn", "Make folder").click();

    cy.get("#folder-name").should("have.class", "is-invalid");
    cy.get("#folder-name-invalid")
      .should("be.visible")
      .and("have.text", "name must be unique");
  });

  it("should be possible to edit a folder name", () => {
    const FOLDER_NAME1 = `FOLDER_NAME_1-${getRandomText()}`;
    const FOLDER_NAME2 = `FOLDER_NAME_2-${getRandomText()}`;

    cy.makeFolder(FOLDER_NAME1);

    cy.extractFolderId("previousFolderId")
      .getFolderShareUrl(FOLDER_NAME1)
      .as("previousShareUrl");
    cy.location("pathname").as("previousPathname");

    cy.get(".bc-toolbar-title").should("contain.text", FOLDER_NAME1);
    cy.contains("Copy public shareable link")
      .next("input")
      .should("have.value", "@previousShareUrl");

    cy.contains(".btn", "Edit").click();

    cy.setFieldByLabel("Folder name", FOLDER_NAME2);
    cy.contains(".btn", "Save changes").click();

    cy.extractFolderId(false)
      .should("eq", "@previousFolderId", "Folder ID must not change after edit")
      .getFolderShareUrl(FOLDER_NAME2)
      .as("newShareUrl")
      .should(
        "not.eq",
        "@previousShareUrl",
        "Share URL should change after edit"
      );

    cy.location("pathname").should(
      "eq",
      "@previousPathname",
      "Pathname should not change after edit"
    );
    cy.get(".bc-toolbar-title").should(
      "contain.text",
      FOLDER_NAME2,
      "Folder title should change after edit"
    );
    cy.contains("Copy public shareable link")
      .next("input")
      .should("have.value", "@newShareUrl");
  });

  it("should trim folder names when editing", () => {
    const FOLDER_NAME = getRandomText();

    cy.makeFolder(FOLDER_NAME);

    cy.extractFolderId();

    cy.contains(".btn", "Edit").click();

    cy.setFieldByLabel(
      "Folder name",
      ` \t   ${FOLDER_NAME} (updated)  ${NBSP} `
    );
    cy.contains(".btn", "Save changes").click();

    cy.get("h4.bc-toolbar-title").should(
      "contain.text",
      FOLDER_NAME + " (updated)"
    );

    cy.contains(".btn", "Edit").click();

    cy.get("input#folder-name").should(
      "have.value",
      FOLDER_NAME + " (updated)"
    );
  });

  it("should return an error if an edited folder name is empty", () => {
    const FOLDER_NAME = getRandomText();

    cy.makeFolder(FOLDER_NAME);

    cy.extractFolderId();

    cy.contains(".btn", "Edit").click();

    cy.get("#folder-name").should("not.have.class", "is-invalid");
    cy.get("#folder-name-invalid").should("not.exist");

    cy.setFieldByLabel("Folder name", " ");
    cy.contains(".btn", "Save changes").click();

    cy.get("#folder-name").should("have.class", "is-invalid");
    cy.get("#folder-name-invalid")
      .should("be.visible")
      .and("have.text", "name cannot be empty");

    // TODO: extract visitFolder command (that works with folder Id alias)
    cy.get<string>("@folderId").then((folderId) =>
      cy.visit(`/folders/${folderId}`)
    );

    cy.get("h4.bc-toolbar-title").should("contain.text", FOLDER_NAME);
  });

  it("should return an error if an updated folder name is already in use within the same space", () => {
    const FOLDER_NAME1 = `FOLDER_NAME_1-${getRandomText()}`;
    const FOLDER_NAME2 = `FOLDER_NAME_2-${getRandomText()}`;

    cy.makeFolder(FOLDER_NAME1);
    cy.location("pathname").as("folder1Url");

    cy.makeFolder(FOLDER_NAME2);

    cy.contains(".btn", "Edit").click(); // Editing folder 2

    cy.get("#folder-name").should("not.have.class", "is-invalid");
    cy.get("#folder-name-invalid").should("not.exist");

    cy.setFieldByLabel("Folder name", FOLDER_NAME1);
    cy.contains(".btn", "Save changes").click();

    cy.get("#folder-name").should("have.class", "is-invalid");
    cy.get("#folder-name-invalid")
      .should("be.visible")
      .and("have.text", "name must be unique");

    cy.visit("@folder1Url");

    cy.contains(".btn", "Edit").click(); // Editing folder 1

    cy.get("#folder-name").should("not.have.class", "is-invalid");
    cy.get("#folder-name-invalid").should("not.exist");

    cy.setFieldByLabel("Folder name", FOLDER_NAME2);
    cy.contains(".btn", "Save changes").click();

    cy.get("#folder-name").should("have.class", "is-invalid");
    cy.get("#folder-name-invalid")
      .should("be.visible")
      .and("have.text", "name must be unique");

    cy.setFieldByLabel("Folder name", FOLDER_NAME1);
    cy.contains(".btn", "Save changes").click();

    cy.location("pathname").should("eq", "@folder1Url");
  });

  it("should be possible to delete a folder", () => {
    const FOLDER_NAME = getRandomText();

    cy.makeFolder(FOLDER_NAME);

    cy.visitSpace({ qs: { limit: 1000 } });

    cy.getTotalNumberOfFolders().as("totalNumberOfFolders", { type: "static" });

    cy.contains("a", FOLDER_NAME).should("exist").click();

    cy.contains(".btn", "Edit").click();

    cy.contains(".btn", "Delete folder").click();

    cy.location("pathname").should(
      "eq",
      `/spaces/${Cypress.env("DEFAULT_SPACE")}`
    );

    cy.get<number>("@totalNumberOfFolders").then((totalNumberOfFolders) => {
      cy.getTotalNumberOfFolders().should("eq", totalNumberOfFolders - 1);
    });

    cy.ensureToast("Folder deleted successfully").closeToast();
    cy.ensureNoToast();

    cy.visitSpace({ qs: { limit: 1000 } });

    cy.contains("a", FOLDER_NAME).should("not.exist");
  });
});
