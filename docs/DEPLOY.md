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
