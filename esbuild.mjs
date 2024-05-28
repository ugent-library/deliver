import * as esbuild from "esbuild";
import { sassPlugin } from "esbuild-sass-plugin";

const config = {
  entryPoints: [
    "assets/ugent/favicon.ico",
    "assets/ugent/fonts/*",
    "assets/ugent/images/*",
    { in: "assets/js/app.js", out: "js/app" },
    { in: "assets/css/app.scss", out: "css/app" },
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
  plugins: [sassPlugin()],
};

if (process.argv.includes("--watch")) {
  const ctx = await esbuild.context(config);
  await ctx.watch();

  console.log("ESBuild running. Watching for changes...");
} else {
  await esbuild.build(config);
}
