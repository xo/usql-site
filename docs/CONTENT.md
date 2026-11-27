# Content

Every page of the site is a Markdown file in `src/content/docs/`. The folder
of a page is its group in the sidebar, and its file name is its address:
`src/content/docs/guides/connect.md` is `https://usql.app/guides/connect/`.

## The groups

| Folder                      | Group of the sidebar | Editor collection                  |
| --------------------------- | -------------------- | ---------------------------------- |
| `src/content/docs/index.md` | the home page        | `cms/collections/00-home.yml`      |
| `start/`                    | Get started          | `cms/collections/10-start.yml`     |
| `guides/`                   | Guides               | `cms/collections/20-guides.yml`    |
| `charts/`                   | Charts               | `cms/collections/30-charts.yml`    |
| `databases/`                | Databases            | `cms/collections/40-databases.yml` |
| `project/`                  | Project              | `cms/collections/90-project.yml`   |

To add a group, add three things in one change: the folder, a group in the
`sidebar` of `astro.config.mjs`, and a file in `cms/collections/`. A
collection with no `fields` gets the fields of `cms/fields/page.yml`.

## The frontmatter

```md
---
title: Connect to a database
description: The URLs that usql opens, the short names of the drivers, and named connections.
sidebar:
  order: 1
---
```

`title` is required. `description` is the text that search engines and link
previews show. `sidebar.order` sets the place of the page in its group, from

1. Starlight's schema in `src/content.config.ts` checks the frontmatter, and
   the build stops on a page that does not match it. A new field goes in
   `src/content.config.ts` and in `cms/fields/page.yml` in the same change.

## Write a page

- Follow the `simple-english` skill. One task or one topic for each page.
- Start with what the reader gets, then the steps.
- Show a command in a code block with the language, such as `sh` or `sql`.
  Leave out the `$` prompt, so that the copy button copies a command that
  runs.
- Run every command in the release of usql that the site documents before
  you write it.
- Link to another page by its address, such as `/guides/connect/`. The
  build stops on a link to a page that does not exist, and on a link to a
  heading that the page does not have. A link to a file of `public/`, such
  as `/casts/tour.cast`, passes.
- Do not copy a list of usql's options, commands or settings into a page.
  See rule 5 of `AGENTS.md`.

## Directives

A page holds Markdown and these directives, and nothing else (D4):

| Directive                                      | What it shows                                                                                        |
| ---------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `::cast[caption]{name="tour"}`                 | the recording `public/casts/tour.cast`. An optional `poster="0:12"` shows that frame before it plays |
| `:::note`, `:::tip`, `:::caution`, `:::danger` | Starlight's asides, ended by `:::`                                                                   |
| `::schemes{name="postgres"}`                   | the table of the URL schemes and aliases of one database, from dburl                                 |
| `::databases`                                  | the table of every database, with a link to each page of `databases/`                                |
| `:::tabs`                                      | tabs, one for each choice of the reader, such as a system. See below                                 |
| `::releases`                                   | the notes of every release of usql. See below                                                        |

[RECORDING.md](RECORDING.md) says how to make a recording.

The build stops on a directive that is wrong, such as a recording that does
not exist or a database that dburl does not know. The editor has a form for
each directive, under its Insert menu.

### Tabs

Use tabs when the reader picks one of several choices and reads only that
one, such as the way to install usql. A `:::tabs` block holds one
`::::tab[label]` section for each tab, and nothing else. A tab has four
colons and ends with `::::`, so that an aside can be in it:

```md
:::tabs
::::tab[macOS]
Markdown for macOS.
::::
::::tab[Linux]
Markdown for Linux.

:::note
An aside in a tab.
:::
::::
:::
```

The first tab is open when the page loads. Do not put a heading in a tab,
because the table of contents of the page lists every heading.

A sync key keeps the choice of the reader. Write it after `:::tabs`, as in
`:::tabs{sync="os"}`, with a-z, 0-9 and `-`. When the reader picks a tab,
every `:::tabs` with the same key opens the tab with the same label, on the
same page and on the next pages that the reader opens. The browser keeps the
label. Give the same key only to tabs with the same labels, such as `macOS`,
`Linux` and `Windows`. The install page uses the key `install`.

### Release notes

`::releases` shows the notes of every release of usql, newest first, under a
heading for each year. `::releases{limit="20"}` shows only the 20 newest.
The notes come from `src/data/generated/releases.json`, which
`go run ./cmd/gen -releases` writes from the releases of usql on GitHub. Do
not copy the notes into a page.

The maintainers write the notes on GitHub, in Markdown. The build shows raw
HTML in them as text, and turns an image into a link to the image, so that
the page loads nothing from another host (D6). A link keeps its target only
when it is an absolute `http`, `https` or `mailto` URL. A heading of the
notes is a level 4 heading or lower, so that the table of contents lists only
the years and the releases.

## The page of a database

The file name of the page of a database in `databases/` is the dialect of
the database in dburl, such as `postgres.md` or `sqlite3.md`. `::databases`
links to the page by that name, and `::schemes{name="..."}` takes the same
name. The dialects are in `src/data/generated/schemes.json`. Each page has the
same headings:

1. `## Connect`, with a short example URL, and the forms of URL that this
   database adds.
2. `## Schemes and aliases`, with `::schemes{name="..."}`.

Do not copy a scheme, an alias or a driver into a page by hand. `cmd/gen`
writes `schemes.json` from the release of dburl that `go.mod` pins. When
`go.mod` moves to another release, run `go run ./cmd/gen`. A test fails until
the file matches.

## Pictures

The editor saves a picture to `public/images/`, and a page shows it as
`![text](/images/name.png)`. Give every picture a text that says what it
shows.
