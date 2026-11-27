# Contributing

## Setup

You need Node.js 22.23.4 or later, and Go for the local check, the
recordings and `cmd/gen`.

```bash
npm ci
npm run dev
```

## Change a page

Open the page in `src/content/docs/`, or in the editor at
`http://localhost:4321/admin/`. [docs/CONTENT.md](docs/CONTENT.md) says how
a page is written, and [docs/RECORDING.md](docs/RECORDING.md) says how to add
a terminal recording.

## Before you open a pull request

```bash
bash scripts/precommit.sh
```

The last line says whether the change passes. The check runs gofmt,
`go vet`, `go test`, the tests of the plugins of `src/lib/`, prettier,
`astro check`, the build, and a check of the built pages.

## Rules

- One pull request does one thing.
- A durable choice gets a decision in `docs/decisions/`.
- Known work goes in `docs/BACKLOG.md`, with the next `W` number.
- Never commit a password, a token, or the address of a database.

# Agent skills

The repository carries two agent skills. A skill is a set of instructions that
a coding agent loads for a task. `simple-english` sets how prose is written,
and `go-pedantry` sets how Go is written. The standing rules in
[AGENTS.md](AGENTS.md) say when to load each.

`skills-lock.json` names the source of each skill. The `skills` command from
npm writes that file, and version 1.7.0 is the one used here. It writes each
skill into two folders. Codex and the other agents read
`.agents/skills/<name>`, and Claude Code reads `.claude/skills/<name>`.

To add a skill or to update one, run this in the repository root:

```bash
npx skills@1.7.0 add AminBlg/SimpleEnglish --skill simple-english --agent codex claude-code --copy -y
```

```bash
npx skills@1.7.0 add oborchers/fractional-cto --skill go-pedantry --agent codex claude-code --copy -y
```

Keep `--copy`. Without it, the command writes `.claude/skills/<name>` as a
symbolic link. A Windows checkout writes a symbolic link as a text file, and
Claude Code then loads no skill and says nothing. `TestSkillsAreCopies` fails
on a link, on a missing copy, and when the two folders differ. See D1.

`.claude/settings.local.json` holds the Claude Code permissions of one person.
The root `.gitignore` ignores it.
