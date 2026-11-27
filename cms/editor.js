// The editor components of the site, for the Markdown body of a page. Each
// one turns a block of the Markdown into a form and back. The Markdown is the
// one that src/lib/cast.mjs reads.

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
