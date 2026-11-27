# D7. The reference is generated from usql

Status: Proposed.

Proposed on 2026-11-27. The reference pages, which list the options, the meta
commands, the settings, the variables and the schemes of usql, are generated
from usql and dburl. A person does not write them.

## Why

usql adds and changes options, commands and settings in most releases. A
list that a person copies goes stale the day usql changes, and nothing tells
the reader that it is wrong. usql follows the same rule for its README, whose
tables of drivers its `gen.go` makes from dburl.

## How

A Go program of this repository will write the lists as JSON, and Astro pages
will draw them (usql-site W2). `go.mod` pins the release of usql that the
site documents. usql must export some of the lists first, which is open
question 4 of [PLAN.md](../PLAN.md).
