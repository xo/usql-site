# Recording

A recording is an asciicast v2 file in `public/casts/`. `cmd/record` makes it
from the script of the same name in `recordings/`. Never edit a recording by
hand, and never record one with another program (D5). A file of
`public/charts/` is also made by `cmd/record`, from the `keep` of a script.

## Make a recording

You need Go, and the release of usql that the site documents on your `PATH`.
`recordings/usql-version` holds the version of that release, as
`usql --version` prints it after `usql `. The recorder runs
`usql --version` before it records. If the version is not the one in the
file, the recorder stops and names both versions.

Run the script, by its name, in the repository root:

```bash
go run ./cmd/record tour
```

`go run ./cmd/record -all` makes every recording again. For a new release of
usql, do these steps:

1. Put the new release first on your `PATH`.
2. Write its version to `recordings/usql-version`.
3. Run `go run ./cmd/record -all`.
4. Review the changes to the recordings and to `public/charts/` before you
   stage them.

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
| `keep`  | files of the home directory to copy to `public/charts/<name>/` after the recording, such as the SVG file of a chart   |

Each step has one key:

| Step           | What it does                                                                                                                  |
| -------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| `type: TEXT`   | types the text, then Enter                                                                                                    |
| `keys: TEXT`   | types the text with no Enter. `\t` is Tab and `\x1b` is Escape                                                                |
| `wait: TEXT`   | waits until the terminal shows the text, for up to 20 seconds                                                                 |
| `sleep: 1.5s`  | waits, so that the reader can read                                                                                            |
| `secret: NAME` | types the value of the variable `NAME` of `needs`, then Enter, at a prompt that does not show it, such as the password prompt |

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
6. It removes `public/charts/<name>/`, and copies each file of `keep` into
   it. It stops when a file of `keep` is missing, or holds the value of a
   variable of `needs`.

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

Load the sample data with `dbrun fixture northwind <selector>`. The
recordings use these variables, each with the URL that `dbrun dsn` prints,
pointed at the database of the northwind fixture:

| Variable                  | What it holds                                     |
| ------------------------- | ------------------------------------------------- |
| `USQL_SITE_PG`            | PostgreSQL 18                                     |
| `USQL_SITE_MY`            | MySQL                                             |
| `USQL_SITE_PG_NOPASSWORD` | PostgreSQL 18, with no password in the URL        |
| `USQL_SITE_PG_PASSWORD`   | the password of that user, for the `secret:` step |

Read the owner that `dbrun` reports. If another owner started the container,
do not load, change or drop anything in it. Stop the containers that you
started when you are done. Never write the address or the password of a
server into a script or a page.

## A chart

The player draws text and no images. A chart page shows the recording of the
typing, and the chart as an SVG file that the same session writes with
`\chart ... file=NAME.svg`. Name the file in `keep`:

```yaml
steps:
  - type: \chart title='Cups sold' file=sales.svg
  - wait: "writing to sales.svg"
keep:
  - sales.svg
```

The page shows the file as a picture:

```md
![A bar chart of the cups of coffee, tea and cocoa sold each month](/charts/chart-bar-line/sales.svg)
```

In usql 0.0.0-85d27b9, `\chart` stops with "graphics not supported" when
the terminal draws no images, even with `file=`. The recorder sets
`TERM_GRAPHICS=none`, so a chart script writes an `init.sql` that turns the
graphics on before the first prompt:

```yaml
files:
  .config/usql/init.sql: |
    \uset quiet on
    \uset term_graphics kitty
```

`quiet` must come first. Without it, `\uset` prints a message, and usql
draws its logo with the graphics of kitty. With `file=`, `\chart` writes the
file and draws nothing, so no image reaches the recording. Do not use this
`init.sql` in a recording of `-f`, because it turns the graphics on for `-f`
too. Remove it when usql writes a chart file with no graphics (W5).

## Show a recording

Add the directive to the page, on a line of its own:

```md
::cast[What the recording shows]{name="tour"}
```

The build stops when the recording does not exist.
