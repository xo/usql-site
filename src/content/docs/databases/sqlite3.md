---
title: SQLite
description: Open a SQLite file with usql, and the URL schemes and aliases that open it.
sidebar:
  order: 4
---

## Connect

A path opens a SQLite file that exists:

```sh
usql books.db
```

To make a new file, name the scheme, such as `sq:`:

```sh
usql sq:books.db
```

## Schemes and aliases

::schemes{name="sqlite3"}
