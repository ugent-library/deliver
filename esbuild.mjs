import * as esbuild from "esbuild";
import { sassPlugin } from "esbuild-sass-plugin";

const ctx = await esbuild.context({
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
});

await ctx.watch();
console.log("ESBuild finished. Watching for changes...");
