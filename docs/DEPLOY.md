# Deploy

The site is the folder `dist/` that `npm run build` makes. It is plain files,
and any static host serves it.

## The build

```bash
npm ci
npm run build
```

`npm run build` runs `scripts/prepare.mjs`, which makes `public/admin/`, then
`astro build`. These variables change it:

| Variable       | What it sets                                           | Default                                  |
| -------------- | ------------------------------------------------------ | ---------------------------------------- |
| `SITE_URL`     | the address of the site, for the sitemap and the links | `https://usql.app`                       |
| `CMS_AUTH_URL` | the service that signs an editor in to GitHub (W3)     | none, so an editor signs in with a token |

## The release notes

The page of the release notes shows `src/data/generated/releases.json`. The
build does not fetch the releases. A person makes the file again after each
release of usql, and commits it:

```bash
go run ./cmd/gen -releases
```

The command reads the releases of usql from the GitHub API. If the variable
`GITHUB_TOKEN` or `GH_TOKEN` holds a token, the command sends it, so that
GitHub allows more requests. It needs no token for the few requests that it
makes. The workflow does not run it yet.

## GitHub Pages

`.github/workflows/site.yml` runs `scripts/precommit.sh` for every push and
pull request. On `main`, it then publishes `dist/` to GitHub Pages, when
the variable `PAGES_ENABLED` of the repository is `true`.
`public/CNAME` holds `usql.app`, the custom domain. GitHub Pages must be
turned on for the repository, with "GitHub Actions" as the source, before a
maintainer sets the variable (W1).

## The editor

The editor is served at `/admin/`. A save is a commit to `main`, which starts
the workflow, so the site shows the change a few minutes later.

An editor signs in with a GitHub token that can write to `xo/usql-site`:
a fine-grained token with "Contents: read and write" for that repository.
"Sign In with GitHub" needs the service of W3.
