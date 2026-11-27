---
title: Size, colors and titles
description: Give a chart a title and a subtitle, name its series, and set its size in pixels.
sidebar:
  order: 3
---

The options of `\chart` set the title, the subtitle and the size of a chart. The names of the columns of the query name the series.

::cast[Name the series with as, give a chart a title and a subtitle, and draw it again in a smaller size.]{name="chart-options"}

The recording writes each chart to an SVG file with `file=`, because the player of the recordings draws no images. In a terminal that draws images, the same `\chart` without `file=` draws the chart in place. The examples use the tables of the [gallery](/charts/gallery/#the-data).

## A title and a subtitle

`title` sets the title above the chart, and `subtitle` sets the line under it. Put a value that holds a space in single quotes:

```sql
select month, coffee as "Coffee", tea as "Tea"
from sales order by month_no
\chart title='Cups sold' subtitle='January to June 2026'
```

![A bar chart titled Cups sold, with the subtitle January to June 2026, of the cups of coffee and tea sold each month.](/charts/chart-options/titled.svg)

## The names of the series

The legend shows the name of each column. To change a name, rename the column in the query with `as`. The query above shows `Coffee` and `Tea` in place of `coffee` and `tea`. Put a name in double quotes, because some databases, such as PostgreSQL, change a name with no quotes to lower case.

## The size

`size` sets the width and the height of the chart in pixels, as `WIDTHxHEIGHT`. The default is `800x600`. The title, the axes and the legend are inside that size.

```sql
\chart title='Cups sold' size=640x320
```

![The same bar chart of coffee and tea, 640 pixels wide and 320 pixels high.](/charts/chart-options/small.svg)

With an empty query buffer, `\chart` draws the last query again, so this line draws the query of the first chart in a new size.

## The colors

usql chooses the color of each series. No option changes the colors of the series.

The option `bg` takes a color, such as `bg=black` or `bg='#f0f0f0'`. In the current release of usql, `bg` changes nothing: a chart has a white background, in the terminal and in a file.
