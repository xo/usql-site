---
title: Transactions
description: Start a transaction with \begin, end it with \commit or \rollback, and what usql does with an open transaction when it exits.
sidebar:
  order: 16
---

A transaction groups statements, so that the database keeps all of their changes or none of them. usql has three meta commands for a transaction. They work the same on every database that has transactions.

## Start and end a transaction

`\begin` starts a transaction. The statements after it run in the transaction:

```sql
\begin
update products set price = price + 1;
select name, price from products;
```

`\commit` keeps the changes of the transaction. `\rollback` drops them:

```sql
\rollback
```

While a transaction is open, the prompt ends with `=~` in place of `=>`.

## The rules

- usql has one transaction at a time. If a transaction is open, `\begin` prints `error: previous transaction exists`.
- If no transaction is open, `\commit` and `\rollback` print `error: no previous transaction exists`.
- You cannot connect to another database while a transaction is open. End the transaction first.

## A read-only transaction

`\begin -read-only` starts a transaction that only reads. After `-read-only`, you can give an isolation level, such as `serializable`:

```sql
\begin -read-only serializable
```

The driver of the database must support it. For example, DuckDB prints an error, because it does not support read-only transactions. SQLite accepts `-read-only`, but a statement in the transaction can still write.

## BEGIN and COMMIT in SQL

You can also send `begin;` and `commit;` as SQL. Then the database starts the transaction, and usql does not know about it. `\commit` and `\rollback` do not end it, and the prompt does not show it. On SQLite, a `\begin` after it fails, because the database already has a transaction. Use `\begin` so that usql knows that a transaction is open.

## An open transaction at exit

If a transaction is open when usql exits, usql rolls it back. The database keeps none of its changes. This happens when you type `\q` or press Ctrl-D, and at the end of a script that `-f` runs.

To keep the changes, end the transaction with `\commit` before you exit. In a script, put `\commit` at the end.
