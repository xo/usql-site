# D1. Every xo repository is set up for coding agents the same way

Status: Decided.

The maintainer decided on 2026-09-27 that every repository in the `xo`
namespace is set up for coding agents the same way. dbmeta records the
standard as its D110, and usql adopts it as its D1. This decision adopts it
for usql-site.

## The files

In the root:

- `README.md`, for a person who finds the project.
- `AGENTS.md`, which holds the rules for a coding agent.
- `CLAUDE.md`, which holds one line, `@AGENTS.md`, so that Claude Code reads
  the same rules. It is a file and not a symbolic link.
- `CONTRIBUTING.md`, for a person who changes the project, with the command
  that installs the agent skills.
- `skills-lock.json`, which names the source of each skill.
- `.gitignore`, which ignores `.claude/settings.local.json`.
- `.gitattributes`, which holds `* text=auto eol=lf`.

In `docs/`: `PLAN.md`, `BACKLOG.md`, and `decisions/` with one file per
decision and an index.

## The skills

The repository commits `simple-english` and `go-pedantry`, each as an ordinary
folder in `.agents/skills/<name>` and in `.claude/skills/<name>`.
`TestSkillsAreCopies` fails on a link, on a missing copy and on two copies
that differ.

## The standing rules

`AGENTS.md` opens with the three rules that are the same in every repository:
stage and do not commit until a maintainer says so, load `simple-english`
before writing prose, and load `go-pedantry` before writing Go.

## The tests

`internal/repo` checks the files: `TestClaudeImportsAgents`,
`TestTheRootHoldsFourDocuments`, `TestSkillsAreCopies` and
`TestTheDecisionIndexIsComplete`.
