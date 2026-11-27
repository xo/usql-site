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
environment variables and the keys of `config.yaml`. A Go program,
`cmd/gen`, will write them as JSON to `src/data/generated/`, and Astro pages
under `src/pages/reference/` will draw them with `StarlightPage`. The pages
are not in the editor, because nobody edits them.

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

## W9. A page for each database

Generate a page for each scheme of dburl's registry, with its aliases, its
driver and an example URL, as usql's `gen.go` makes the README tables. A
person adds notes to the main databases first.

## W10. Tabs for each system

A container directive, `:::tabs`, with one `::tab[macOS]` and so on in it,
for the install page and the paths of the configuration directory. It needs
a plugin in `src/lib/` and an editor component in `cms/editor.js` (D4).

## W11. Move the guide of the README to the site

Once the guides are on the site, the README of usql can shrink to an
introduction, the install commands and a link to usql.app. A maintainer
decides it in usql. Open question 3 of [PLAN.md](PLAN.md).

## W12. Release notes

A page for each release of usql, from its GitHub release, made by `cmd/gen`
(W2).

## W13. Check the links

Add `starlight-links-validator` to the build, so that a link to a page that
does not exist stops it.
