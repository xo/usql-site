# Decisions

Every decision this project made is a file in this folder, one per decision,
named by its number and its title. This table is the index. Find the number
here, then open the file.

Each file opens with its status. "Decided" means a maintainer chose it.
"Proposed" means an agent or a peer session suggested it and no maintainer has
confirmed it. "Open" means nobody has chosen yet. A decision that changes an
earlier one says so in its status, as "Amends D1", and the earlier one says it
back, as "Amended by D3". Read the status before the decision.

A new decision gets the next number and a file of its own, named
`D<nnn>-<title>.md` with the number in three digits. Add its row here.
`TestTheDecisionIndexIsComplete` fails when a decision has no row or a row is
wrong.

usql, `dburl` and `dbmeta` keep their own decision logs with the same
numbering. A reference to one of theirs names the repository, as in "usql's
D6".

| #                                                                          | Decision                                                     | Status   |
| -------------------------------------------------------------------------- | ------------------------------------------------------------ | -------- |
| [D1](D001-every-xo-repository-is-set-up-for-coding-agents-the-same-way.md) | Every xo repository is set up for coding agents the same way | Decided  |
| [D2](D002-the-site-is-static-and-its-pages-are-in-this-repository.md)      | The site is static, and its pages are in this repository     | Decided  |
| [D3](D003-astro-and-starlight-build-the-site.md)                           | Astro and Starlight build the site                           | Proposed |
| [D4](D004-a-page-is-markdown-and-a-recording-is-a-directive.md)            | A page is Markdown, and a recording is a directive           | Proposed |
| [D5](D005-a-program-makes-the-recordings-from-scripts.md)                  | A program makes the recordings from scripts                  | Proposed |
| [D6](D006-the-pages-load-nothing-from-another-host.md)                     | The pages load nothing from another host                     | Proposed |
| [D7](D007-the-reference-is-generated-from-usql.md)                         | The reference is generated from usql                         | Proposed |
