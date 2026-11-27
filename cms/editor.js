// The editor components of the site, for the Markdown body of a page. Each
// one turns a block of the Markdown into a form and back. Each block is the
// one that a plugin of src/lib/ reads.

// ::cast[caption]{name="tour"}, which src/lib/cast.mjs reads.
const castRE = /^::cast\[([^\]\n]*)\]\{name="([a-z0-9-]+)"(?: poster="(\d{1,2}:\d{2})")?\}$/m;

/** escape returns s with the characters that HTML reads replaced. */
function escape(s) {
  return String(s ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}

CMS.registerEditorComponent({
  id: "cast",
  label: "Recording",
  fields: [
    {
      name: "name",
      label: "Recording (the file name in public/casts, without .cast)",
      widget: "string",
      pattern: ["^[a-z0-9][a-z0-9-]*$", "Use a-z, 0-9 and -"],
    },
    { name: "caption", label: "Caption", widget: "string", required: false },
    {
      name: "poster",
      label: "Frame to show before it plays (M:SS)",
      widget: "string",
      required: false,
      pattern: ["^(\\d{1,2}:\\d{2})?$", "Use M:SS, such as 0:12"],
    },
  ],
  pattern: castRE,
  fromBlock: (match) => ({ caption: match[1], name: match[2], poster: match[3] ?? "" }),
  toBlock: ({ name, caption, poster }) =>
    `::cast[${(caption ?? "").replaceAll("]", "")}]{name="${name}"` +
    (poster ? ` poster="${poster}"` : "") +
    "}",
  toPreview: ({ name, caption }) =>
    `<p><strong>Recording ${escape(name)}</strong>: ${escape(caption)}</p>`,
});

// ::schemes{name="postgres"} and ::databases, which src/lib/schemes.mjs reads.
const schemesRE = /^::schemes\{name="([a-z0-9-]+)"\}$/m;

CMS.registerEditorComponent({
  id: "schemes",
  label: "Schemes of a database",
  fields: [
    {
      name: "name",
      label: "Database (its dialect in dburl, such as postgres or sqlite3)",
      widget: "string",
      pattern: ["^[a-z0-9][a-z0-9-]*$", "Use a-z, 0-9 and -"],
    },
  ],
  pattern: schemesRE,
  fromBlock: (match) => ({ name: match[1] }),
  toBlock: ({ name }) => `::schemes{name="${name ?? ""}"}`,
  toPreview: ({ name }) => `<p><strong>The schemes and aliases of ${escape(name)}</strong></p>`,
});

CMS.registerEditorComponent({
  id: "databases",
  label: "Table of every database",
  fields: [],
  pattern: /^::databases$/m,
  fromBlock: () => ({}),
  toBlock: () => "::databases",
  toPreview: () => "<p><strong>The table of every database</strong></p>",
});

// :::tabs, or :::tabs{sync="key"}, which src/lib/tabs.mjs reads. A tab is
// ::::tab[label], its Markdown, and :::: on a line of its own. The block ends
// at the first line ::: that follows a line ::::, so an aside in a tab, which
// ends with :::, does not end it.
const tabsRE = /^:::tabs(?:\{sync="([a-z0-9-]+)"\})?\n([\s\S]*?\n::::)\n:::$/m;
const tabRE = /^::::tab\[([^\]\n]*)\]\n([\s\S]*?)\n?::::$/gm;

CMS.registerEditorComponent({
  id: "tabs",
  label: "Tabs",
  fields: [
    {
      name: "sync",
      label:
        "Sync key (optional). Tabs with the same key keep the choice of the reader across pages",
      widget: "string",
      required: false,
      pattern: ["^([a-z0-9][a-z0-9-]*)?$", "Use a-z, 0-9 and -"],
    },
    {
      name: "tabs",
      label: "Tabs",
      widget: "list",
      fields: [
        { name: "label", label: "Label, such as macOS", widget: "string" },
        { name: "body", label: "Content", widget: "markdown" },
      ],
    },
  ],
  pattern: tabsRE,
  fromBlock: (match) => ({
    sync: match[1] ?? "",
    tabs: [...match[2].matchAll(tabRE)].map((m) => ({ label: m[1], body: m[2] })),
  }),
  toBlock: ({ sync, tabs }) =>
    [
      sync ? `:::tabs{sync="${sync}"}` : ":::tabs",
      ...(tabs ?? []).map(
        ({ label, body }) =>
          `::::tab[${(label ?? "").replaceAll("]", "")}]\n` +
          ((body ?? "").trim() === "" ? "" : `${body.trim()}\n`) +
          "::::"
      ),
      ":::",
    ].join("\n"),
  toPreview: ({ sync, tabs }) =>
    `<p><strong>Tabs</strong>${sync ? ` (sync key ${escape(sync)})` : ""}: ` +
    `${escape((tabs ?? []).map((t) => t.label).join(", "))}</p>`,
});

// ::releases and ::releases{limit="20"}, which src/lib/releases.mjs reads.
const releasesRE = /^::releases(?:\{limit="([1-9][0-9]*)"\})?$/m;

CMS.registerEditorComponent({
  id: "releases",
  label: "Release notes",
  fields: [
    {
      name: "limit",
      label: "Number of the newest releases to show (empty for every release)",
      widget: "string",
      required: false,
      pattern: ["^([1-9][0-9]*)?$", "Use a number above 0, or leave it empty"],
    },
  ],
  pattern: releasesRE,
  fromBlock: (match) => ({ limit: match[1] ?? "" }),
  toBlock: ({ limit }) => (limit ? `::releases{limit="${limit}"}` : "::releases"),
  toPreview: ({ limit }) =>
    `<p><strong>Release notes</strong>: ${limit ? `the ${escape(limit)} newest releases` : "every release"}</p>`,
});
