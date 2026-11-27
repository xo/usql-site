---
title: Configure usql
description: Set usql up with config.yaml and init.sql, change a setting with \uset, change the prompt, and find the history file.
sidebar:
  order: 5
---

usql reads two files when it starts. `config.yaml` holds settings, print options and named connections. `init.sql` holds commands that usql runs before the first prompt. Both are in the configuration directory of usql.

::cast[A config.yaml file sets the prompt, the text of NULL and a named connection, and init.sql sets a variable.]{name="configure"}

## Where the files are

| System | Configuration directory |
| --- | --- |
| Linux and other Unix systems | `$XDG_CONFIG_HOME/usql`, which is `~/.config/usql` when `XDG_CONFIG_HOME` is not set |
| macOS | `~/Library/Application Support/usql` |
| Windows | `%AppData%\usql` |

macOS does not read `XDG_CONFIG_HOME`. `\? variables` prints the path of `config.yaml` on your system.

The `USQL_CONFIG` environment variable names another `config.yaml`, and the `USQLRC` environment variable names another `init.sql`.

## config.yaml

`config.yaml` is a [YAML](https://yaml.org) file. Each section of the file is a map of names and values:

```yaml
usql:
  prompt1: "%033[1;32m%/%033[0m%R%# "
print:
  null: (none)
connections:
  books: sq:books.db
```

| Section | What it holds | The command that sets it in a session |
| --- | --- | --- |
| `usql` | the settings of usql | `\uset` |
| `print` | the print options, which change the output | `\pset` |
| `connections` | the [named connections](/guides/named-connections/) | `\cset` |

The file has other sections too. The [README of usql](https://github.com/xo/usql#readme) describes each of them.

usql prints nothing when it reads the file. If usql cannot use a value, it prints an error with the file, the line and the name of the setting. Then it reads the other settings.

## init.sql

`init.sql` is a usql script, as `~/.psqlrc` is for `psql`. usql runs it after it opens the database of the command line, and before `-c`, `-f` and the first prompt. Put meta commands in it, such as `\set`, `\pset` and `\uset`:

```sql
\set since 1965
```

A SQL statement in `init.sql` runs on the database of the command line. If you start usql with no database, the statement fails with `error: not connected`. If you sometimes start usql with no database, keep SQL out of `init.sql`.

To start usql without `init.sql`, give `-X`:

```sh
usql -X books.db
```

`-X` does not skip `config.yaml`.

## Settings

A setting changes how usql works, such as the pager, the editor or the prompt. `\uset` with no argument lists every setting and its value. `\uset` with a name and a value changes the setting for the session:

```sql
\uset on_error_stop on
```

`\config` shows and sets every kind of setting by its full name, as in `usql.prompt1` or `print.null`:

```sql
\config print.null
```

usql does not save a change that `\uset`, `\pset` or `\config` makes. To keep a change, put it in `config.yaml` or in `init.sql`.

`\? variables` describes each setting. The [README of usql](https://github.com/xo/usql#readme) has the full list.

## Which value wins

usql applies the values in this order. A later value wins over an earlier one:

1. The defaults of usql. Some come from the environment, such as `PAGER` and `EDITOR`.
2. `config.yaml`.
3. The environment variables of usql, such as `USQL_PAGER` and `USQL_EDITOR`.
4. The options of the command line.
5. `init.sql`.
6. The commands that you type in the session.

## The prompt

The settings `prompt1` and `prompt2` hold the prompt. usql shows `prompt1` on the first row of a statement and `prompt2` on the next rows. The default for both is `%S%N%m%/%R%# `.

A `%` and a letter in the prompt is an escape. usql puts a value in its place. These escapes make the default prompt:

| Escape | What usql shows |
| --- | --- |
| `%S` | the short name of the driver and a colon, such as `pg:` |
| `%N` | the user name and `@`, when the URL has a user |
| `%m` | the host name, up to the first dot |
| `%/` | the database name |
| `%R` | `=` for a new statement, and `-` while a statement continues |
| `%#` | `>`, or `~` in a transaction |

`%033` puts the character with the octal code 033 in the prompt. That character is Escape, which starts a color code for the terminal. This prompt shows the database name in bold green:

```sql
\uset prompt1 '%033[1;32m%/%033[0m%R%# '
```

## The history

At a terminal, usql saves each statement and command that you type in the file `history` in its state directory:

| System | History file |
| --- | --- |
| Linux and other Unix systems | `$XDG_STATE_HOME/usql/history`, which is `~/.local/state/usql/history` when `XDG_STATE_HOME` is not set |
| macOS | `~/Library/Application Support/usql/history` |
| Windows | `%LocalAppData%\usql\history` |

The `USQL_HISTORY` environment variable names another file. usql does not save the history of `-c` and `-f`.

`\s` prints the history. `\hist 10` prints the last 10 entries, with their numbers. [Line editing and completion](/guides/line-editing/) shows how to search the history.
