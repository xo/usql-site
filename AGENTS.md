# Rules for agents working on usql-site

This repository is the documentation site of usql, at `https://usql.app`.
It is a static site built with Astro and Starlight. Its pages are Markdown
files in this repository, and people edit them on the site with Sveltia CMS.
This file holds the rules that an agent can break without noticing. The
reasons are in the documents listed below.

## Standing rules

These hold in every `xo` repository, for every coding agent (D1).

1. Stage changes for review. Commit and push only when a maintainer says so.
2. Load the `simple-english` skill before you write any text that a person
   reads: a page of the site, a document, a code comment, an error message or
   a commit message. Follow it for that text.
3. Load the `go-pedantry` skill before you write or review Go code. Follow it
   where it does not conflict with a rule in this file. A rule here wins.

`CLAUDE.md` holds one line that imports this file, so that Claude Code and
every other agent read the same rules. Edit this file, not that one.

## Which document to read

| Task                                                              | Document                                             |
| ----------------------------------------------------------------- | ---------------------------------------------------- |
| Writing or changing a page                                        | [docs/CONTENT.md](docs/CONTENT.md)                   |
| Making a terminal recording                                       | [docs/RECORDING.md](docs/RECORDING.md)               |
| Building and publishing the site                                  | [docs/DEPLOY.md](docs/DEPLOY.md)                     |
| The plan of the pages, and the questions open for the maintainers | [docs/PLAN.md](docs/PLAN.md)                         |
| Planned work, and how to cite it                                  | [docs/BACKLOG.md](docs/BACKLOG.md)                   |
| Why something is the way it is                                    | [docs/decisions/README.md](docs/decisions/README.md) |
| Setting up, checking a change, and the agent skills               | [CONTRIBUTING.md](CONTRIBUTING.md)                   |
| What the repository is, for a person                              | [README.md](README.md)                               |

A document that is not in this table does not exist. If you find one, either
add a row for it or delete it.

## Hard rules

1. Do not edit generated files. `public/admin/` is made by
   `scripts/prepare.mjs` from `cms/`. A file of `public/casts/` is made by
   `cmd/record` from the script of the same name in `recordings/`.
   `src/data/generated/` is made by `cmd/gen` from the Go modules that
   `go.mod` pins. Change the source, then make the file again.

2. A page is a `.md` file, never `.mdx`, so that the web editor can open it.
   A page imports no component. Anything that is more than Markdown is a
   directive, such as `::cast[...]{name="..."}`. A directive has a plugin in
   `src/lib/` and an editor component in `cms/editor.js`. Add or change both
   in the same change. D4 says why.

3. A field of a page is in two places: the schema in
   `src/content.config.ts` and the editor fields in `cms/fields/page.yml` or
   `cms/collections/`. Change both in the same change. A new folder of pages
   also needs a group in the sidebar in `astro.config.mjs` and a file in
   `cms/collections/`.

4. The pages load nothing from another host: no font, no script, no style,
   no image and no analytics. `scripts/precommit.sh` checks the built pages.
   The web editor at `/admin/` is the one exception, and only for what Sveltia
   CMS loads that is not copied here. D6 says why.

5. Facts about usql come from usql. Do not copy a list of options, commands,
   settings, schemes or drivers into a page by hand. Such a list goes stale
   the day usql changes. The schemes of the databases come from dburl,
   through `::schemes` and `::databases`. The reference pages will be
   generated from usql (W2). Until then, a page links to the README of usql
   for the full list.

6. Every command that a page shows must work in the usql release that the
   site documents. Run it before you write it. A recording is the best proof.

7. Start a database for a recording only with dbmeta's `dbrun`, with
   `DBMETA_OWNER_NAME=usql`. Never use a container command of your own. Check
   the owner that `dbrun start` reports. If someone else started it, do not
   load, change or drop anything in it.

8. A recording script never holds a URL with a password, a token or a
   host name of a machine. It names the variables it needs in `needs:`, and
   writes them into the temporary home directory with `files:`. The recorder
   fails when a recording shows the value of one of them.

9. Do not record a recording on a machine with another usql than the release
   that the site documents. `docs/RECORDING.md` says how to check it.

10. The repository root holds four documents: `README.md`, `AGENTS.md`,
    `CLAUDE.md` and `CONTRIBUTING.md`. Every other document goes under
    `docs/`. See D1.

11. Backlog identifiers are append only. Never reuse a `W` number and never
    renumber one. See the rules at the top of `docs/BACKLOG.md`.

12. npm versions are exact (`.npmrc`), and the lockfile is committed. Update
    a package with `npx npm-check-updates`, then run the local check.

## Writing

The pages and the documents are written in plain English. Short sentences,
active voice, no semicolons, and no em dashes. State the fact rather than its
importance. Use the words of usql: a "meta command" is a backslash command, a
"setting" is what `\uset` sets, and a "print option" is what `\pset` sets.

Do not name a maintainer in a page, a document or a comment. Write "the
maintainer" for a decision that was made, with its date, and "a maintainer"
for a rule or for something still to decide.

## The local check

```bash
bash scripts/precommit.sh
```

It runs gofmt, `go vet`, `go test`, prettier, `astro check`, the build, and a
check of the built pages, and prints one verdict on its last line.

`go.mod` sits next to `node_modules`, which can hold Go files. Name the Go
packages, `./cmd/...` and `./internal/...`, and never `./...`.

## Sibling repositories

| Repository  | What it holds                                                                                                |
| ----------- | ------------------------------------------------------------------------------------------------------------ |
| `xo/usql`   | The program that the site documents. Its README is the source of most pages until the reference is generated |
| `xo/dburl`  | The connection URL schemes, which `cmd/gen` reads for the pages of the databases                             |
| `xo/dbmeta` | Database metadata, and `dbrun`, which starts the databases of the recordings                                 |

Do not change a sibling repository from here. If usql must change for a page,
for example to export a list that a reference page needs, add an item to
`docs/BACKLOG.md` here and tell a maintainer.
