# D2. The site is static, and its pages are in this repository

Status: Decided.

The maintainer decided on 2026-11-27 that the documentation of usql is a
static site at `usql.app`, edited with Sveltia CMS, which saves to GitHub. The
pages are in the same repository as the code of the site. There is no second
repository for the content.

## What follows

- A save in the editor is a commit to `main` of `xo/usql-site`, and the site
  is built from `main`. There is no release step between a save and the
  site.
- A change to the shape of a page, such as a new field, changes the schema,
  the editor and the pages in one commit. Two repositories could not do that.
- The editor is served by the site itself, at `/admin/`. `scripts/prepare.mjs`
  makes it at each build from `cms/` and the npm package of Sveltia CMS.
- An editor needs write access to `xo/usql-site`. A person without it opens a
  pull request.

## Why

A separate repository for the content lets a release of the site differ from
the content that people save. The site of usql has no such need. One
repository keeps the schema and the pages in step, and a pull request from
anyone works the same way as a save in the editor.

The old site, a Hugo site that described xo and usql together, was removed on
2026-11-27. Its history stays in git.
