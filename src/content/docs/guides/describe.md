---
title: Describe a database
description: List the tables, views, indexes, functions, schemas and databases with the \d commands, and describe one table.
sidebar:
  order: 13
---

The describe commands list what a database holds and describe one object, such as a table. They are the meta commands that start with `\d`, and `\l`. They work as they do in psql, on every database that usql has the metadata for. Metadata is the data about the tables, columns and indexes of a database.

::cast[List the objects of a DuckDB file, describe a table, list its views, indexes and functions, and see the error of a command that DuckDB has no metadata for.]{name="describe"}

## List the objects

| Command | What it lists                           |
| ------- | --------------------------------------- |
| `\d`    | the tables, the views and the sequences |
| `\dt`   | the tables                              |
| `\dv`   | the views                               |
| `\di`   | the indexes, each with its table        |
| `\df`   | the functions                           |
| `\dn`   | the schemas                             |
| `\l`    | the databases                           |

Each command takes a pattern, to list only the names that match. `*` matches any text, and a dot separates the schema from the name:

```sql
\dt prod*
\dt main.*
```

## Describe one object

`\d` with a name describes that object. For a table, it prints the columns, their types, which columns can be `NULL`, the default values, the indexes and the foreign keys:

```sql
\d orders
```

`\d` describes a view in the same way.

## More detail

A letter after the command changes what it prints:

| Letter | What it does                                                    |
| ------ | --------------------------------------------------------------- |
| `+`    | prints more columns, such as the description of each object     |
| `S`    | lists the objects of the system too, such as built-in functions |
| `x`    | prints the result in expanded mode, one column on each line     |

For example, `\dt+` lists the tables with more detail, and `\dfS` lists the built-in functions too.

## SQLite and DuckDB

A SQLite file has one schema, `main`. `\l` lists `main` as its database. SQLite keeps no functions in a database file, so `\df` lists functions that the SQLite library has, such as those of its extensions.

A DuckDB file can hold more schemas. A table that you make with no schema goes into `main`. `\l` lists the database of the file, by the name of the file. `\df` lists the macros that you made.

## Where the metadata comes from

usql reads the metadata with the library dbmeta. A describe command works only where dbmeta has the metadata for that database. If it does not, usql prints an error and runs nothing:

```txt
error: \du not supported by duckdb driver
```

For example, neither SQLite nor DuckDB has users or roles, so `\du` does not work on them. `\dT`, which lists data types, works on DuckDB and not on SQLite.

`\?` lists every describe command. Most of them are for PostgreSQL and the databases like it.
