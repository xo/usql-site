// Runs before `astro dev` and `astro build`. It makes public/admin/, the web
// editor of the site, from cms/ and from the Sveltia CMS package of
// package.json:
//
// 1. Removes public/admin/, which an earlier run made.
// 2. Writes public/admin/config.yml: cms/config.yml, with one collection for
//    each file of cms/collections/, in file name order. A collection with no
//    fields gets the fields of cms/fields/page.yml. CMS_AUTH_URL fills in the
//    sign-in service. Without it, base_url is left out, and an editor signs
//    in with a GitHub token.
// 3. Copies cms/admin.html, cms/editor.js, the script of Sveltia CMS and its
//    chunks. The script loads its fonts from cdn.jsdelivr.net. They are
//    copied from their npm packages, and the script points at the copies.
//
//   node scripts/prepare.mjs

import { cp, mkdir, readdir, readFile, rm, writeFile } from "node:fs/promises";
import { join } from "node:path";
import { parse, stringify } from "yaml";

const root = new URL("..", import.meta.url).pathname;
const adminDir = join(root, "public", "admin");
const sveltiaDir = join(root, "node_modules", "@sveltia", "cms", "dist");

await rm(adminDir, { recursive: true, force: true });

// config.yml
const config = parse(await readFile(join(root, "cms", "config.yml"), "utf8"));
const authURL = process.env.CMS_AUTH_URL ?? "";
if (authURL === "") {
  delete config.backend.base_url;
} else {
  config.backend.base_url = authURL.replace(/\/$/, "");
}
const pageFields = parse(await readFile(join(root, "cms", "fields", "page.yml"), "utf8"));
config.collections = [];
const files = (await readdir(join(root, "cms", "collections"))).filter((f) => f.endsWith(".yml"));
for (const file of files.sort()) {
  const collection = parse(await readFile(join(root, "cms", "collections", file), "utf8"));
  if (!collection?.name) throw new Error(`cms/collections/${file} has no name`);
  if (!collection.files && !collection.fields) collection.fields = pageFields;
  config.collections.push(collection);
}
await mkdir(adminDir, { recursive: true });
const header = "# Made by scripts/prepare.mjs from cms/. Do not edit.\n";
await writeFile(join(adminDir, "config.yml"), header + stringify(config, { lineWidth: 0 }));
await cp(join(root, "cms", "admin.html"), join(adminDir, "index.html"));
await cp(join(root, "cms", "editor.js"), join(adminDir, "editor.js"));

// The fonts of Sveltia CMS, by the address that its script loads them from.
// If a new version of Sveltia CMS asks for another font, the build stops
// here: add its npm package to package.json and a line here.
const fonts = [
  [
    "https://cdn.jsdelivr.net/fontsource/fonts/source-sans-3:vf@5.3.1/latin-wght-normal.woff2",
    "@fontsource-variable/source-sans-3/files/source-sans-3-latin-wght-normal.woff2",
  ],
  [
    "https://cdn.jsdelivr.net/fontsource/fonts/source-sans-3:vf@5.3.1/latin-wght-italic.woff2",
    "@fontsource-variable/source-sans-3/files/source-sans-3-latin-wght-italic.woff2",
  ],
  [
    "https://cdn.jsdelivr.net/fontsource/fonts/noto-mono@5.3.0/latin-400-normal.woff2",
    "@fontsource/noto-mono/files/noto-mono-latin-400-normal.woff2",
  ],
  [
    "https://cdn.jsdelivr.net/fontsource/fonts/material-symbols-outlined:vf@5.3.9/latin-wght-normal.woff2",
    "@fontsource-variable/material-symbols-outlined/files/material-symbols-outlined-latin-wght-normal.woff2",
  ],
];
let script = await readFile(join(sveltiaDir, "sveltia-cms.js"), "utf8");
await mkdir(join(adminDir, "fonts"), { recursive: true });
for (const [url, file] of fonts) {
  const name = file.split("/").pop();
  if (!script.includes(url)) throw new Error(`cms: sveltia-cms.js no longer loads ${url}`);
  script = script.replaceAll(url, `/admin/fonts/${name}`);
  await cp(join(root, "node_modules", file), join(adminDir, "fonts", name));
}
const left = script.match(/https:\/\/cdn\.jsdelivr\.net\/fontsource\/[^`"')]+/);
if (left) throw new Error(`cms: sveltia-cms.js loads a font that is not served here: ${left[0]}`);
await writeFile(join(adminDir, "sveltia-cms.js"), script);

// The script loads its chunks from the directory that it is in, and from
// unpkg.com only when they are not there.
await cp(join(sveltiaDir, "chunks"), join(adminDir, "chunks"), {
  recursive: true,
  filter: (src) => !src.endsWith(".map"),
});

console.log(
  `cms: wrote public/admin/ with ${config.collections.length} collections` +
    (authURL === "" ? ", and sign-in with a token" : `, and sign-in through ${authURL}`)
);
