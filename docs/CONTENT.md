# Content

Every page of the site is a Markdown file in `src/content/docs/`. The folder
of a page is its group in the sidebar, and its file name is its address:
`src/content/docs/guides/connect.md` is `https://usql.app/guides/connect/`.

## The groups

| Folder                      | Group of the sidebar | Editor collection                |
| --------------------------- | -------------------- | -------------------------------- |
| `src/content/docs/index.md` | the home page        | `cms/collections/00-home.yml`    |
| `start/`                    | Get started          | `cms/collections/10-start.yml`   |
| `guides/`                   | Guides               | `cms/collections/20-guides.yml`  |
| `charts/`                   | Charts               | `cms/collections/30-charts.yml`  |
| `project/`                  | Project              | `cms/collections/90-project.yml` |

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
- Link to another page by its address, such as `/guides/connect/`.
- Do not copy a list of usql's options, commands or settings into a page.
  See rule 5 of `AGENTS.md`.

## Directives

A page holds Markdown and these directives, and nothing else (D4):

| Directive                                      | What it shows                                                                                        |
| ---------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| `::cast[caption]{name="tour"}`                 | the recording `public/casts/tour.cast`. An optional `poster="0:12"` shows that frame before it plays |
| `:::note`, `:::tip`, `:::caution`, `:::danger` | Starlight's asides, ended by `:::`                                                                   |

[RECORDING.md](RECORDING.md) says how to make a recording.

## Pictures

The editor saves a picture to `public/images/`, and a page shows it as
`![text](/images/name.png)`. Give every picture a text that says what it
shows.
