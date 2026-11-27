---
title: Run queries and scripts
description: The query buffer, the meta commands that run a query, and scripts that usql runs from a file or from the command line.
sidebar:
  order: 10
---

This page explains how usql collects a statement, the meta commands that run it, and how usql runs a script. A meta command is a command that starts with a backslash, such as `\g`. usql runs a meta command itself, and does not send it to the database.

::cast[Run a query on two lines, print the query buffer, show a row in expanded mode, keep a result in variables, run a query each second, and run one query from the shell.]{name="run-queries"}

## The query buffer

usql keeps the statement that you type in the query buffer. A statement can span lines. While the buffer holds text, the prompt ends with `->` in place of `=>`.

usql runs the statement when a line ends with a semicolon. It also runs it when you type `\g`:

```sql
select name, price from products
where price > 5
\g
```

A semicolon inside a string, a quoted name or a comment does not end the statement.

These meta commands work on the query buffer:

| Command | What it does                                                     |
| ------- | ---------------------------------------------------------------- |
| `\p`    | prints the query buffer, or the last query that ran              |
| `\r`    | clears the query buffer                                          |
| `\e`    | opens the query buffer, or the last query that ran, in an editor |

## Edit a query

`\e` writes the query to a temporary file and opens the file in your editor. When you close the editor, usql reads the file again. If the text ends with a semicolon, usql runs it. If it does not, the text stays in the query buffer.

usql uses the `editor` setting. Its default comes from the `EDITOR` or the `VISUAL` environment variable. If neither is set, usql uses `vi`. To use another editor, set it with `\uset`:

```sql
\uset editor 'code --wait'
```

## Run a query in another way

`\g` can send the result to a file, or to a command after `|`:

```sql
select name from products order by name \g names.txt
select name from products \g |wc -l
```

`\gx` runs the query and prints the result in expanded mode. Expanded mode prints each row as a list of its columns, one column on each line. It helps when a row is too wide for the terminal:

```sql
select * from products where id = 2 \gx
```

`\gset` runs a query that returns one row. It keeps each column in a variable with the name of the column. Add a prefix after `\gset` to put it before each name:

```sql
select count(*) as n, sum(qty) as items from orders \gset
\echo :n orders, :items items
select count(*) as n from orders \gset orders_
\echo :orders_n
```

If the query returns no row, or more than one row, `\gset` prints an error and sets no variable. [Variables and conditionals](/guides/variables/) explains the variables.

## Run a query again and again

`\watch` runs the query again and again, and prints the time before each result. Give the interval in seconds. The default is 2 seconds:

```sql
select count(*) as orders from orders \watch 1
```

To stop it, press Ctrl-C.

## Run a script

A script is a file of SQL statements and meta commands. In usql, `\i` runs a script in the open session:

```sql
\i report.sql
```

On the command line, `-f` runs a script and then exits. `-c` runs one command and then exits:

```sh
usql -f report.sql shop.db
usql -c 'select sum(qty) from orders' shop.db
```

usql also reads a script from a pipe:

```sh
echo 'select count(*) from orders' | usql shop.db
```

In a script, the last statement runs even when it has no semicolon.

usql runs the startup script `init.sql` before a script, if the file exists. To start usql without it, add `-X`:

```sh
usql -X -f report.sql shop.db
```

## Errors in a script

If a statement fails, usql prints the error and runs the rest of the script. To stop the script at the first error, put this line at the top of the script:

```sql
\uset on_error_stop on
```

## More commands

`\?` lists every meta command. The [README of usql](https://github.com/xo/usql#readme) describes the command line options.
