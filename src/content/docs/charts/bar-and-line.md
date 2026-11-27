---
title: Bar and line charts
description: How usql makes the series and the labels of a chart from the columns of a result, and how to draw bars or lines.
sidebar:
  order: 2
---

`\chart` draws the result of a query as a bar chart or a line chart. The columns of the result give the labels and the series of the chart. The option `type` chooses bars or lines.

::cast[Draw the sales of each month as bars and as lines, and the cups of each hour as a line, and write each chart to an SVG file.]{name="chart-bar-line"}

The recording writes each chart to an SVG file with `file=`, because the player of the recordings draws no images. In a terminal that draws images, the same `\chart` without `file=` draws the chart in place. The examples use the tables of a small cafe. The [gallery](/charts/gallery/#the-data) has the SQL that makes them.

## The columns of the result

usql reads the columns of the result in this way:

1. A column is a number when each value in it is a number.
2. The first column that is not a number gives the labels of the x axis.
3. Each other column that is a number is a series. The legend shows the name of the column.

This query has a column of labels, `month`, and three series:

```sql
select month, coffee, tea, cocoa from sales order by month_no;
\chart title='Cups sold'
```

![A bar chart of the cups of coffee, tea and cocoa sold each month from January to June. Each month has three bars side by side.](/charts/chart-bar-line/sales-bar.svg)

A result with a column of labels gives a bar chart. Each label has one bar for each series, side by side.

`\chart` with an empty query buffer draws the last query again. The query above ends with a semicolon, so usql runs it and prints the rows. The `\chart` on the next line draws the same query.

## A line chart

`type=line` draws the same result as lines:

```sql
\chart title='Cups sold' type=line
```

![A line chart of the cups of coffee, tea and cocoa sold each month. Coffee goes up, and tea and cocoa go down.](/charts/chart-bar-line/sales-line.svg)

A line chart shows the name of each series at the end of its line.

When each column of the result is a number, the first column gives the x axis, and usql draws a line chart:

```sql
select hour, cups from hours order by hour
\chart title='Cups by hour'
```

![A line chart of the cups sold each hour from 7 to 17, with the most cups at 8.](/charts/chart-bar-line/hours.svg)

To draw this result as bars, add `type=bar`.

## The labels

The labels of the x axis are in the order of the rows. To change the order, use `order by` in the query.

The title of the x axis is the name of the column of labels. The title of the y axis is `value`, and the title of the legend is `series`. No option changes them. To change the name of a series, change the name of its column with `as`. [Size, colors and titles](/charts/size-colors-titles/#the-names-of-the-series) shows it.

## A NULL value

A column that holds a NULL is not a number. usql does not draw it as a series, and it does not tell you. To draw the column, replace NULL with a number in the query:

```sql
select month, coffee, coalesce(cocoa, 0) as cocoa
from sales
order by month_no
\chart
```

## The precision

`prec=N` rounds each value to `N` decimal places before usql draws it. `prec=0`, the default, does not round the values.

```sql
select month, coffee / 30.0 as coffee_per_day from sales order by month_no
\chart title='Cups of coffee each day' prec=1
```

To round to whole numbers, use `round` in the query.
