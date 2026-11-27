---
title: Edit this site
description: Change a page of usql.app in the web editor or with a pull request.
sidebar:
  order: 1
---

The pages of usql.app are Markdown files in the repository [xo/usql-site](https://github.com/xo/usql-site), in `src/content/docs`. You can change them in two ways.

## In the web editor

The editor is at [usql.app/admin/](/admin/). It is Sveltia CMS, and it saves each change as a commit to the repository. You need write access to the repository.

## With a pull request

1. Fork [xo/usql-site](https://github.com/xo/usql-site).
2. Change the Markdown file of the page. The link "Edit page" at the end of each page opens it on GitHub.
3. Open a pull request.

## Recordings

A recording is a file in `public/casts`, which the program `cmd/record` of the repository writes from a script in `recordings`. A page shows it with a line such as:

```md
::cast[A first session with usql]{name="tour"}
```

Do not edit a recording by hand. Change its script and record it again.
