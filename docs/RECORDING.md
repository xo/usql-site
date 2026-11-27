# Recording

A recording is an asciicast v2 file in `public/casts/`. `cmd/record` makes it
from the script of the same name in `recordings/`. Never edit a recording by
hand, and never record one with another program (D5).

## Make a recording

You need Go, and the release of usql that the site documents on your `PATH`.
Make sure that it is that release:

```bash
usql --version
```

Then run the script, by its name, in the repository root:

```bash
go run ./cmd/record tour
```

`go run ./cmd/record -all` makes every recording again. Do that for each
release of usql, and review the changes before you stage them.

## Write a script

```yaml
# What the recording shows, and the page that shows it.
title: A first session with usql
size: 100x26
files:
  books.sql: |
    create table books (id integer primary key, title text not null);
setup:
  - usql -X -q -f books.sql sq:books.db
steps:
  - type: usql books.db
  - wait: "=> "
  - type: \dt
  - wait: "=> "
  - sleep: 1.5s
  - type: \q
  - wait: "$ "
```

| Key     | What it does                                                                                                          |
| ------- | --------------------------------------------------------------------------------------------------------------------- |
| `title` | the title of the recording                                                                                            |
| `size`  | the size of the terminal, as `COLSxROWS`. The default is `100x26`                                                     |
| `needs` | the variables of the environment that the script needs                                                                |
| `files` | files to write into the home directory before the recording. `${NAME}` is the value of the variable `NAME` of `needs` |
| `setup` | commands to run with `sh -c` before the recording. They are not recorded                                              |
| `steps` | what to type, in order                                                                                                |

Each step has one key:

| Step          | What it does                                                   |
| ------------- | -------------------------------------------------------------- |
| `type: TEXT`  | types the text, then Enter                                     |
| `keys: TEXT`  | types the text with no Enter. `\t` is Tab and `\x1b` is Escape |
| `wait: TEXT`  | waits until the terminal shows the text, for up to 20 seconds  |
| `sleep: 1.5s` | waits, so that the reader can read                             |

Wait for the prompt (`"=> "`, or `"$ "` for the shell) after each command,
and end the script at the shell prompt.

## What the recorder does

1. It makes a new temporary directory, and uses it as the home directory and
   the working directory. No configuration, history or `init.sql` of yours
   reaches the recording.
2. It writes `files`, and runs `setup`.
3. It starts `bash` in a pseudo-terminal of `size`, with the prompt `$ `,
   `TERM=xterm-256color`, and `TERM_GRAPHICS=none`, so that usql asks the
   terminal no questions. Only `PATH` comes from your environment.
4. It types each step at a fixed speed, and records what the terminal shows.
5. It stops when the recording shows the value of a variable of `needs`, and
   writes `public/casts/<name>.cast` when it does not.

The player shortens every pause to 2 seconds at most.

## A database server

A script that needs a server names a variable in `needs`, and writes it into
the configuration of usql as a named connection, so that the recording types
a name and shows no address:

```yaml
needs: [USQL_SITE_PG]
files:
  .config/usql/config.yaml: |
    connections:
      pg: ${USQL_SITE_PG}
steps:
  - type: usql pg
```

Start the server only with dbmeta's `dbrun`, with `DBMETA_OWNER_NAME=usql`.
Run it in the `test` directory of a checkout of `xo/dbmeta`:

```bash
DBMETA_OWNER_NAME=usql go run ./cmd/dbrun start postgres-18
```

Read the owner that `dbrun` reports. If another owner started the container,
do not load, change or drop anything in it. Stop the containers that you
started when you are done. Never write the address or the password of a
server into a script or a page.

## Show a recording

Add the directive to the page, on a line of its own:

```md
::cast[What the recording shows]{name="tour"}
```

The build stops when the recording does not exist.
