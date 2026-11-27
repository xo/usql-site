---
title: Dump a database
description: Write a database as a usql script with \dump, choose what goes into it, and write it for another database.
sidebar:
  order: 15
---

`\dump` writes a database as a usql script. The script holds the schema and the rows. When usql runs the script, it makes the database again. A dump can also be for another database, such as a SQLite file that you load into DuckDB.

::cast[Dump the schema of one table, write a dump of a SQLite file for DuckDB, load it into a new DuckDB file, and query it.]{name="dump"}

## Write a dump

`\dump` with no parameter writes the dump of the open database to the output. To write it to a file, give the name of the file:

```sql
\dump shop.sql
```

usql writes a file whole or not at all, and only its owner can read it.

To send it to a command, put `|` before the command:

```sql
\dump |gzip > shop.sql.gz
```

## Load a dump

A dump is a usql script. Run it with `-f`, and give the URL of the database to load it into:

```sh
usql -f shop.sql sq:copy.db
```

The script turns on the `on_error_stop` setting, so that it stops at the first error.

## Choose what goes into the dump

Options come before the file name. Each option starts with `-`:

| Option                   | What it does                                         |
| ------------------------ | ---------------------------------------------------- |
| `-schema-only`           | dumps the schema and no rows                         |
| `-data-only`             | dumps the rows and no schema                         |
| `-table PATTERN`         | dumps only the tables that match the pattern         |
| `-exclude-table PATTERN` | dumps all tables except those that match the pattern |

A pattern is a pattern of SQL `LIKE`: `%` matches any text. You can give `-table` and `-exclude-table` more than once:

```sql
\dump -schema-only -table orders
\dump -data-only -exclude-table 'log%' data.sql
```

`\?` shows every option of `\dump`, and the [README of usql](https://github.com/xo/usql#dumping-a-database) describes them.

## Dump for another database

`-target` writes the dump for another database. usql changes the types and the statements of the schema to those of the target:

```sql
\dump -target duckdb shop-duckdb.sql
```

Then load the dump into a new DuckDB file:

```sh
usql -f shop-duckdb.sql duckdb:shop.duckdb
```

The target is the name of a driver, or one of its aliases, such as `duckdb`, `sqlite3`, `postgres` or `mysql`.

usql cannot change every statement. If it cannot change one, such as the text of a view, it prints a warning, and the dump holds the statement as a comment. Read the warnings, and make those objects by hand.

## A dump that agrees with itself

usql reads the schema and the rows in one transaction, so that they agree. If you started a transaction with `\begin`, `\dump` reads in that transaction, and the dump holds the changes that are not committed yet. [Transactions](/guides/transactions/) explains `\begin`.
