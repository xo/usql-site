---
title: Syntax highlighting
description: Choose the colors of the SQL that you type, or turn them off.
sidebar:
  order: 7
---

usql colors the SQL and the meta commands as you type them. It reads each statement in the grammar of the database that it is connected to. A style sets the colors. The default style is `monokai`.

::cast[The same statement in the styles monokai, dracula and catppuccin-mocha, and then with no highlighting.]{name="syntax-highlighting"}

## Choose a style

`\? styles` lists the names of the styles, and names the style in use. `\uset syntax_hl_style` changes the style:

```sql
\uset syntax_hl_style dracula
```

To keep a style, put it in the `usql` section of `config.yaml`:

```yaml
usql:
  syntax_hl_style: dracula
```

[Configure usql](/guides/configure/) says where `config.yaml` is.

## The settings

| Setting | What it does |
| --- | --- |
| `syntax_hl` | turns highlighting on or off |
| `syntax_hl_style` | the name of the style |
| `syntax_hl_format` | the number of colors: `terminal` is 16, `terminal256` is 256, `terminal16m` is 24-bit color, and `noop` is none |
| `syntax_hl_override_bg` | when it is `true`, usql keeps the background color of the terminal and not the background of the style |

usql sets `syntax_hl_format` from what your terminal can show. If the colors look wrong, set it to `terminal256` or `terminal`:

```sql
\uset syntax_hl_format terminal256
```

## Turn highlighting off

```sql
\uset syntax_hl off
```

usql turns highlighting off when the terminal shows no color, or when the `NO_COLOR` environment variable is set.
