---
title: Gallery
description: Each chart of the chart pages, with the query and the \chart line that make it.
sidebar:
  order: 9
---

Each chart of the chart pages is here, with the query and the `\chart` line that make it. Each chart is the SVG file that `file=` writes. In a terminal that draws images, the same `\chart` without `file=` draws the chart in place.

## The data

The charts use two tables of a small cafe: the cups sold each month, and the cups sold each hour of a day. Save this SQL as `cafe.sql`:

```sql
create table sales (
  month_no integer primary key,
  month text not null,
  coffee integer not null,
  tea integer not null,
  cocoa integer not null
);
insert into sales values
  (1, 'Jan', 1240, 860, 610), (2, 'Feb', 1180, 820, 560),
  (3, 'Mar', 1310, 790, 420), (4, 'Apr', 1360, 740, 280),
  (5, 'May', 1420, 700, 150), (6, 'Jun', 1490, 680, 90);
create table hours (hour integer primary key, cups integer not null);
insert into hours values
  (7, 38), (8, 92), (9, 71), (10, 54), (11, 47), (12, 66),
  (13, 58), (14, 41), (15, 49), (16, 35), (17, 22);
```

Then make the SQLite file `cafe.db` from it, and open it:

```sh
usql -X -q -f cafe.sql sq:cafe.db
usql cafe.db
```

## Bars of three series

```sql
select month, coffee, tea, cocoa from sales order by month_no
\chart title='Cups sold' file=sales-bar.svg
```

![A bar chart of the cups of coffee, tea and cocoa sold each month from January to June.](/charts/chart-bar-line/sales-bar.svg)

## Lines of three series

```sql
select month, coffee, tea, cocoa from sales order by month_no
\chart title='Cups sold' type=line file=sales-line.svg
```

![A line chart of the cups of coffee, tea and cocoa sold each month. Coffee goes up, and tea and cocoa go down.](/charts/chart-bar-line/sales-line.svg)

## A line of numbers only

```sql
select hour, cups from hours order by hour
\chart title='Cups by hour' file=hours.svg
```

![A line chart of the cups sold each hour from 7 to 17, with the most cups at 8.](/charts/chart-bar-line/hours.svg)

## A title and a subtitle

```sql
select month, coffee as "Coffee", tea as "Tea"
from sales order by month_no
\chart title='Cups sold' subtitle='January to June 2026' file=titled.svg
```

![A bar chart titled Cups sold, with the subtitle January to June 2026, of the cups of coffee and tea sold each month.](/charts/chart-options/titled.svg)

## A smaller size

```sql
select month, coffee as "Coffee", tea as "Tea"
from sales order by month_no
\chart title='Cups sold' size=640x320 file=small.svg
```

![A bar chart of coffee and tea, 640 pixels wide and 320 pixels high.](/charts/chart-options/small.svg)

## One series

```sql
select month, coffee from sales order by month_no
\chart title='Coffee sold' file=coffee.svg
```

![A bar chart titled Coffee sold, with one bar for each month from January to June.](/charts/chart-file/coffee.svg)
