// A Markdown plugin for the release notes of usql. It reads
// src/data/generated/releases.json, which `go run ./cmd/gen -releases` writes
// from the releases of usql on GitHub (D7). A page shows them with a leaf
// directive on a line of its own:
//
//   ::releases
//   ::releases{limit="20"}
//
// limit shows only that number of the newest releases. The plugin groups the
// releases by year, newest first. Each release is a heading with its tag, the
// date and a link to the release on GitHub, and its notes.
//
// The notes are Markdown that the maintainers write on GitHub. The plugin
// parses them with the same parser as the pages, and makes them safe for a
// page of this site:
//
//   - Raw HTML shows as text.
//   - An image becomes a link to the image, so that no page loads it (D6).
//   - A link keeps its target only when it is an absolute http, https or
//     mailto URL. Another link shows as its text.
//   - A heading of the notes moves two levels down, and is at least a level 4
//     heading, under the level 3 heading of its release. The table of
//     contents lists only the years and the releases.
//   - Directives, footnotes and frontmatter are not read.

import { existsSync, readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { markdownToMdast } from "satteri";

const releasesFile = new URL("../data/generated/releases.json", import.meta.url);
const limitRE = /^[1-9][0-9]*$/;

/**
 * @typedef {{ tag: string, name: string, date: string, body: string, url: string }} Release
 * @typedef {import("satteri").MdastNode} Node
 */

/** The node types whose children are phrasing content, such as text. */
const phrasing = new Set([
  "paragraph",
  "heading",
  "tableCell",
  "emphasis",
  "strong",
  "delete",
  "link",
  "linkReference",
]);

/** safeURL reports whether a link of the notes can keep its target. */
export function safeURL(url) {
  if (!URL.canParse(url)) return false;
  return ["http:", "https:", "mailto:"].includes(new URL(url).protocol);
}

/** text returns a text node. */
function text(value) {
  return { type: "text", value };
}

/** link returns a link node. */
function link(url, children) {
  return { type: "link", url, title: null, children };
}

/**
 * clean returns the nodes of the notes, made safe as the comment at the top of
 * this file says. defs holds the URLs of the link definitions of the notes,
 * by identifier. inPhrasing is true when the nodes are phrasing content, and
 * inLink is true when they are inside a link.
 *
 * @param {Node[]} nodes
 * @param {Map<string, string>} defs
 * @param {boolean} inPhrasing
 * @param {boolean} inLink
 * @returns {Node[]}
 */
function clean(nodes, defs, inPhrasing, inLink) {
  /** @type {Node[]} */
  const out = [];
  for (const node of nodes) {
    switch (node.type) {
      case "html":
        out.push(
          inPhrasing ? text(node.value) : { type: "paragraph", children: [text(node.value)] }
        );
        break;
      case "image":
      case "imageReference": {
        const url = node.type === "image" ? node.url : (defs.get(node.identifier) ?? "");
        const label = node.alt || url;
        out.push(inLink || !safeURL(url) ? text(label) : link(url, [text(label)]));
        break;
      }
      case "link":
      case "linkReference": {
        const url = node.type === "link" ? node.url : (defs.get(node.identifier) ?? "");
        const children = clean(node.children, defs, true, true);
        if (inLink || !safeURL(url)) {
          out.push(...children);
        } else {
          out.push(link(url, children));
        }
        break;
      }
      case "definition":
      case "footnoteDefinition":
      case "footnoteReference":
      case "yaml":
      case "toml":
        break;
      case "heading":
        out.push({
          ...node,
          depth: Math.min(Math.max(node.depth + 2, 4), 6),
          children: clean(node.children, defs, true, inLink),
        });
        break;
      default:
        if ("children" in node) {
          out.push({
            ...node,
            children: clean(node.children, defs, phrasing.has(node.type), inLink),
          });
        } else {
          out.push(node);
        }
    }
  }
  return out;
}

/** definitions returns the URLs of the link definitions of a tree, by identifier. */
function definitions(node, defs = new Map()) {
  if (node.type === "definition" && !defs.has(node.identifier)) {
    defs.set(node.identifier, node.url);
  }
  for (const child of node.children ?? []) definitions(child, defs);
  return defs;
}

/**
 * notes returns the mdast nodes of the Markdown of the notes of a release,
 * made safe for a page of the site.
 *
 * @param {string} body
 * @returns {Node[]}
 */
export function notes(body) {
  const root = markdownToMdast(body, {
    position: false,
    features: { gfm: { footnotes: false }, frontmatter: false },
  });
  return clean(root.children, definitions(root), false, false);
}

/**
 * releaseNodes returns the mdast nodes of one release: its heading, the line
 * with its date, and its notes. The line names the release when its name is
 * more than its tag, as in "usql v0.21.6".
 *
 * @param {Release} r
 * @returns {Node[]}
 */
function releaseNodes(r) {
  const named = r.name !== r.tag && r.name !== `usql ${r.tag}`;
  const line = [text(`${named ? `${r.name}. ` : ""}Released on ${r.date}.`)];
  if (safeURL(r.url)) {
    line.push(text(" "), link(r.url, [text("The release on GitHub")]), text(" has its files."));
  }
  return [
    { type: "heading", depth: 3, children: [text(r.tag)] },
    { type: "paragraph", children: line },
    ...notes(r.body),
  ];
}

/**
 * The Satteri plugin that turns ::releases directives into the release notes.
 *
 * @returns {NonNullable<import("@astrojs/markdown-satteri").SatteriProcessorOptions["mdastPlugins"]>[number]}
 */
export function releases() {
  /** @type {Release[] | undefined} */
  let all;
  return {
    name: "usql-releases",
    leafDirective(node, ctx) {
      if (node.name !== "releases") return;
      const where = ctx.fileURL ? fileURLToPath(ctx.fileURL) : "a page";
      const limit = node.attributes?.limit ?? "";
      if (limit !== "" && !limitRE.test(limit)) {
        throw new Error(`${where}: ::releases limit must be a number above 0, not ${limit}`);
      }
      if (all === undefined) {
        if (!existsSync(releasesFile)) {
          throw new Error(
            `${where}: src/data/generated/releases.json does not exist. Run go run ./cmd/gen -releases`
          );
        }
        all = JSON.parse(readFileSync(releasesFile, "utf8")).releases;
      }
      const shown = limit === "" ? all : all.slice(0, Number(limit));
      /** @type {Node[]} */
      const children = [];
      let year = "";
      for (const r of shown) {
        if (r.date.slice(0, 4) !== year) {
          year = r.date.slice(0, 4);
          children.push({ type: "heading", depth: 2, children: [text(year)] });
        }
        children.push(...releaseNodes(r));
      }
      return { type: "paragraph", data: { hName: "div" }, children };
    },
  };
}
