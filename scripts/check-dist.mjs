// Checks the built site in dist/ for a page that loads something from
// another host (D6): a script, a style, a font, an image, a frame or a media
// file. An SVG file, such as a chart, is checked for an image, a script or a
// style that it loads. A link that a reader follows is not a load, and is not
// checked. The editor at dist/admin/ is not checked.
//
//   node scripts/check-dist.mjs

import { readdir, readFile } from "node:fs/promises";
import { join, relative } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const dist = join(root, "dist");
const own = new URL(process.env.SITE_URL ?? "https://usql.app").host;

async function* walk(dir) {
  for (const entry of await readdir(dir, { withFileTypes: true })) {
    const path = join(dir, entry.name);
    if (entry.isDirectory()) {
      if (path === join(dist, "admin")) continue;
      yield* walk(path);
    } else {
      yield path;
    }
  }
}

function foreign(url) {
  if (!/^(https?:)?\/\//.test(url)) return false;
  return new URL(url, "https://" + own).host !== own;
}

// The tags and attributes that load what they name. A <link> loads only for
// some values of rel.
const loads =
  /<(script|img|iframe|source|video|audio|embed|object)\b[^>]*?\s(?:src|data|poster)="([^"]+)"/gi;
const links = /<link\b([^>]*)>/gi;
const loadingRel =
  /\brel="[^"]*\b(stylesheet|preload|modulepreload|prefetch|icon|manifest|apple-touch-icon)\b/i;
const cssURL = /url\(\s*["']?([^"')]+)["']?\s*\)|@import\s+["']([^"']+)["']/gi;
// The attributes of SVG that load what they name. xmlns names a namespace and
// loads nothing.
const svgLoads = /\s(?:href|xlink:href|src)="([^"]+)"/gi;

const problems = [];
for await (const path of walk(dist)) {
  const name = relative(root, path);
  if (path.endsWith(".html")) {
    const html = await readFile(path, "utf8");
    for (const [, tag, url] of html.matchAll(loads)) {
      if (foreign(url)) problems.push(`${name}: <${tag}> loads ${url}`);
    }
    for (const [, attrs] of html.matchAll(links)) {
      const href = attrs.match(/\bhref="([^"]+)"/)?.[1];
      if (href && loadingRel.test(attrs) && foreign(href)) {
        problems.push(`${name}: <link> loads ${href}`);
      }
    }
  } else if (path.endsWith(".svg")) {
    const svg = await readFile(path, "utf8");
    for (const [, url] of svg.matchAll(svgLoads)) {
      if (foreign(url)) problems.push(`${name}: loads ${url}`);
    }
    for (const [, a, b] of svg.matchAll(cssURL)) {
      const url = a ?? b;
      if (foreign(url)) problems.push(`${name}: loads ${url}`);
    }
  } else if (path.endsWith(".css")) {
    const css = await readFile(path, "utf8");
    for (const [, a, b] of css.matchAll(cssURL)) {
      const url = a ?? b;
      if (foreign(url)) problems.push(`${name}: loads ${url}`);
    }
  }
}

if (problems.length > 0) {
  console.error("a page loads from another host:\n" + problems.join("\n"));
  process.exit(1);
}
console.log(`no page loads from another host than ${own}`);
