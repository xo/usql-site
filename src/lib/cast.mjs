// A Markdown plugin for the recordings of the site. A page names a recording
// with a leaf directive on a line of its own:
//
//   ::cast[A first session with usql]{name="tour"}
//
// The label in brackets is the caption, and name is the file
// public/casts/<name>.cast, which cmd/record writes. An optional poster="0:12"
// shows the frame at that time before the recording plays. The plugin turns
// the directive into a figure that src/scripts/cast.ts fills with the player.
// A missing recording stops the build. See docs/RECORDING.md.

import { existsSync } from "node:fs";
import { fileURLToPath } from "node:url";

const castDir = new URL("../../public/casts/", import.meta.url);
const nameRE = /^[a-z0-9][a-z0-9-]*$/;
const posterRE = /^\d{1,2}:\d{2}$/;

/** escape returns s with the characters that HTML reads replaced. */
function escape(s) {
  return s
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}

/**
 * The Satteri plugin that turns ::cast directives into figures.
 *
 * @returns {NonNullable<import("@astrojs/markdown-satteri").SatteriProcessorOptions["mdastPlugins"]>[number]}
 */
export function cast() {
  return {
    name: "usql-cast",
    leafDirective(node, ctx) {
      if (node.name !== "cast") return;
      const where = ctx.fileURL ? fileURLToPath(ctx.fileURL) : "a page";
      const name = node.attributes?.name ?? "";
      if (!nameRE.test(name)) {
        throw new Error(`${where}: ::cast needs name="<recording>", with a-z, 0-9 and -`);
      }
      if (!existsSync(new URL(`${name}.cast`, castDir))) {
        throw new Error(
          `${where}: public/casts/${name}.cast does not exist. Run go run ./cmd/record ${name}`
        );
      }
      const poster = node.attributes?.poster ?? "";
      if (poster !== "" && !posterRE.test(poster)) {
        throw new Error(`${where}: ::cast poster must be M:SS, not ${poster}`);
      }
      const caption = ctx.textContent(node).trim();
      const src = `/casts/${name}.cast`;
      const attrs = [`class="usql-cast"`, `data-cast="${src}"`];
      if (poster !== "") attrs.push(`data-poster="npt:${poster}"`);
      const label = caption === "" ? name : caption;
      return {
        type: "html",
        value:
          `<figure ${attrs.join(" ")}>` +
          `<div class="usql-cast__player"><a href="${src}">${escape(label)}</a></div>` +
          (caption === "" ? "" : `<figcaption>${escape(caption)}</figcaption>`) +
          `</figure>`,
      };
    },
  };
}
