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

## Schemes and aliases

::schemes{name="postgres"}
