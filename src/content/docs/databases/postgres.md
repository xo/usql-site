---
title: PostgreSQL
description: Connect usql to PostgreSQL, and the URL schemes and aliases that open it.
sidebar:
  order: 2
---

## Connect

```sh
usql pg://user:pass@host/dbname
```

Without a host, usql connects as the local user, through the Unix domain socket of the local server:

```sh
usql pg://
```

To connect without TLS, add `sslmode=disable`:

```sh
usql "pg://user:pass@host/dbname?sslmode=disable"
```

The driver also reads the password from the `PGPASSWORD` environment variable, and from the `~/.pgpass` file or the file that `PGPASSFILE` names, as `psql` does. [Passwords](/guides/passwords/) shows the ways that usql itself takes a password.

::cast[List the tables of the northwind database on PostgreSQL, and run a query over three lines.]{name="postgres"}

## Schemes and aliases

::schemes{name="postgres"}
