---
title: Coming from psql
description: What works in usql as it does in psql, what usql adds, and what it does not have.
sidebar:
  order: 3
---

usql follows PostgreSQL's `psql`. If you know `psql`, you already know most of usql. This page lists what stays the same, what usql adds, and what is missing.

## What works as in psql

- A statement ends at a `;`, or at a meta command such as `\g`. A meta command is a command that starts with a backslash. usql reads its input as `psql` does, on every database.
- The prompt shows `=>` for a new statement and `->` while a statement continues on another row.
- `help` shows a short help, and `quit` or `exit` leaves usql. Ctrl-C clears the line, and Ctrl-D at an empty prompt leaves usql.
- The meta commands have the names and the arguments of `psql`. For example, `\d`, `\dt`, `\x`, `\pset`, `\g`, `\gx`, `\gset`, `\gexec`, `\watch`, `\e`, `\i`, `\o`, `\timing`, and `\if` with `\elif`, `\else` and `\endif`.
- `\set` sets a variable, and `:name`, `:'name'` and `:"name"` put its value into a statement.
- A backtick in the argument of a meta command runs a shell command, as in `psql`.
- `-c` runs one command, `-f` runs a file, and `-X` skips the startup file.

## What usql adds

- Many databases. usql opens PostgreSQL, MySQL, SQLite, SQL Server, Oracle, DuckDB and many more from a URL. `\c` opens a database of another kind in the same session. [Connect to a database](/guides/connect/) explains the URLs.
- `\copy` between databases. usql's `\copy` runs a query on one database and inserts the rows into a table of another database. It does not read or write a file as the `\copy` of `psql` does.
- `\dump` writes a database as a usql script. The script can make the database again, also on a database of another kind.
- `\chart` draws the result of a query as a chart in a terminal that draws images. [Your first chart](/charts/first-chart/) shows how.
- [Named connections](/guides/named-connections/), with `\cset` and the `connections:` section of `config.yaml`.
- [Completion](/guides/line-editing/) of meta commands, SQL words and names, in a numbered menu.
- [Syntax highlighting](/guides/syntax-highlighting/) of what you type, in the grammar of the database.

## Settings are not variables

In `psql`, `\set` sets both your own variables and the special variables, such as `ON_ERROR_STOP`. In usql, `\set` sets only your own variables. The behavior of usql comes from settings, which `\uset` sets, and their names are in lower case:

```sql
\uset on_error_stop on
```

Do not use `\set ON_ERROR_STOP on` in usql. usql accepts it, but it is an ordinary variable and a script does not stop at an error. [Configure usql](/guides/configure/) explains the settings.

## Files that usql reads in place of psql's

| psql | usql |
| --- | --- |
| `~/.psqlrc` | `init.sql` in the configuration directory of usql |
| `~/.psql_history` | `history` in the state directory of usql |
| `~/.inputrc` | the same file. A `$if usql` block applies only to usql |

[Configure usql](/guides/configure/) says where these directories are on each system.

## What usql does not have

- `\h` and `\help`. usql has no help for SQL statements. `\?` lists the meta commands.
- The `\copy` of `psql`, which copies between a table and a file.
- The meta commands of the extended query protocol, such as `\parse`, `\bind_named` and `\startpipeline`.
- `\lo_import` and `\lo_export`.
- `\encoding`. usql reads a script as UTF-8. The setting `encoding` names another encoding for a script.

The [README of usql](https://github.com/xo/usql#readme) lists every meta command that usql has.
