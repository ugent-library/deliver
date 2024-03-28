import { defineConfig } from "cypress";
import * as path from "path";
import * as fs from "fs";

export default defineConfig({
  projectId: "fm8w1c",

  env: {
    OIDC_ORIGIN: "http://localhost:3102",
    SUPER_ADMIN_USER_NAME: "superdeliver",
    SPACE_ADMIN_USER_NAME: "deliver",
    DEFAULT_SPACE: "CYPRESS",
  },

  e2e: {
    baseUrl: "http://localhost:3101/",

    experimentalStudio: true,
    experimentalRunAllSpecs: true,

    // Increase viewport width because GitHub Actions may render a wider font which
    // may cause button clicks to be prevented by overlaying elements.
    viewportWidth: 1200,

    setupNodeEvents(on, config) {
      on("task", {
        clearDownloads() {
          const downloadsFolder = path.resolve(config.downloadsFolder);

          if (!fs.existsSync(downloadsFolder)) {
            fs.mkdirSync(downloadsFolder, { recursive: true });
          }

          const downloads = fs.readdirSync(downloadsFolder);

          downloads.forEach((file) => {
            fs.unlinkSync(path.join(downloadsFolder, file));
          });

          return downloads;
        },
      });

      return config;
    },
  },
});
