---
title: Databases
description: Every database that usql connects to, with the URL schemes and aliases that open it.
sidebar:
  order: 1
  label: All databases
---

usql opens a database from a URL. The scheme of the URL, the part before the `:`, names the driver. Each driver also has aliases, and every alias opens the same driver. [Connect to a database](/guides/connect/) explains the parts of a URL.

The table below lists every database of the release of usql that the site documents. It comes from the scheme registry of [dburl](https://github.com/xo/dburl), the package that reads the URLs of usql. A database with a link has a page with notes.

::databases
