# D3. Astro and Starlight build the site

Status: Proposed.

Proposed on 2026-11-27. Astro builds the site, and Starlight, the
documentation theme of Astro, gives it its layout.

## Why

Starlight has what a documentation site needs, with no code of our own: the
sidebar, the search, the themes, code blocks with a copy button and the frame
of a terminal, a table of contents, and the link to edit a page. Its search,
Pagefind, runs in the browser from files of the site, so it needs no service.

Starlight reads plain Markdown files with a frontmatter, which is the format
that Sveltia CMS writes. Its schema of the frontmatter checks every page at
build time, so a bad save stops the build and does not reach the site.

Three flagship models, asked which open source documentation sites to learn
from, named Starlight as the theme to use. See [PLAN.md](../PLAN.md).

## Costs

Starlight changes quickly before its 1.0 release, and a new version can
change its configuration. 0.39 changed the form of the sidebar, for example.
Its version is pinned, and an update runs the local check.
