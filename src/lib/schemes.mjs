// A Markdown plugin for the URL schemes of the databases. It reads
// src/data/generated/schemes.json, which cmd/gen writes from dburl's
// registry, so that no page copies a scheme by hand (D7). A page names a
// table with a leaf directive on a line of its own:
//
//   ::schemes{name="postgres"}
//   ::databases
//
// ::schemes is the table of the schemes and aliases of one database. Its
// name is the dialect of the database in dburl, which is also the file name
// of the page of the database in src/content/docs/databases/. ::databases is
// the table of every database, with a link to each page that exists. A name
// that dburl does not know stops the build.

import { existsSync, readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";

const schemesFile = new URL("../data/generated/schemes.json", import.meta.url);
const pagesDir = new URL("../content/docs/databases/", import.meta.url);

/**
 * @typedef {{
 *   name: string, desc: string, dialect: string, aliases: string[],
 *   opaque: boolean, transport: string[], file_type: boolean,
 *   driver: string, driver_url: string, home: string, cgo: boolean,
 *   deployment: string[],
 * }} Scheme
 * @typedef {{ name: string, title: string, schemes: Scheme[] }} Database
 */

/** escape returns s with the characters that HTML reads replaced. */
function escape(s) {
  return s
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}

/** code returns each name as inline code, joined with commas. */
function code(names) {
  return names.map((n) => `<code>${escape(n)}</code>`).join(", ");
}

/** table returns an HTML table with the header head and the rows. */
function table(head, rows) {
  return (
    `<table><thead><tr>${head.map((h) => `<th>${h}</th>`).join("")}</tr></thead><tbody>` +
    rows.map((r) => `<tr>${r.map((c) => `<td>${c}</td>`).join("")}</tr>`).join("") +
    `</tbody></table>`
  );
}

/**
 * databases returns the databases of schemes.json, by dialect, sorted by
 * title. A scheme with no dialect, such as file, is not a database. The
 * scheme whose name is the dialect comes first in each database.
 *
 * @returns {Map<string, Database>}
 */
function databases() {
  /** @type {{ dburl: string, schemes: Scheme[] }} */
  const registry = JSON.parse(readFileSync(schemesFile, "utf8"));
  /** @type {Map<string, Database>} */
  const byName = new Map();
  for (const s of registry.schemes) {
    if (s.dialect === "") continue;
    const db = byName.get(s.dialect) ?? { name: s.dialect, title: "", schemes: [] };
    db.schemes.push(s);
    byName.set(s.dialect, db);
  }
  for (const db of byName.values()) {
    db.schemes.sort((a, b) => (a.name === db.name ? -1 : b.name === db.name ? 1 : 0));
    db.title = db.schemes[0].desc;
  }
  const sorted = [...byName.values()].sort((a, b) =>
    a.title.localeCompare(b.title, "en", { sensitivity: "base" })
  );
  return new Map(sorted.map((db) => [db.name, db]));
}

/** runsAs returns what the deployment kinds of a database say to a reader. */
function runsAs(db) {
  const words = { server: "a server", embedded: "a library in usql", hosted: "a hosted service" };
  const kinds = new Set(db.schemes.flatMap((s) => s.deployment));
  return Object.entries(words)
    .filter(([k]) => kinds.has(k))
    .map(([, w]) => w)
    .join(", ");
}

/** notes returns the notes of the row of a scheme. */
function notes(s) {
  const out = [];
  if (s.opaque) {
    out.push(`The rest of the URL is a path, as in <code>${escape(s.name)}:file</code>.`);
  }
  if (s.file_type) {
    out.push("usql also opens its files from a path alone, and with <code>file:</code>.");
  }
  const transports = s.transport.filter((t) => t !== "any").map((t) => `+${t}`);
  if (transports.length > 0) {
    out.push(`Transports: ${code(transports)}.`);
  }
  if (s.transport.includes("any")) {
    out.push("The transport is the name of a driver, after <code>+</code>.");
  }
  if (s.cgo) {
    out.push("The driver needs cgo.");
  }
  return out.join(" ");
}

/** schemesTable returns the table of the schemes of one database. */
function schemesTable(db) {
  return table(
    ["Scheme", "Aliases", "Go driver", "Notes"],
    db.schemes.map((s) => [
      code([s.name]),
      code(s.aliases),
      s.driver_url === ""
        ? code([s.driver])
        : `<a href="${escape(s.driver_url)}">${code([s.driver])}</a>`,
      notes(s),
    ])
  );
}

/** databasesTable returns the table of every database. */
function databasesTable(all) {
  return table(
    ["Database", "Schemes and aliases", "Runs as"],
    [...all.values()].map((db) => {
      const page = existsSync(new URL(`${db.name}.md`, pagesDir));
      return [
        page ? `<a href="/databases/${db.name}/">${escape(db.title)}</a>` : escape(db.title),
        code(db.schemes.flatMap((s) => [s.name, ...s.aliases])),
        runsAs(db),
      ];
    })
  );
}

/**
 * The Satteri plugin that turns ::schemes and ::databases directives into
 * tables.
 *
 * @returns {NonNullable<import("@astrojs/markdown-satteri").SatteriProcessorOptions["mdastPlugins"]>[number]}
 */
export function schemes() {
  /** @type {Map<string, Database> | undefined} */
  let all;
  return {
    name: "usql-schemes",
    leafDirective(node, ctx) {
      if (node.name !== "schemes" && node.name !== "databases") return;
      const where = ctx.fileURL ? fileURLToPath(ctx.fileURL) : "a page";
      all ??= databases();
      if (node.name === "databases") {
        return { type: "html", value: databasesTable(all) };
      }
      const name = node.attributes?.name ?? "";
      const db = all.get(name);
      if (db === undefined) {
        throw new Error(
          `${where}: ::schemes{name="${name}"} names no database of dburl. ` +
            `Use a dialect of src/data/generated/schemes.json, such as postgres`
        );
      }
      return { type: "html", value: schemesTable(db) };
    },
  };
}
