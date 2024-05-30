import * as esbuild from "esbuild";
import { clean } from "esbuild-plugin-clean";
import { sassPlugin } from "esbuild-sass-plugin";
import manifestPlugin from "esbuild-plugin-manifest";
import fs from "fs";

const config = {
  entryPoints: [
    { in: "assets/js/app.js", out: "js/app" },
    { in: "assets/css/app.scss", out: "css/app" },
    "assets/ugent/images/**/*",
    "assets/ugent/favicon.ico",
    "assets/ugent/fonts/**/*",
  ],
  outdir: "static/",
  bundle: true,
  minify: true,
  sourcemap: true,
  legalComments: "linked",
  loader: {
    ".ico": "copy",
    ".woff": "copy",
    ".woff2": "copy",
    ".svg": "copy",
    ".png": "copy",
  },
  plugins: [
    clean({ patterns: ["./static/*"] }),
    sassPlugin({
      embedded: true,
    }),
    manifestPlugin({
      filter: (fn) => !fn.endsWith(".map") && !fn.endsWith(".LEGAL.txt"),
      generate: generateManifest,
    }),
  ],
};

if (process.argv.includes("--watch")) {
  const ctx = await esbuild.context(config);
  await ctx.watch();

  console.log("ESBuild running. Watching for changes...");
} else {
  console.log(
    "---------------------------------- Building assets ----------------------------------",
  );

  const result = await esbuild.build(config);

  fs.writeFileSync("meta.json", JSON.stringify(result.metafile));

  console.log(
    "-------------------------------------------------------------------------------------",
  );
}

function generateManifest(input) {
  return Object.entries(input).reduce((out, [k, v]) => {
    // Remove "static" from keys
    out[k.replace("static", "")] = `/${v}`;

    return out;
  }, {});
}
