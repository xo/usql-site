---
title: First steps
description: Open a database, list its tables, run a query and change the output format.
sidebar:
  order: 2
---

This page follows the recording below. It uses a SQLite file, so you need no database server.

::cast[Open a SQLite file, list its tables, run a query, and print the result as CSV.]{name="tour"}

## Open a database

Give usql the database as a URL, or as the path of a file:

```sh
usql books.db
```

usql opens a file that exists with the SQLite or DuckDB driver. For a server, give its URL, such as `pg://user@localhost/books`. [Connect to a database](/guides/connect/) explains the URLs.

The prompt starts with the short name of the driver and the database, as in `sq:books.db=>`.

## List the tables

`\dt` lists the tables. `\d books` describes the table `books`. `\?` lists every backslash command.

## Run a query

Type SQL and end it with a semicolon. usql runs the statement when you press Enter after the semicolon:

```sql
select a.name, count(*) as books
from authors a join books b on b.author_id = a.id
group by a.name order by books desc;
```

A statement can span lines. The prompt changes to `->` until the statement ends.

## Change the output format

`\pset format` changes how usql prints a result. The formats include `aligned`, `csv`, `json`, `html` and `markdown`:

```sql
\pset format csv
```

## Leave

`\q` closes the connection and leaves usql. Ctrl-D does the same at an empty prompt.
