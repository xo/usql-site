---
title: Save a chart to a file
description: Write a chart to an SVG file with \chart file=, to open it in a browser or to put it in a document.
sidebar:
  order: 4
---

`file=NAME` writes the chart to an SVG file, and usql draws nothing in the terminal. A browser opens the file, and a page or a document can show it.

::cast[Write a bar chart of the coffee sold each month to coffee.svg, and look at the file.]{name="chart-file"}

```sql
select month, coffee from sales order by month_no
\chart title='Coffee sold' file=coffee.svg
```

usql prints `writing to coffee.svg`. This is the file:

![A bar chart titled Coffee sold, with one bar for each month from January to June.](/charts/chart-file/coffee.svg)

In a terminal that draws images, the same `\chart` without `file=` draws the chart in place. The example uses the tables of the [gallery](/charts/gallery/#the-data).

## The file

- usql writes SVG, whatever the name of the file ends with. Give the file a name that ends with `.svg`.
- A name with no directory is in the working directory of usql.
- The directory of the file must exist. If it does not, `\chart` stops with an error such as `open charts/coffee.svg: no such file or directory`.
- usql replaces a file that exists, and does not ask first.
- `size` sets the size of the SVG file, as it does for a chart in the terminal. The default is `800x600`.

## In a terminal that draws no images

In the current release of usql, `\chart` needs graphics also for `file=`. In a terminal that draws no images, it stops with an error that ends with this text, even with `file=`:

```text
\chart: graphics not supported in terminal
```

To write the file there, choose a graphics protocol first. With `file=`, usql then writes the file and draws nothing:

```sql
\uset term_graphics kitty
\chart title='Coffee sold' file=coffee.svg
```

[Terminal graphics](/guides/terminal-graphics/) says how usql chooses the graphics.
