---
title: Named connections
description: Give a database URL a short name with \cset or in config.yaml, and open it by that name.
sidebar:
  order: 2
---

A named connection is a short name for the URL of a database. You type the name in place of the URL, on the command line and with `\c`. A name set with `\cset` lasts until usql exits. A name in `config.yaml` is there each time that usql starts.

::cast[Open a SQLite file by a name from config.yaml, name a DuckDB file with \cset, and change between the two.]{name="named-connections"}

## Name a connection in the session

`\cset` takes a name and a URL:

```sql
\cset shelf duckdb:shelf.duckdb
```

`\c` opens the connection by its name:

```sql
\c shelf
```

`\cset` with no argument lists the names and their URLs. It shows a password in a URL as `xxxxx`:

```txt
sq:books.db=> \cset
books 'sq:books.db'
shelf 'duckdb:shelf.duckdb'
```

`\cset` with only a name removes that name. usql does not save a name that `\cset` sets.

## Name a connection in config.yaml

The `connections:` section of `config.yaml` holds the names that usql knows when it starts. [Configure usql](/guides/configure/) says where the file is on your system.

```yaml
connections:
  books: sq:books.db
  app:
    protocol: postgres
    username: app
    hostname: localhost
    port: 5432
    database: app
```

A name has a URL, as `books` has, or a map of the parts of a URL, as `app` has. usql makes the URL `postgres://app@localhost:5432/app` from the map of `app`.

Give the name on the command line in place of the URL:

```sh
usql books
```

## Keep the password safe

`config.yaml` is a plain text file. If a URL in it holds a password, make sure that only you can read the file. On Linux:

```sh
chmod 600 ~/.config/usql/config.yaml
```

On another system, use the path that `\? variables` prints.

## When a name is also a file

usql looks for a name before it reads the argument as a path. If a file in the current directory has the same name as a connection, usql opens the connection. To open the file, give a path such as `./books`.

The [README of usql](https://github.com/xo/usql#readme) describes the other forms of a named connection.
