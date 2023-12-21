import { defineConfig } from 'cypress'
import * as dotenvPlugin from 'cypress-dotenv'
import * as path from 'path'
import * as fs from 'fs'

export default defineConfig({
  projectId: 'fm8w1c',

  e2e: {
    baseUrl: 'https://deliver.libtest.ugent.be/',
    experimentalStudio: true,
    experimentalRunAllSpecs: true,

    // Increase viewport width because GitHub Actions may render a wider font which
    // may cause button clicks to be prevented by overlaying elements.
    viewportWidth: 1200,

    setupNodeEvents(on, config) {
      on('task', {
        clearDownloads() {
          const downloadsFolder = path.resolve(config.downloadsFolder)

          if (!fs.existsSync(downloadsFolder)) {
            fs.mkdirSync(downloadsFolder, { recursive: true })
          }

          const downloads = fs.readdirSync(downloadsFolder)

          downloads.forEach(file => {
            fs.unlinkSync(path.join(downloadsFolder, file))
          })

          return downloads
        },
      })

      config = dotenvPlugin(config)

      return config
    },
  },
})
