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

## W4. Record against database servers

A script names the variables that it needs, such as `USQL_SITE_PG`, and
`files:` writes them into the `config.yaml` of the temporary home as named
connections. The recording then types `usql pg`, and shows no address.
Add a script for each guide that needs a server. Start each server with
dbmeta's `dbrun`, as `docs/RECORDING.md` says, and add the sample data that
the guides use, such as the `sales` table of usql's `_samples/chart.sql`.

## W5. Show the chart of a recording

The player draws text and no images. A chart page shows the recording of the
typing, and the chart as the SVG file that the same session writes with
`\chart file=`. Two changes:

1. usql: `\chart file=` must work when the terminal draws no images. Today
   `doExecChart` in usql's `handler/handler.go` checks for an encoder before
   it reads the options, and stops with "graphics not supported".
2. `cmd/record`: a `keep:` list of files that it copies from the temporary
   home to `public/charts/<name>/` after the recording.

## W6. Charts from a Vega-Lite file

usql reads the directory `charts.path` at startup, and no command uses it.
The guide "Charts from a Vega-Lite file" waits for usql to draw a chart from
a file in that directory, with the result of a query as its data.

## W7. Pin the usql of the recordings

Write the release of usql that the recordings use to
`recordings/usql-version`, or take it from `go.mod` once W2 pins it, and make `cmd/record` stop when `usql --version` on the `PATH`
prints another version. Record every script again for each release with
`go run ./cmd/record -all`.

## W8. Write the guides of the plan

The pages of [PLAN.md](PLAN.md) that are not written yet. Take each from the
README of usql, run every command in the release, and add a recording where
the plan says so.

On 2026-11-27, three agents wrote 14 guides and 10 recordings, with a build
of usql main at 85d27b9 that held the SQLite and DuckDB drivers. What is
left:

- Passwords, and Databases in Kubernetes. They need a server or a cluster
  to run the commands against.
- The chart pages after "Your first chart", which wait for W5 and W6.
- The recipes.
- The command line options that the guides leave out, such as `-1`, `-o`,
  `-U` and `--dump`, until the move of usql to ox (usql's W83) is released.
- The parts that the bugs of W14 hide: `\copy` with named connections, the
  `DRIVER PARAMS` form of `\c`, `\dn` and `\d` of an index on DuckDB.

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
`cms/editor.js` holds a list of tabs, each with a label and its Markdown. The
install page uses it. The paths of the configuration directory belong to
the guide "Configure usql" of W8, which can use the tabs.

The tabs do not remember the choice of the reader from one page to the
next, as Starlight's `syncKey` does. Add it when two pages have tabs with
the same labels.

## W11. Move the guide of the README to the site

Once the guides are on the site, the README of usql can shrink to an
introduction, the install commands and a link to usql.app. A maintainer
decides it in usql. Open question 3 of [PLAN.md](PLAN.md).

## W12. Release notes

A page for each release of usql, from its GitHub release, made by `cmd/gen`
(W2).

## W13. Check the links: Done

`starlight-links-validator` runs in the build, so that a link to a page that
does not exist, or to a heading that a page does not have, stops it. A link
to a file of `public/` passes. `astro.config.mjs` excludes `/admin/`, which
`scripts/prepare.mjs` makes.

## W14. Fix the usql bugs found while writing the guides

The agents that wrote the guides of W8 found these with usql main at
85d27b9. Each one is a change to usql, which a maintainer makes there. The
first three were run again by hand, and failed the same way.

1. `\copy` does not look up a named connection, though the README says it
   does. After `\cset shop sq:shop.db`,
   `\copy shop dk:warehouse.duckdb 'select id, name from products' products`
   prints `error: invalid database scheme`. `metacmd.Copy` passes the name
   to `dburl.Parse`.
2. `\c DRIVER PARAMS` fails, and so do `\cset NAME DRIVER PARAMS` and the
   list form of `connections:` in `config.yaml`. `\c sqlite3 books.db`
   prints `error: : driver not available`. `handler.Open` builds the URL
   with no `SchemeName`.
3. `usql nothere.db` makes an empty SQLite file. The README says that a path
   with no driver must exist.
4. `\timing on` from `init.sql`, from `\i`, or from an earlier `-c` prints
   "Timing is on.", and the next query prints no time.
5. When the last statement of `init.sql` fails, usql exits, though
   `on_error_stop` is off. `\i` of a file whose last statement fails prints
   the error twice and skips the next `-c`.
6. `\set ON_ERROR_STOP on` is accepted and does nothing.
7. `\begin -read-only` on SQLite lets an `INSERT` run.
8. DuckDB: `\d products` shows no primary key, `\d` of an index lists no
   columns, `\dn` leaves out `main`, and `\dnS` lists `main` three times.
   These come from dbmeta.
9. Time: after `\pset timezone America/New_York`, a DuckDB `date` of
   2026-11-02 prints as Nov 1. `\pset timezone` prints an empty line. A
   SQLite `timestamp` prints in the time zone of the machine, and DuckDB
   prints it as stored.
10. On SQLite, a statement prints the row count of an earlier one, such as
    `CREATE INDEX 3`, `BEGIN 1` and `COMMIT 1`.
11. Smaller ones: DuckDB `4.50::decimal(8,2)` prints `4.5`. With `\t`,
    `\pset border 2` and `\pset linestyle unicode`, the top border is
    missing. The error of `\pset format bogus` leaves `vertical` out of its
    list. A script that fails only in `\gset` exits 0. The prompt escape `%x`
    does nothing. `\uset term_graphics bogus` is accepted.
