import getLabel from "./get-label";
import getFolderCount from "./get-folder-count";
import getParams from "./get-params";

Cypress.Commands.addQuery("getLabel", getLabel);
Cypress.Commands.addQuery("getFolderCount", getFolderCount);
Cypress.Commands.addQuery("getParams", getParams);
