# usql-site backlog

This file records the planned work for the site. Each item names the files it
touches, so that the work can start without finding the context again.

## How to refer to an item

Every item has an identifier of the form `W` followed by a number. Three rules
govern it:

1. Identifiers are append only. A new item takes the next unused number.
2. An identifier is never reused.
3. An identifier is never renumbered.

An item that is finished or abandoned keeps its heading and gains a status in
that heading, as in `## W1. Publish the site at usql.app: Done`. The statuses
in use are `Done`, `Dropped` and `Superseded by Wn`. A heading with no status
is open.

These numbers are separate from the `W` numbers of usql. A reference to an
item of usql names the repository, as in "usql's W83".

## W1. Publish the site at usql.app

`.github/workflows/site.yml` runs the local check on every push and pull
request, and publishes `dist/` to GitHub Pages from `main`. `public/CNAME`
holds `usql.app`. A maintainer must:

1. Turn on GitHub Pages for `xo/usql-site`, with "GitHub Actions" as the
   source.
2. Set the variable `PAGES_ENABLED` of the repository to `true`.
3. Point the DNS of `usql.app` at GitHub Pages, and turn on HTTPS.
4. Remove the old site from Netlify.

Open question 1 of [PLAN.md](PLAN.md) asks whether GitHub Pages is the right
host.

## W2. Generate the reference pages from usql

The reference pages list what usql has: the options of the command line, the
meta commands, the settings, the print options, the variables, the
environment variables and the keys of `config.yaml`. `cmd/gen`, which
writes the schemes of W9, will write them as JSON to `src/data/generated/`,
and Astro pages under `src/pages/reference/` will draw them with
`StarlightPage`. The pages are not in the editor, because nobody edits them.

`go.mod` pins the release of usql that the site documents, as usql pins
dburl. usql must export the lists first. The options are in package main
(`cmdline.go`), and the names of the settings are unexported in
`env/list.go`. Open question 4 of [PLAN.md](PLAN.md).

## W3. Sign in to the editor with GitHub

`cms/config.yml` has `base_url: __AUTH_URL__`, and `scripts/prepare.mjs`
fills it from `CMS_AUTH_URL`. It needs a GitHub OAuth application of the `xo`
organization, and a service that exchanges its code for a token, such as
`sveltia/sveltia-cms-auth`. Set `CMS_AUTH_URL` in the workflow when the
service exists. Until then, an editor signs in with a token. Open question 2
of [PLAN.md](PLAN.md).

## W4. Record against database servers: Done

A script names the variables that it needs, such as `USQL_SITE_PG`, and
`files:` writes them into the `config.yaml` of the temporary home as named
connections. The recording then types `usql pg`, and shows no address.
`docs/RECORDING.md` lists the variables. The recordings `connect`,
`postgres`, `mysql` and `passwords` use PostgreSQL 18 and MySQL with the
northwind fixture of dbrun. The step `secret:` types the value of a variable
at the password prompt, which does not show it.

The example URLs of the pages of PostgreSQL, MySQL and SQL Server were run
against their servers on 2026-11-27. The named instance form of SQL Server,
and the pages of Oracle and ClickHouse, were not run.

## W5. Show the chart of a recording

The player draws text and no images. A chart page shows the recording of the
typing, and the chart as the SVG file that the same session writes with
`\chart file=`. Two changes:

1. usql: `\chart file=` must work when the terminal draws no images. Today
   `doExecChart` in usql's `handler/handler.go` checks for an encoder before
   it reads the options, and stops with "graphics not supported". Until it
   is fixed, a chart script writes an `init.sql` with `\uset quiet on` and
   `\uset term_graphics kitty`, as `docs/RECORDING.md` says. When usql
   fixes it, remove that `init.sql` from the chart scripts, and the notes
   about the current release from the chart pages.
2. `cmd/record`: `keep:` copies files from the temporary home to
   `public/charts/<name>/` after the recording. Done.

## W6. Charts from a Vega-Lite file

usql reads the directory `charts.path` at startup, and no command uses it.
The guide "Charts from a Vega-Lite file" waits for usql to draw a chart from
a file in that directory, with the result of a query as its data.

## W7. Pin the usql of the recordings: Done

`recordings/usql-version` holds the version of usql that the recordings use,
and `cmd/record` stops when `usql --version` on the `PATH` prints another
one. Record every script again for each release with
`go run ./cmd/record -all`. Once W2 pins usql in `go.mod`, take the version
from there.

## W8. Write the guides of the plan

The pages of [PLAN.md](PLAN.md) that are not written yet. Take each from the
README of usql, run every command in the release, and add a recording where
the plan says so.

On 2026-11-27, three agents wrote 14 guides and 10 recordings, with a build
of usql main at 85d27b9 that held the SQLite and DuckDB drivers. What is
left:

