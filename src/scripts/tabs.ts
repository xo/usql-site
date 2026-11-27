// Switches the tabs that src/lib/tabs.mjs made. It follows the script of
// Starlight's Tabs component: a click or the arrow, Home and End keys select
// a tab and show its panel.
//
// Tabs with a sync key, in data-sync, keep the choice of the reader. A choice
// opens the tab with the same label in every tabs element with the same key
// on the page, and the browser keeps the label in localStorage for the next
// pages. When the browser keeps nothing, the first tab is open.

export {};

/** storageKey returns the localStorage key of a sync key. */
function storageKey(sync: string): string {
  return `usql-tabs:${sync}`;
}

/** load returns the label that the browser keeps for a sync key, or null. */
function load(sync: string): string | null {
  try {
    return localStorage.getItem(storageKey(sync));
  } catch {
    return null;
  }
}

/** save keeps the label for a sync key, when the browser allows it. */
function save(sync: string, label: string) {
  try {
    localStorage.setItem(storageKey(sync), label);
  } catch {
    // The choice holds on this page only.
  }
}

class UsqlTabs extends HTMLElement {
  tabs: HTMLAnchorElement[] = [];
  panels: HTMLElement[] = [];

  connectedCallback() {
    const tablist = this.querySelector<HTMLUListElement>('[role="tablist"]');
    if (tablist === null) return;
    this.tabs = [...tablist.querySelectorAll<HTMLAnchorElement>('[role="tab"]')];
    this.panels = [...this.querySelectorAll<HTMLElement>(':scope > [role="tabpanel"]')];
    const sync = this.dataset.sync;
    if (sync !== undefined) {
      const label = load(sync);
      if (label !== null) this.selectLabel(label);
    }
    this.tabs.forEach((tab, i) => {
      tab.addEventListener("click", (e) => {
        e.preventDefault();
        this.choose(i);
      });
      tab.addEventListener("keydown", (e) => {
        const next =
          e.key === "ArrowLeft"
            ? i - 1
            : e.key === "ArrowRight"
              ? i + 1
              : e.key === "Home"
                ? 0
                : e.key === "End"
                  ? this.tabs.length - 1
                  : null;
        if (next === null || this.tabs[next] === undefined) return;
        e.preventDefault();
        this.choose(next);
        this.tabs[next].focus();
      });
    });
  }

  /** label returns the label of tab i. */
  label(i: number): string {
    return this.tabs[i]?.textContent?.trim() ?? "";
  }

  /** selectLabel selects the tab with the label, if this element has one. */
  selectLabel(label: string) {
    const i = this.tabs.findIndex((_, j) => this.label(j) === label);
    if (i >= 0) this.select(i);
  }

  /**
   * choose selects tab i for the reader. With a sync key, it also selects the
   * tab with the same label in the other tabs of the page with the same key,
   * and keeps the label. The tab of the reader stays at the same place on the
   * screen when the tabs above it change height.
   */
  choose(i: number) {
    this.select(i);
    const sync = this.dataset.sync;
    if (sync === undefined) return;
    const label = this.label(i);
    save(sync, label);
    const top = this.getBoundingClientRect().top;
    for (const el of document.querySelectorAll<UsqlTabs>("usql-tabs[data-sync]")) {
      if (el !== this && el.dataset.sync === sync) el.selectLabel(label);
    }
    window.scrollBy(0, this.getBoundingClientRect().top - top);
  }

  /** select shows the panel of tab i and hides the others. */
  select(i: number) {
    this.tabs.forEach((tab, j) => {
      tab.setAttribute("aria-selected", String(i === j));
      if (i === j) {
        tab.removeAttribute("tabindex");
      } else {
        tab.setAttribute("tabindex", "-1");
      }
    });
    this.panels.forEach((panel, j) => {
      panel.hidden = i !== j;
    });
  }
}

if (customElements.get("usql-tabs") === undefined) {
  customElements.define("usql-tabs", UsqlTabs);
}
