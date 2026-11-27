---
title: Format the output
description: Print a result as a table, CSV, JSON or another format, show one row at a time, use the pager, and choose how times look.
sidebar:
  order: 12
---

usql prints a result as a table by default. A print option changes how it prints. `\pset` sets a print option, and `\pset` with no name lists every print option and its value. This page explains the print options that you need most.

::cast[Print a result as a table and as CSV, show one row in expanded mode, and change the time format.]{name="format-output"}

## Choose a format

`\pset format` sets the format of the output:

```sql
\pset format csv
select item, qty from orders;
```

These are the formats that you need most:

| Format      | What it prints                                          |
| ----------- | ------------------------------------------------------- |
| `aligned`   | a table with aligned columns. This is the default       |
| `unaligned` | the values of a row on one line, with `\|` between them |
| `csv`       | comma-separated values, with a header line              |
| `json`      | an array with one object for each row                   |
| `markdown`  | a table in Markdown                                     |
| `html`      | a table in HTML                                         |

usql has more formats. `\? variables` lists them with the other print options, and the [README of usql](https://github.com/xo/usql#display-formatting-print-variables) describes the print options.

`\a` switches between `aligned` and `unaligned`.

## Show one row at a time

`\x` turns expanded mode on and off. In expanded mode, usql prints each row as a list of its columns, one column on each line. A wide row is easier to read this way:

```sql
\x
select * from orders where id = 3;
```

`\x auto` uses expanded mode only when a result is too wide for the terminal.

To print only one result in expanded mode, end the query with `\gx` in place of the semicolon.

## Print the rows only

`\t` turns off the names of the columns and the row count at the end, and turns them on again. It helps when another program reads the output.

## Show NULL values

usql prints a `NULL` value as an empty cell. To tell a `NULL` value from an empty string, set the text that it prints:

```sql
\pset null '(null)'
```

## Lines and borders

`\pset border` sets the border of a table: `0`, `1` or `2`. `\pset linestyle unicode` draws the lines with box characters in place of `-`, `|` and `+`:

```sql
\pset border 2
\pset linestyle unicode
```

## The pager

A pager is a program, such as `less`, that shows long text one screen at a time. If a result is taller or wider than the terminal, usql sends it to the pager.

The `pager` print option controls this:

| Value    | What usql does                                           |
| -------- | -------------------------------------------------------- |
| `on`     | uses the pager when the result does not fit the terminal |
| `off`    | never uses the pager                                     |
| `always` | uses the pager for every result                          |

```sql
\pset pager off
```

The pager program is a setting of usql. Its default comes from the `PAGER` environment variable. If `PAGER` is not set, usql uses `less` or `more`. To use another program, set it with `\uset`:

```sql
\uset pager 'less -S'
```

## Time formats

usql prints a date or a time in the RFC3339Nano format by default, such as `2026-11-02T09:15:00Z`. `\pset time` sets another format:

```sql
\pset time '2006-01-02 15:04'
select id, placed_at from orders;
```

The format is a layout of the Go programming language. You write the reference time, Monday, January 2, 2006 at 15:04:05, in the form that you want to see. For example, `'Jan 2, 2006'` prints `Nov 2, 2026`.

`\pset time` also takes the name of a format of Go, such as `RFC1123` or `Kitchen`. The [README of usql](https://github.com/xo/usql#time-constants) lists the names.

`\pset timezone` sets the time zone that usql prints a time in, such as `UTC` or `Europe/Berlin`.

## Set the print options at each start

To use the same print options in each session, put the `\pset` commands in the startup script `init.sql`. The [README of usql](https://github.com/xo/usql#startup-script) says where the file is.
