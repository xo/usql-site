---
title: Charts in a script
description: Run a script of charts with \i at the prompt, and why usql -f draws no chart.
sidebar:
  order: 5
---

A script can hold `\chart` as it holds any other meta command. Run the script with `\i` at the prompt, and usql draws each chart. A script that runs with `-f` draws no chart.

This script, `charts.sql`, writes two charts to SVG files:

```sql
-- The sales of each month, and the cups of each hour.
select month, coffee, tea, cocoa from sales order by month_no
\chart title='Cups sold' file=sales.svg
select hour, cups from hours order by hour
\chart title='Cups by hour' file=hours.svg
```

The script uses the tables of the [gallery](/charts/gallery/#the-data).

## Run a script with \i

`\i` runs a script at the prompt of usql:

```sql
\i charts.sql
```

::cast[Show charts.sql, run it with \i, and list the two SVG files that it writes.]{name="chart-script"}

usql writes `sales.svg` and `hours.svg`:

![A bar chart of the cups of coffee, tea and cocoa sold each month.](/charts/chart-script/sales.svg)

![A line chart of the cups sold each hour from 7 to 17.](/charts/chart-script/hours.svg)

In a terminal that draws images, a script with the same `\chart` lines without `file=` draws each chart in place.

## Why -f draws no chart

When usql runs a script from `-f`, or a command from `-c`, it turns the terminal graphics off, so that the escape codes of an image do not go to a pipe or a file. `\chart` then stops with an error:

```sh
usql -f charts.sql cafe.db
```

```text
error: sqlite3: \chart: graphics not supported in terminal
```

To turn the graphics on for `-f`, choose the protocol of your terminal on the command line with `-U`, such as `kitty`, `iterm` or `sixel`:

```sh
usql -U term_graphics=kitty -f charts.sql cafe.db
```

::cast[Run charts.sql with -f, which fails at each chart, and run it again with -U term_graphics=kitty, which writes the files.]{name="chart-script-file"}

usql then draws each chart without `file=` to the standard output. Do this only when the output goes to your terminal.

In the current release of usql, `\chart` needs graphics also for `file=`. So a script that writes its charts to files with `-f` needs `-U term_graphics=kitty` too, as the recording shows. With `file=`, usql writes the file and draws nothing. [Terminal graphics](/guides/terminal-graphics/) says more about the protocols.
