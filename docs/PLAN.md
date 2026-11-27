# Plan

usql-site is the documentation site of usql, at `https://usql.app`. This file
holds the plan of the site and the questions that are open for the
maintainers. The work items are in [BACKLOG.md](BACKLOG.md), each with a `W`
number. The decisions are in [decisions/](decisions/README.md), each with a
`D` number.

## Direction

- The site is static. Its pages are Markdown files in this repository, and
  Sveltia CMS edits them in the browser at `/admin/` (D2).
- Astro and Starlight build it. Starlight gives the sidebar, the search
  (Pagefind, which runs in the browser), the dark and light themes, and the
  code blocks with a copy button (D3).
- A page is Markdown and nothing else, so that the editor can open every
  page. A recording is a directive in the Markdown (D4).
- Every recording is made by a program from a script, and is made again for
  each release of usql (D5).
- The pages load nothing from another host (D6).
- The reference pages are generated from usql, so that they cannot drift from
  the program (D7).

## The sites to learn from

Three flagship models were asked on 2026-11-27 which open source
documentation sites to learn from (GPT-6 Astra, Gemini 3.1 Pro and DeepSeek
V4 Pro). All three named DuckDB and ClickHouse, and two put PostgreSQL
first. Both answers that reached the tools named Starlight and Expressive
Code. The answer of DeepSeek stopped after three sites.

| Site                      | Source                                | What to take                                                                                                                                                              |
| ------------------------- | ------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| PostgreSQL, the psql page | `postgres/postgres`, `postgres/pgweb` | The complete reference of each meta command and variable, with its syntax, its default and its caveats. usql follows psql, so its reference has the same shape            |
| DuckDB                    | `duckdb/duckdb-web`                   | An install page with one block for each system, and short SQL examples that the reader can run                                                                            |
| ClickHouse                | `ClickHouse/clickhouse-docs`          | The split of quick starts, guides of integrations, the SQL reference and the reference of settings. A table of settings with the name, the type, the default and the text |
| rclone                    | `rclone/rclone` (`docs/`)             | One page for each backend, with the same headings on each. usql needs the same for its databases                                                                          |
| Nushell                   | `nushell/nushell.github.io`           | A book that goes from the first steps to recipes, with input and output on every page                                                                                     |
| Miller                    | `johnkerl/miller` (`docs/`)           | Every example shows the command and its exact output, which the build makes by running the command                                                                        |
| Observable Plot           | `observablehq/plot` (`docs/`)         | A gallery of charts, each with its result and the least code that makes it                                                                                                |
| Supabase                  | `supabase/supabase` (`apps/docs`)     | The reference of the command line, one page for each command                                                                                                              |
| Charm                     | `charmbracelet/vhs`                   | Terminal recordings that a script makes, so that they are made again for each release                                                                                     |
| Astro                     | `withastro/docs`                      | How a Starlight site is laid out, and how its pages are written                                                                                                           |

The tools that the models named, and what the site does with them:

- Starlight (`withastro/starlight`) is the theme of the site.
- Expressive Code is part of Starlight. It draws the code blocks, with the
  frame of a terminal and a copy button.
- asciinema-player (`asciinema/asciinema-player`) plays the recordings. It is
  served from this site.
- Pagefind is part of Starlight. It is the search.
- starlight-versions (`HiDeoo/starlight-versions`) can keep the pages of older
  releases. The site documents only the latest release for now.

## The pages

The site has six groups. A page marked "generated" comes from usql or dburl,
and a person does not write it. A page marked "recording" has at least one.
A page marked "Written" is on the site. The rest are W8.

### Get started

| Page             | Notes                                                                   |
| ---------------- | ----------------------------------------------------------------------- |
| Install usql     | Written. One tab for each way to install it (W10)                       |
| First steps      | Written, with a recording                                               |
| Coming from psql | Written. What works the same, what usql adds, and what it does not have |
| Get help         | Written. `\?`, `\help`, `--help`, and GitHub Discussions                |

### Guides

