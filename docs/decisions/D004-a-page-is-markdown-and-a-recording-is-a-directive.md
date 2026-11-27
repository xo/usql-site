# D4. A page is Markdown, and a recording is a directive

Status: Proposed.

Proposed on 2026-11-27. A page is a `.md` file, never `.mdx`. Anything that is
more than Markdown is a directive, such as:

    ::cast[A first session with usql]{name="tour"}

## Why

Sveltia CMS edits Markdown. MDX can import a component anywhere in a page,
and the editor can neither show nor keep such a page safely. A directive is
one line of text that the editor keeps as it is. An editor component in
`cms/editor.js` turns it into a form, so that a person who does not know the
syntax can add a recording.

Starlight already reads directives for its asides, such as `:::note`, so
the syntax is the same on every page.

## How

Each directive has two parts, which change together:

- a plugin in `src/lib/`, which turns the directive into HTML at build time,
  and stops the build when it is wrong, such as a recording that does not
  exist.
- an editor component in `cms/editor.js`, with the same pattern.
