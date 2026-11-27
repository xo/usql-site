---
title: Get help
description: The help that usql prints, and where to ask a question.
sidebar:
  order: 4
---

usql prints help about its own commands at the prompt and on the command line. For a question that the help does not answer, ask on GitHub.

## At the prompt

`help` prints a short list of where to start:

```txt
sq:books.db=> help
You are using usql, the universal command-line interface for SQL databases.
Type:  \copyright for distribution terms
       \? for help with usql commands
       \g or terminate with semicolon to execute query
       \q to quit
```

`\?` lists every meta command, in groups, with its arguments. A meta command is a command that starts with a backslash. Add a word to `\?` for other help:

| Command | What it shows |
| --- | --- |
| `\?` | the meta commands |
| `\? commands` | the meta commands, with their aliases and the aliases of `config.yaml` |
| `\? options` | the options of the command line |
| `\? variables` | the settings, the variables and the environment variables, and the path of `config.yaml` on your system |
| `\? styles` | the styles of [syntax highlighting](/guides/syntax-highlighting/), and the style in use |

usql shows long help in the pager. The pager is the program that `PAGER` names, or `less` when `PAGER` is not set. In `less`, press `q` to go back to the prompt.

usql has no `\h` or `\help`, so it gives no help on SQL statements. Read the manual of your database for the SQL.

## On the command line

```sh
usql --help
```

usql prints its options and exits.

```sh
usql --version
```

usql prints its version. Give the version when you ask a question or report a problem.

## Ask a question

Ask a question in the [discussions of usql on GitHub](https://github.com/xo/usql/discussions). Search the discussions first. Another person probably asked the same question.

If usql does something wrong, open an issue in the [issues of usql](https://github.com/xo/usql/issues). Give the version of usql, the database and its version, the commands that you typed, and what usql printed.

The [README of usql](https://github.com/xo/usql#readme) describes every feature of usql.
