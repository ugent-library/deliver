// Parent commands
import login from "./login";
import loginAsSpaceAdmin from "./login-as-space-admin";
import loginAsSuperAdmin from "./login-as-super-admin";
import logout from "./logout";
import ensureToast from "./ensure-toast";
import ensureNoToast from "./ensure-no-toast";
import setFieldByLabel from "./set-field-by-label";
import getClipboardText from "./get-clipboard-text";
import setClipboardText from "./set-clipboard-text";
import extractFolderId from "./extract-folder-id";
import ensureModal from "./ensure-modal";
import ensureNoModal from "./ensure-no-modal";
import visitSpace from "./visit-space";
import makeFolder from "./make-folder";
import cleanUp from "./clean-up";
import disableHtmx from "./disable-htmx";

// Child commands
import finishLog from "./finish-log";
import closeToast from "./close-toast";
import setField from "./set-field";
import submitForm from "./submit-form";

// Dual commands
import getFolderShareUrl from "./get-folder-share-url";
import closeModal from "./close-modal";

// Parent commands
Cypress.Commands.addAll({
  login,

  loginAsSpaceAdmin,

  loginAsSuperAdmin,

  logout,

  ensureToast,

  ensureNoToast,

  setFieldByLabel,

  getClipboardText,

  setClipboardText,

  extractFolderId,

  ensureModal,

  ensureNoModal,

  visitSpace,

  makeFolder,

  cleanUp,

  disableHtmx,
});

// Child commands
Cypress.Commands.addAll(
  { prevSubject: true },
  {
    finishLog,

    closeToast,

    setField,

    submitForm,
  }
);

// Dual commands
Cypress.Commands.addAll(
  {
    prevSubject: "optional",
  },
  {
    getFolderShareUrl,

    closeModal,
  }
);
