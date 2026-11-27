// Switches the tabs that src/lib/tabs.mjs made. It follows the script of
// Starlight's Tabs component: a click or the arrow, Home and End keys select
// a tab and show its panel.

export {};

class UsqlTabs extends HTMLElement {
  tabs: HTMLAnchorElement[] = [];
  panels: HTMLElement[] = [];

  connectedCallback() {
    const tablist = this.querySelector<HTMLUListElement>('[role="tablist"]');
    if (tablist === null) return;
    this.tabs = [...tablist.querySelectorAll<HTMLAnchorElement>('[role="tab"]')];
    this.panels = [...this.querySelectorAll<HTMLElement>(':scope > [role="tabpanel"]')];
    this.tabs.forEach((tab, i) => {
      tab.addEventListener("click", (e) => {
        e.preventDefault();
        this.select(i);
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
        this.select(next);
        this.tabs[next].focus();
      });
    });
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
