// Fills each recording figure that src/lib/cast.mjs made with the asciinema
// player. The player is loaded only on a page that has a recording.

export {};

const figures = document.querySelectorAll<HTMLElement>("figure.usql-cast[data-cast]");

if (figures.length > 0) {
  const [{ create }] = await Promise.all([
    import("asciinema-player"),
    import("asciinema-player/dist/bundle/asciinema-player.css"),
  ]);
  for (const figure of figures) {
    const target = figure.querySelector<HTMLElement>(".usql-cast__player");
    if (target === null) continue;
    target.replaceChildren();
    create(figure.dataset.cast as string, target, {
      fit: "width",
      theme: "usql",
      terminalFontFamily:
        'ui-monospace, SFMono-Regular, Menlo, Consolas, "Liberation Mono", "DejaVu Sans Mono", monospace',
      idleTimeLimit: 2,
      poster: figure.dataset.poster,
      preload: true,
    });
  }
}
