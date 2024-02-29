import getLabel from "./get-label";
import getNumberOfDisplayedFolders from "./get-number-of-displayed-folders";
import getTotalNumberOfFolders from "./get-total-number-of-folders";
import getParams from "./get-params";

Cypress.Commands.addQuery("getLabel", getLabel);
Cypress.Commands.addQuery(
  "getNumberOfDisplayedFolders",
  getNumberOfDisplayedFolders
);
Cypress.Commands.addQuery("getTotalNumberOfFolders", getTotalNumberOfFolders);
Cypress.Commands.addQuery("getParams", getParams);
