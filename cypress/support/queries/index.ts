import getLabel from "./get-label";
import getFolderCount from "./get-folder-count";

Cypress.Commands.addQuery("getLabel", getLabel);
Cypress.Commands.addQuery("getFolderCount", getFolderCount);
