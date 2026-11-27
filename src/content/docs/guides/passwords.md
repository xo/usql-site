---
title: Passwords
description: Give usql the password of a database without showing it to other users, with the password prompt, -W, -w, USQL_DSN and named connections.
sidebar:
  order: 3
---

usql takes the password of a database from the URL, or asks for it at a prompt. This page shows where a password can go, and who can read it there.

::cast[With -w, usql stops when the server refuses the login. Without -w, usql asks for the password, and the prompt does not show what is typed.]{name="passwords"}

## A password in a URL

A URL holds the password after the user name and a `:`:

```sh
usql pg://app:secret@localhost/app
```

Every user of the machine can read the command line of a running program, for example with `ps`. Your shell also writes the command to its history file. Do not give a URL with a password on the command line of a machine that other people use.

If the password holds a character that has a meaning in a URL, such as `@`, `:` or `/`, encode it as `%` and its hex code. The password `a@b` is `a%40b` in a URL:

```sh
usql pg://app:a%40b@localhost/app
```

usql does not show a password that it knows. `\conninfo` and `\cset` show it as `xxxxx`. But if you type a URL with a password at the prompt of usql, for example with `\c`, the history file of usql keeps the line. Only you can read the history file.

## The password prompt

If the server refuses the login, usql prints the error and asks for the password:

```txt
$ usql pg://app@localhost/app
error: postgres: 28P01: password authentication failed for user "app"
Enter password:
```

The prompt does not show what you type. usql connects again with the password that you typed. `\c` asks in the same way.

usql asks only when all of these are true:

- usql reads from a terminal. In a pipe or a job with no terminal, usql prints the error and stops.
- The driver tells usql that the server refused the password. The drivers of PostgreSQL and MySQL do.
- You did not give `-w`.

## Ask before the first try with -W

`-W` (`--password`) makes usql ask for the password before it connects:

```sh
usql -W pg://app@localhost/app
```

`-W` also works with a URL in `USQL_DSN`:

```sh
USQL_DSN=pg://app@localhost/app usql -W
```

:::caution
In this release, `-W` does not work with a named connection. `usql -W app` prints `error: named connection is not a url`. Leave out `-W`, and usql asks for the password when the server refuses the login.
:::

## Never ask with -w

`-w` (`--no-password`) makes usql never ask. If the server refuses the login, usql prints the error and stops. Give `-w` to a script that runs at a terminal, so that it stops and does not wait for a password:

```sh
usql -w pg://app@localhost/app -f report.sql
```

## The USQL_DSN environment variable

When the command line names no database, usql opens the URL in the `USQL_DSN` environment variable. Other users cannot read the environment of your program, as they can read its command line:

```sh
USQL_DSN="pg://app:$PASSWORD@localhost/app" usql
```

`USQL_DSN` holds a URL or the name of a connection. usql removes the variable from its own environment after it reads it, so the shell of `\!`, the pager and the editor do not get the password. `\c` does not read the variable.

usql has no environment variable for the password alone. The PostgreSQL driver reads `PGPASSWORD`, as [the page of PostgreSQL](/databases/postgres/) says.

## Named connections in config.yaml

A [named connection](/guides/named-connections/) in `config.yaml` keeps the password off the command line. The connection holds the password in its URL, or in the `password` key of its map:

```yaml
connections:
  app: pg://app:secret@localhost/app
  shop:
    protocol: mysql
    username: shop
    password: secret
    hostname: localhost
    database: shop
```

Then give the name in place of the URL:

```sh
usql app
```

Leave the password out of a connection, and usql asks for it when the server refuses the login. The recording at the top of this page does that.

`config.yaml` is a plain text file, and usql does not check who can read it. If the file holds a password, make sure that only you can read it. On Linux:

```sh
chmod 600 ~/.config/usql/config.yaml
```

[Configure usql](/guides/configure/) says where the file is on other systems.

## No password file

usql reads no password file of its own. Earlier releases read `~/.usqlpass`, and this release does not. Put each password in a named connection in `config.yaml`, or let usql ask for it.
