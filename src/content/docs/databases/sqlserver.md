---
title: SQL Server
description: Connect usql to Microsoft SQL Server, and the URL schemes and aliases that open it.
sidebar:
  order: 5
---

## Connect

```sh
usql ms://user:pass@host/dbname
```

To connect to a named instance, put the instance before the database:

```sh
usql ms://user:pass@host/instancename/dbname
```

## Schemes and aliases

::schemes{name="sqlserver"}
