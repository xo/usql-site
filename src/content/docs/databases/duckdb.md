---
title: DuckDB
description: Open a DuckDB file with usql, and the URL schemes and aliases that open it.
sidebar:
  order: 7
---

## Connect

A path opens a DuckDB file that exists:

```sh
usql analytics.duckdb
```

To make a new file, name the scheme, such as `dk:`:

```sh
usql dk:analytics.duckdb
```

## Schemes and aliases

::schemes{name="duckdb"}