- Databases in Kubernetes. It needs a cluster to run the commands against.
- Charts from a Vega-Lite file, which waits for W6.
- The recipes.
- The command line options that the guides leave out, such as `-1`, `-o`,
  `-U` and `--dump`, until the move of usql to ox (usql's W83) is released.
- The parts that the bugs of W14 hid, once usql's fixes are released:
  `\copy` with named connections, the `DRIVER PARAMS` form of `\c`, and
  `\dn` and `\d` of an index on DuckDB.

A second round on the same day wrote Passwords, with recordings against
PostgreSQL and MySQL (W4), and five chart pages with eight SVG charts (W5).

The recorder types each key on its own, so a script cannot type a key that
sends more than one byte, such as Up. The line-editing recording uses
Ctrl-P.

## W9. A page for each database: Done

`cmd/gen` writes `src/data/generated/schemes.json` from the registry of the
release of dburl that `go.mod` pins, as usql's `gen.go` makes the README
tables. `TestSchemesAreCurrent` fails when the file is out of date.
`src/lib/schemes.mjs` draws two directives from it: `::schemes{name="..."}`,
the schemes, aliases and drivers of one database, and `::databases`, the
table of every database. `src/content/docs/databases/` has the list, and a
page with an example URL for PostgreSQL, MySQL, SQLite, SQL Server, Oracle,
DuckDB and ClickHouse. A page for another database is a new Markdown file
named by its dialect, as `docs/CONTENT.md` says.

The example URLs of the pages of a server were not run against a server.
usql read each of them and chose the right driver, and the connection was
refused. Run each against its server when W4 has them.

## W10. Tabs for each system: Done

`:::tabs`, with one `::::tab[label]` section for each tab, is drawn by
`src/lib/tabs.mjs` with the markup and the look of Starlight's Tabs
component. `src/scripts/tabs.ts` switches the tabs. The editor component in
`cms/editor.js` holds a list of tabs, each with a label and its Markdown.

`:::tabs{sync="key"}` keeps the choice of the reader for every tabs with the
same key, on the page and on the next pages, in the storage of the browser.
A short script in the tabs opens the saved tab before the page shows. The
install page uses `sync="install"`.

## W11. Move the guide of the README to the site

Once the guides are on the site, the README of usql can shrink to an
introduction, the install commands and a link to usql.app. A maintainer
decides it in usql. Open question 3 of [PLAN.md](PLAN.md).

## W12. Release notes: Done

`go run ./cmd/gen -releases` reads the releases of usql from GitHub and
writes `src/data/generated/releases.json`. It is the only mode of `cmd/gen`
that uses the network. `::releases`, on the page `project/releases`, draws
them. The notes are Markdown, and the plugin draws them safely: raw HTML
shows as text, an image becomes a link, and a link keeps its target only
when it is http, https or mailto. `src/lib/releases.test.mjs` checks these
rules, and `scripts/precommit.sh` runs it.

Run `go run ./cmd/gen -releases` after each release of usql. The workflow
does not run it yet.

## W13. Check the links: Done

`starlight-links-validator` runs in the build, so that a link to a page that
does not exist, or to a heading that a page does not have, stops it. A link
to a file of `public/` passes. `astro.config.mjs` excludes `/admin/`, which
`scripts/prepare.mjs` makes.

## W14. Fix the usql bugs found while writing the guides

The agents that wrote the guides of W8 found these with usql main at
85d27b9. On 2026-11-27 they were fixed in usql, each with a test, and the
fixes wait in usql's checkout for a maintainer to commit them. usql's
`docs/RELEASE.md` lists them under "Fixes found while writing usql.app":

- `\copy` with named connections, and `\c DRIVER PARAMS`.
- A bare path that does not exist, which created an empty SQLite file.
- `-W` with a named connection.
- `\set ON_ERROR_STOP`, the `\pset format` error, `\gset` exit codes, and
  the prompt escape `%x`.
- `\timing` lost after `\i` and `init.sql`, the error of the last statement
  of `init.sql`, and `\i` that printed its error twice.
- `\begin -read-only` on SQLite, which is now an error.
- `\pset timezone`, which moved a date.
- Stale row counts on SQLite, and the scale of a DuckDB decimal.
- `\chart bg=`, and a NULL in a series of `\chart`.
- `\dnS` on DuckDB. dbmeta v0.50.1 fixed the other gaps of `\d` and `\dn`
  on DuckDB.

What is left:

1. dburl: `file:` with an absolute path that does not exist opens the
   directory above it as PostgreSQL, so it does not create a SQLite file as
   usql's README says. The dburl session confirmed it on 2026-11-27, and a
   maintainer decides the fix, because it changes how a bare path resolves
   too.
2. usql: an error of a meta command starts with the name of the driver, as
   in `error: sqlite3: \chart: graphics not supported in terminal`, and the
   error of `\chart file=` comes after `writing to` was printed.
3. usql: a meta command after a failed statement resets the exit status of a
   script.
4. usql: a `\o` in a file of `\i` is lost when the file ends.
5. usql: `\chart file=` needs a terminal that draws images (W5).

The missing top border with `\t` and `\pset border 2` is not a bug: psql
18.7 draws the same, and tblfmt matches it. A SQLite `timestamp` that prints
in the time zone of the machine is not a bug either: dburl's D67 sets
`_loc=auto` on purpose, and `?_loc=UTC` in the URL keeps UTC.
`\uset term_graphics bogus` is accepted on purpose, as none.
