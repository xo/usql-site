---
title: Connect to a database
description: The URLs that usql opens, the short names of the drivers, and named connections.
sidebar:
  order: 1
---

usql opens a database from a URL. The command line, `\connect` (or `\c`) and `\copy` all take the same URLs. A URL is also called a DSN, a data source name.

## The forms of a URL

```txt
driver+transport://user:pass@host/dbname?opt1=a&opt2=b
driver:/path/to/file
/path/to/file
name
```

| Part | What it is |
| --- | --- |
| `driver` | the name of the driver, or one of its aliases |
| `transport` | `tcp`, `udp`, `unix`, or the name of a driver for ODBC and ADODB |
| `user`, `pass` | the user name and the password |
| `host` | the name of the host |
| `dbname` | the database name, instance, or service name |
| `?opt1=a&...` | options of the driver |
| `/path/to/file` | a path on disk |
| `name` | a connection that `\cset` or the configuration file named |

## Short names

Every driver has a short name of two letters, which is usually the start of its name. For example, `pg` is PostgreSQL, `my` is MySQL, `ms` is SQL Server, `or` is Oracle and `sq` is SQLite:

```sh
usql pg://postgres@localhost/postgres
usql my://root@localhost/mysql
usql sq:books.db
```

## Paths on disk

If a URL has no driver, usql reads it as a path. usql opens a Unix domain socket with the MySQL driver, a directory with the PostgreSQL driver, and a file with the SQLite or DuckDB driver.

## Defaults

Most parts of a URL are optional. usql fills in what it can. This command connects to PostgreSQL as the local user, through the socket in `/var/run/postgresql`:

```sh
usql pg://
```

## Keep a password off the command line

Every user of a machine can read the command line of a program. If the URL holds a password, give it in the `USQL_DSN` environment variable:

```sh
USQL_DSN="pg://app:$PASSWORD@localhost:15432/app" usql
```

usql reads `USQL_DSN` only when the command line names no database. It removes the variable from its own environment after it reads it, so the shell of `\!`, the pager and the editor do not get the password.
