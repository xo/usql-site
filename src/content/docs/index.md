---
title: usql
description: A universal command line interface for SQL databases, modeled on psql.
template: splash
hero:
  tagline: One command line for PostgreSQL, MySQL, SQLite, SQL Server, Oracle, DuckDB and many other databases. It works like psql.
  image:
    file: ../../assets/logo.svg
  actions:
    - text: Install usql
      link: /start/install/
      icon: right-arrow
    - text: Take the tour
      link: /start/first-steps/
      variant: minimal
---

::cast[A first session: open a SQLite file, list its tables, and run a query.]{name="tour"}

usql opens a database from a URL, such as `pg://user@host/db` or `my://localhost/db`, or from a file on disk. The prompt takes SQL and psql's backslash commands, such as `\d`, `\copy`, `\set` and `\pset`. It also has commands that psql does not have: `\copy` between two databases, `\dump`, and `\chart`, which draws a chart in the terminal.
