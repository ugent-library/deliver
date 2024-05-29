import * as esbuild from "esbuild";
import { sassPlugin } from "esbuild-sass-plugin";
import manifestPlugin from "esbuild-plugin-manifest";
import { readFileSync } from "fs";
import { createHash } from "crypto";

const config = {
  entryPoints: [
    { in: "assets/js/app.js", out: "js/app" },
    { in: "assets/css/app.scss", out: "css/app" },
    "assets/ugent/images/*",
    "assets/ugent/favicon.ico",
    "assets/ugent/fonts/*",
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
  external: ["*.woff", "*.woff2"],
  plugins: [
    sassPlugin({
      embedded: true,
    }),
    manifestPlugin({
      hash: false,
      filter: (fn) => !fn.endsWith(".map") && !fn.endsWith(".LEGAL.txt"),
      generate: generateManifest,
    }),
  ],
};

console.log(
  "------------------------------ (Re)building assets ------------------------------",
);
if (process.argv.includes("--watch")) {
  const ctx = await esbuild.context(config);
  await ctx.watch();

  console.log("ESBuild running. Watching for changes...");
} else {
  await esbuild.build(config);
}
console.log(
  "---------------------------------------------------------------------------------",
);

function generateManifest(input) {
  return Object.entries(input).reduce((out, [k, v]) => {
    // Create md5 hash from file
    const contents = readFileSync(v, { encoding: "utf8" });
    const hash = createHash("md5").update(contents).digest("hex");

    // Remove "static" from paths
    out[k.replace("static/", "/")] = `${v.replace("static/", "/")}?id=${hash}`;

    return out;
  }, {});
}
