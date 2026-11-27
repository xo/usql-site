# usql-site

The documentation site of [usql](https://github.com/xo/usql), at
`https://usql.app`. It is a static site built with
[Astro](https://astro.build) and [Starlight](https://starlight.astro.build).
The pages are Markdown files in this repository. People edit them in the
browser at `https://usql.app/admin/`, which is
[Sveltia CMS](https://sveltiacms.app), or with a pull request.

## Run it

You need Node.js 22.23.4 or later. Go is needed only to make recordings.

```bash
npm ci
npm run dev
```

The site is at `http://localhost:4321/`, and the editor is at
`http://localhost:4321/admin/`. Choose "Work with Local Repository" in the
editor to edit the files of your checkout. `npm run build` writes the site to
`dist/`.

## Check a change

```bash
bash scripts/precommit.sh
```

Read the last line for the verdict.

## Layout

| Path                                             | What it holds                                                              |
| ------------------------------------------------ | -------------------------------------------------------------------------- |
| `src/content/docs/`                              | the pages, one Markdown file each                                          |
| `src/lib/`                                       | the Markdown plugins, such as the one for recordings                       |
| `src/scripts/`, `src/components/`, `src/styles/` | the player of the recordings, and the look of the site                     |
| `cms/`                                           | the web editor. `scripts/prepare.mjs` makes `/admin/` from it              |
| `recordings/`                                    | the scripts of the terminal recordings                                     |
| `public/casts/`                                  | the recordings, which `cmd/record` makes from the scripts                  |
| `cmd/record/`                                    | the Go program that makes the recordings                                   |
| `docs/`                                          | the plan, the backlog, the decisions, and how to write, record and publish |

Read `AGENTS.md` and `CONTRIBUTING.md` before you change anything.
