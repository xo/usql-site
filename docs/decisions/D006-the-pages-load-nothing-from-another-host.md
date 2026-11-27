# D6. The pages load nothing from another host

Status: Proposed.

Proposed on 2026-11-27. A page of the site loads no font, script, style,
image or analytics from another host. `scripts/precommit.sh` checks the
built pages.

## Why

A reader of the site does not tell a third party which pages they read. The
site works when a CDN is down or blocked. The player of the recordings and
the search are served from the site, so they work too.

## The editor

The editor at `/admin/` is not a page that readers see. Sveltia CMS loads
its fonts from a CDN, which `scripts/prepare.mjs` copies to `/admin/fonts/`.
It loads its chunks from its own directory first, and they are copied too.
What it loads beyond that, such as the translations of its interface, comes
from unpkg.com, and the check does not look at `/admin/`.
