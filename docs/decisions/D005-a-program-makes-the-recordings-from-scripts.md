# D5. A program makes the recordings from scripts

Status: Proposed.

Proposed on 2026-11-27. Every terminal recording of the site is made by
`cmd/record` from a script in `recordings/`, and is never recorded by hand.

## Why

A recording shows usql as it was on the day it was made. usql changes its
prompt, its output and its colors between releases, so a recording made by
hand goes stale, and nobody can make it again the same way. A script makes the
recording again for each release with one command.

The program runs the script in a new temporary home directory, so that no
configuration, history or prompt of the person who records it reaches the
recording. The address of a database comes from a variable of the
environment, and the program stops when a recording shows its value.

The recordings are asciicast v2 files, which asciinema-player plays as text.
Text is small, sharp at every size, and can be selected and copied. A video
or a GIF has none of that.

## Costs

The player draws no images, so a recording of `\chart` cannot show the chart.
A chart page shows the chart as a separate SVG file (usql-site W5).

Typing at a fixed speed looks less natural than a person. Each step can wait
or sleep to give the reader time to read.
