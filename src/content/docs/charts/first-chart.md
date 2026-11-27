---
title: Your first chart
description: Draw a bar chart and a line chart of a query in the terminal with \chart.
sidebar:
  order: 1
---

`\chart` runs the query in the buffer and draws its result as a chart in the terminal. It works in a terminal that draws images, such as kitty, ghostty, WezTerm, iTerm2 or foot.

## The shape of the result

The first column that is not a number gives the labels of the x axis. Each column that is a number is one series of the chart.

```sql
select to_char(month, 'Mon') as month, coffee, tea, cocoa
from sales
order by sales.month
\chart title='Cups sold in 2026' subtitle='by month' type=bar
```

## The options

| Option | What it does |
| --- | --- |
| `title` | the title of the chart |
| `subtitle` | the subtitle |
| `type` | `bar` or `line` |
| `size` | the size in pixels, as `WIDTHxHEIGHT`. The default is `800x600` |
| `bg` | the color of the background. In the current release of usql, it changes nothing |
| `prec` | the number of decimal places of the values |
| `file` | write the chart to this SVG file, and draw nothing |

A line chart in another size:

```sql
select to_char(month, 'Mon') as month, coffee, tea, cocoa
from sales
order by sales.month
\chart title='Cups sold in 2026' type=line size=900x450
```

## When usql draws nothing

usql asks the terminal at startup which images it draws. If the terminal does not answer, set `USQL_TERM_GRAPHICS` to `kitty`, `iterm` or `sixel`. usql turns the graphics off when it reads a file from `-f`, so run a chart script with `\i` at the prompt.
