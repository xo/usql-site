// A Markdown plugin for tabs, one for each choice the reader can make, such
// as the system to install usql on. A page writes them as a container
// directive that holds one container directive for each tab:
//
//   :::tabs
//   ::::tab[macOS]
//   Markdown for macOS.
//   ::::
//   ::::tab[Linux]
//   Markdown for Linux.
//   ::::
//   :::
//
// A tab has four colons, so that an aside such as :::note can sit in it. The
// label in brackets is the text of the tab. The plugin turns the directive
// into the markup of Starlight's Tabs component, and src/scripts/tabs.ts
// switches the tabs. A tabs directive with anything other than tabs in it,
// a tab with no label, and a tab outside a tabs directive stop the build.
//
// An optional sync key, as in :::tabs{sync="os"}, keeps the choice of the
// reader, as Starlight's syncKey does. When the reader picks a tab, every
// tabs directive with the same key, on this page and on the next pages,
// opens the tab with the same label. The browser keeps the label in its
// localStorage. A short script after the tabs opens the kept tab before the
// page shows, and the page works without it.

import { fileURLToPath } from "node:url";

const syncRE = /^[a-z0-9][a-z0-9-]*$/;

/**
 * restore opens the tab whose label the browser keeps for the sync key of a
 * usql-tabs element. It runs inline, in a script that is the last child of
 * the element, so that the page never shows the first tab first. It does the
 * same as the select method of src/scripts/tabs.ts, which runs later.
 */
function restore() {
  try {
    const el = document.currentScript?.parentElement;
    const label = localStorage.getItem(`usql-tabs:${el?.dataset.sync}`);
    if (!el || label === null) return;
    const tabs = [...el.querySelectorAll(':scope > div > [role="tablist"] [role="tab"]')];
    const i = tabs.findIndex((tab) => tab.textContent?.trim() === label);
    if (i < 0) return;
    const panels = [...el.querySelectorAll(':scope > [role="tabpanel"]')];
    tabs.forEach((tab, j) => {
      tab.setAttribute("aria-selected", String(i === j));
      tab.setAttribute("tabindex", i === j ? "0" : "-1");
    });
    panels.forEach((panel, j) => {
      if (panel instanceof HTMLElement) panel.hidden = i !== j;
    });
  } catch {
    // The browser keeps no choice, so the first tab stays open.
  }
}

/** escape returns s with the characters that HTML reads replaced. */
function escape(s) {
  return s
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}

/** element returns an mdast node that becomes the HTML element tag. */
function element(tag, properties, children) {
  return { type: "paragraph", data: { hName: tag, hProperties: properties }, children };
}

/**
 * The Satteri plugin that turns :::tabs directives into tabs. It is a
 * factory, so that each page counts its own tabs for their ids.
 *
 * @returns {NonNullable<import("@astrojs/markdown-satteri").SatteriProcessorOptions["mdastPlugins"]>[number]}
 */
export function tabs() {
  return () => {
    let count = 0;
    return {
      name: "usql-tabs",
      containerDirective(node, ctx) {
        const where = ctx.fileURL ? fileURLToPath(ctx.fileURL) : "a page";
        if (node.name === "tab") {
          const parent = ctx.parent(node);
          if (parent?.type !== "containerDirective" || parent.name !== "tabs") {
            throw new Error(`${where}: ::::tab must be inside :::tabs`);
          }
          return;
        }
        if (node.name !== "tabs") return;
        const sync = node.attributes?.sync ?? "";
        if (sync !== "" && !syncRE.test(sync)) {
          throw new Error(`${where}: :::tabs sync must use a-z, 0-9 and -, not ${sync}`);
        }
        const set = count++;
        const labels = [];
        const panels = [];
        for (const [i, tab] of node.children.entries()) {
          if (tab.type !== "containerDirective" || tab.name !== "tab") {
            throw new Error(
              `${where}: :::tabs can hold only ::::tab[label] sections, ended by ::::`
            );
          }
          const children = [...tab.children];
          const first = children[0];
          const label =
            first?.type === "paragraph" && first.data?.directiveLabel
              ? ctx.textContent(first).trim()
              : "";
          if (label === "") {
            throw new Error(`${where}: a tab needs a label, as in ::::tab[macOS]`);
          }
          children.shift();
          const tabID = `usql-tab-${set}-${i}`;
          const panelID = `usql-tab-panel-${set}-${i}`;
          labels.push(
            `<li role="presentation" class="tab">` +
              `<a role="tab" href="#${panelID}" id="${tabID}"` +
              ` aria-selected="${i === 0}" tabindex="${i === 0 ? 0 : -1}">${escape(label)}</a>` +
              `</li>`
          );
          const properties = { id: panelID, role: "tabpanel", "aria-labelledby": tabID };
          if (i !== 0) properties.hidden = true;
          panels.push(element("div", properties, children));
        }
        if (panels.length === 0) {
          throw new Error(`${where}: :::tabs holds no ::::tab[label] section`);
        }
        const children = [
          {
            type: "html",
            value: `<div class="tablist-wrapper not-content"><ul role="tablist">${labels.join("")}</ul></div>`,
          },
          ...panels,
        ];
        if (sync === "") return element("usql-tabs", {}, children);
        children.push({ type: "html", value: `<script>(${restore})()</script>` });
        return element("usql-tabs", { "data-sync": sync }, children);
      },
    };
  };
}
