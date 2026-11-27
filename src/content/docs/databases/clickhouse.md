---
title: ClickHouse
description: Connect usql to ClickHouse, and the URL schemes and aliases that open it.
sidebar:
  order: 8
---

## Connect

```sh
usql ch://user:pass@host/dbname
```

The driver connects to the HTTP interface of ClickHouse, on port 8123. With `tls=true`, it uses HTTPS on port 8443.

## Schemes and aliases

::schemes{name="clickhouse"}
