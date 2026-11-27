---
title: Variables and conditionals
description: Set variables with \set, put them into a query, fill them from a command, and choose the lines that run with \if.
sidebar:
  order: 11
---

A variable holds a value that usql puts into a query or a meta command for you. This page explains how to set a variable, how to use it, and how `\if` runs lines only when a condition is true.

::cast[Set two variables, use them in a query, fill a variable from a shell command, and choose a line with \if.]{name="variables"}

## Set a variable

`\set` sets a variable. Put quotes around a value that has spaces:

```sql
\set col name
\set item 'Coffee beans'
```

`\set` with no name lists every variable. usql sets some variables itself, such as `ROW_COUNT`, which holds the number of rows of the last query.

`\unset` removes a variable:

```sql
\unset item
```

A variable of `\set` is not a setting of usql. The settings, such as the pager and the editor, have their own command, `\uset`.

## Use a variable

usql puts the value of a variable into the query buffer in three forms:

| Form      | What usql writes                                        |
| --------- | ------------------------------------------------------- |
| `:name`   | the value as it is                                      |
| `:'name'` | the value as a string, in single quotes                 |
| `:"name"` | the value as a name, such as a column, in double quotes |

```sql
\set tbl products
\set col name
\set item 'Coffee beans'
select * from :tbl where :"col" = :'item';
```

usql sends this statement to the database:

```sql
select * from products where "name" = 'Coffee beans';
```

`\p` prints the query buffer with the values in it. `\raw` prints it as you typed it.

usql does not change a variable name inside a string. `select ':item';` sends the text `:item`.

`:{?name}` is `TRUE` when the variable is set, and `FALSE` when it is not:

```sql
\echo :{?item}
```

Meta commands take variables too, such as `\echo :item` or `\c :url`.

## Set a variable from a command

A backtick runs a command in the shell, and gives its output to a meta command. usql uses the shell that `USQL_SHELL` names, or `SHELL` when `USQL_SHELL` is not set:

```sql
\set day `date +%A`
\echo Today is :day
```

A backtick works only in a meta command. In SQL, a backtick is a character of the statement.

`\gset` sets variables from the result of a query. [Run queries and scripts](/guides/run-queries/) explains it.

## Run lines only when a condition is true

`\if` starts a block of lines. usql runs the lines of the first branch whose condition is true, and skips the rest. `\elif` adds a branch with another condition. `\else` adds a branch that runs when no condition is true. `\endif` ends the block:

```sql
select count(*) > 2 as many, count(*) = 0 as none from products \gset
\if :none
  \echo No products
\elif :many
  \echo More than two products
\else
  \echo One or two products
\endif
```

A condition is a boolean value, such as `true`, `false`, `on`, `off`, `yes`, `no`, `1` or `0`. Any other value is an error. To compare values, let the database do it, and keep the result with `\gset`, as above.

A line in a branch that does not run sends no statement, runs no meta command and changes no variable. At the prompt, usql shows `@>` in such a branch.

## More about variables

The [README of usql](https://github.com/xo/usql#variables) describes the other kinds of variables: the named connections of `\cset`, the print options of `\pset` and the settings of `\uset`.
