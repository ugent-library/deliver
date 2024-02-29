import getActivePage from "./get-active-page";
import getLabel from "./get-label";
import getFolderCount from "./get-folder-count";
import getParams from "./get-params";

Cypress.Commands.addQuery("getActivePage", getActivePage);
Cypress.Commands.addQuery("getLabel", getLabel);
Cypress.Commands.addQuery("getFolderCount", getFolderCount);
Cypress.Commands.addQuery("getParams", getParams);