| Page                        | Notes                                                                                        |
| --------------------------- | -------------------------------------------------------------------------------------------- |
| Connect to a database       | Written. The URLs, the short names, the paths on disk and `USQL_DSN`                         |
| Named connections           | Written. `\cset`, the `connections:` of `config.yaml`, recording                             |
| Passwords                   | `-W`, `-w`, `USQL_DSN`, and why there is no password file                                    |
| Databases in Kubernetes     | The `k8s` URLs of dbfwd                                                                      |
| Run queries and scripts     | Written. The query buffer, `\g`, `\gx`, `\gset`, `\watch`, `\e`, `-c`, `-f`, `\i`, recording |
| Variables and conditionals  | Written. `\set`, interpolation, backticks, `\if`, recording                                  |
| Format the output           | Written. `\pset`, the formats, `\x`, the pager, the time formats, recording                  |
| Describe a database         | Written. The `\d` commands, which read dbmeta, recording                                     |
| Copy data between databases | Written. `\copy` between two URLs, and the import of a CSV file, recording                   |
| Dump a database             | Written. `\dump`, `--dump`, and a dump for another database, recording                       |
| Transactions                | Written. `\begin`, `\commit`, `-1`                                                           |
| Configure usql              | Written. `config.yaml`, `init.sql`, `\uset`, `-U`, the prompt, the history, recording        |
| Line editing and completion | Written. The keys, the `inputrc` file, the completion of names, recording                    |
| Syntax highlighting         | Written. The styles, recording                                                               |
| Terminal graphics           | Written. The detection, the variables, and the terminals that draw images                    |

### Charts

| Page                         | Notes                                                                |
| ---------------------------- | -------------------------------------------------------------------- |
| Your first chart             | Written. The shape of the result and the options                     |
| Bar and line charts          | The series, the labels and the precision, with the chart images (W5) |
| Size, colors and titles      | The options, each with an image                                      |
| Save a chart to a file       | `file=`. Needs the change to usql of W5                              |
| Charts in a script           | `\i`, and why `-f` draws no chart                                    |
| Charts from a Vega-Lite file | Needs the feature of usql of W6                                      |
| Gallery                      | Every chart of the pages above, each with the query that makes it    |

The player of the recordings draws text and no images, so a chart page shows
two things: the recording of the typing, and the chart as an SVG file that
the same session writes with `\chart file=`. W5 holds the steps.

### Databases

The list of every database that usql connects to, and a page for each main
database, with the schemes and aliases from the scheme registry of dburl and
an example URL (W9). The main databases have pages: PostgreSQL, MySQL,
SQLite, SQL Server, Oracle, DuckDB and ClickHouse. A person adds notes to the
page of a database, such as the options of its driver and what it does not
support, and adds a page for another database.

### Reference

All generated (W2):

- The command line options.
- The meta commands, with their arguments.
- The settings of `\uset`, the print options of `\pset`, the variables, the
  connection variables, and the settings of `\chart`.
- The environment variables.
- The keys of `config.yaml`.
- The URL schemes and their aliases.

### Recipes and project

Recipes are short pages for one task each, such as "Export a query to CSV",
"Import a CSV file", "Compare two databases" and "Run usql in CI". The
project group has "Edit this site" (written), the release notes (W12), and
how to contribute to usql.

## Recordings

`cmd/record` makes every recording from a script in `recordings/`. It runs a
shell in a pseudo-terminal, in a temporary home directory, types each step,
and writes an asciicast v2 file. A page shows it with
`::cast[caption]{name="..."}`. [RECORDING.md](RECORDING.md) has the details.

- A recording that needs only SQLite or DuckDB makes its database in its
  `setup:`. The first one, `tour`, does that.
- A recording that needs a server reads its address from a variable that it
  names in `needs:`. The person who records starts the server with dbmeta's
  `dbrun` (W4).
- The recordings are made with the release of usql that the site documents,
  and are made again for each release (W7).

## Open questions for the maintainers

1. Hosting. The plan is GitHub Pages, built by
   `.github/workflows/site.yml`, with the custom domain `usql.app`. A
   maintainer must turn on Pages for `xo/usql-site` with "GitHub Actions" as
   the source, and point the DNS of `usql.app` at it. The old site was on
   Netlify. Is GitHub Pages right, or should the site stay on Netlify or move
   to Cloudflare Pages? (W1)
2. Sign-in to the editor. "Sign In with GitHub" needs a GitHub OAuth
   application and a small service that holds its secret, such as
   `sveltia/sveltia-cms-auth` on Cloudflare Workers. Until then, an editor
   signs in with a GitHub token that can write to the repository. Is a token
   enough, or should the service be set up? (W3)
3. The README of usql. Once the site holds the guides, the README of usql can
   shrink to an introduction, the install commands and a link to the site.
   Should it? (W11)
4. The reference needs lists that usql does not export: the options of the
   command line are in package main, and the names of the settings and
   variables are not exported from `env`. Should usql export them, or should
   the generator read `usql --help` and `\set` output? (W2)
5. `\chart file=` stops with "graphics not supported" when the terminal draws
   no images, though it only writes a file. A recording cannot write a chart
   until usql allows it. (W5)
6. `charts.path` names a directory of chart files, and usql reads it at
   startup, but no command uses it yet. Complex charts from a Vega-Lite file
   need that feature in usql. (W6)
