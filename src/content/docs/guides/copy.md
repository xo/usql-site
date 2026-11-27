---
title: Copy data between databases
description: Copy the result of a query from one database into a table of another database with \copy, and import a CSV file.
sidebar:
  order: 14
---

`\copy` runs a query on one database and writes its rows into a table of another database. The two databases can use different drivers, such as SQLite and DuckDB. This page explains how to copy rows, and how to import a CSV file.

`\copy` in usql is not the `\copy` of psql. It does not read or write a file in the format of the PostgreSQL `COPY` command.

::cast[Copy a table from a SQLite file into a DuckDB file, import a CSV file into the DuckDB file, and query the result.]{name="copy"}

## Copy rows

`\copy` takes four parameters:

```txt
\copy SRC DST QUERY TABLE
```

| Parameter | What it is                                                         |
| --------- | ------------------------------------------------------------------ |
| `SRC`     | the URL of the source database, where the query runs               |
| `DST`     | the URL of the destination database, which holds the table         |
| `QUERY`   | the query to run on the source database                            |
| `TABLE`   | the table to write the rows into, with an optional list of columns |

This command copies the products of a SQLite file into a DuckDB file:

```sql
\copy sq:shop.db dk:warehouse.duckdb 'select id, name, price from products' products
```

usql prints `COPY` and the number of rows that it wrote:

```txt
COPY 3
```

`\copy` opens its own connections to the two databases. You do not have to connect to a database first, and `\copy` does not change the connection that is open. [Connect to a database](/guides/connect/) explains the URLs.

## The rules for the table and the query

- The table must exist in the destination database before you copy. `\copy` does not make it.
- The query must return one column for each column of the table. To write fewer columns, name them after the table:

  ```sql
  \copy sq:shop.db dk:warehouse.duckdb 'select name, price from products' 'products(name, price)'
  ```

- Put quotes around the query, and around a table with a list of columns, because they hold spaces.
- `\copy` does not change the type of a value. If a column of the table needs another type, cast the value in the query, with the functions of the source database.

## Keep the URLs in variables

To copy more than one table, keep the URLs in variables. `\copy` puts the values of variables into its parameters, as other meta commands do:

```sql
\set src sq:shop.db
\set dst dk:warehouse.duckdb
\copy :src :dst 'select id, name, price from products' products
```

[Variables and conditionals](/guides/variables/) explains the variables.

## Import a CSV file

To import a CSV file, use a source database that can read the file as a table. DuckDB reads a CSV file with its function `read_csv`. `dk:` with no file name opens a DuckDB database in memory, which can read the file:

```sql
\copy dk: dk:warehouse.duckdb "select * from read_csv('customers.csv')" customers
```

The destination does not have to be DuckDB. The same command writes into a SQLite file when `DST` is `sq:shop.db`.

The README of usql shows another way, with the `csvq` driver, which reads each CSV file of a directory as a table. See [Importing Data from CSV](https://github.com/xo/usql#importing-data-from-csv).

## Large tables

`\copy` reads each row through usql. For a very large table, use the tools of the database, such as its bulk loader. They are faster.
