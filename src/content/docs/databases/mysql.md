---
title: MySQL
description: Connect usql to MySQL and MariaDB, and the URL schemes and aliases that open them.
sidebar:
  order: 3
---

## Connect

```sh
usql my://user:pass@host/dbname
```

Without a host, usql connects to the server on `localhost`:

```sh
usql my://
```

## Schemes and aliases

::schemes{name="mysql"}
