---
title: Line editing and completion
description: The keys of the usql prompt, the search of the history, Tab completion, and the inputrc file.
sidebar:
  order: 6
---

The usql prompt edits a line with the keys of GNU Readline, as the prompt of `psql` and of `bash` does. Tab completes meta commands, SQL words and the names of tables and columns.

::cast[A statement over three rows, Ctrl-P to bring it back, Ctrl-R to search the history, and Tab to complete a table name and a meta command.]{name="line-editing"}

## Statements over several rows

If you press Enter before the `;` that ends a statement, usql starts a new row of the same statement. The prompt shows `->` on the new rows. usql runs the statement when you press Enter after the `;`.

The rows of a statement are one entry of the history. Up brings back the whole statement, and you can edit any of its rows.

## Keys

These are the keys that you use most. The other keys of GNU Readline work too.

| Key | What it does |
| --- | --- |
| Left, Right | move the cursor one character |
| Alt-B, Alt-F | move the cursor one word back or forward |
| Ctrl-A, Ctrl-E | move the cursor to the start or the end of the line |
| Ctrl-W | delete the word before the cursor |
| Ctrl-K | delete from the cursor to the end of the line |
| Ctrl-U | delete from the start of the line to the cursor |
| Ctrl-\_ | undo the last change |
| Up or Ctrl-P, Down or Ctrl-N | show the previous or the next entry of the history |
| Ctrl-R | search the history |
| Tab | complete the word at the cursor |
| Ctrl-C | clear the line |
| Ctrl-D | leave usql, at an empty prompt |

## Search the history

Press Ctrl-R, then type a part of an earlier entry. usql shows the latest entry that holds it. Press Ctrl-R again for an older entry. Press Enter to run the entry. To edit the entry before you run it, press Ctrl-E. Ctrl-C stops the search and clears the line.

[Configure usql](/guides/configure/) says where usql keeps the history.

## Completion

Press Tab to complete the word at the cursor. If only one word fits, usql completes it. If more than one word fits, usql completes the part that they share and shows a numbered menu below the line:

```txt
sq:books.db=> \dt
 1 \dt      4 \dtx     7 \dtx+
 2 \dt+     5 \dtS+    8 \dtSx+
 3 \dtS     6 \dtSx
```

In the menu, press Tab or an arrow key to move to an entry. Press Enter or the number of an entry to take it. Press Ctrl-C to close the menu.

usql completes the names of tables and columns from the database that it is connected to. Not all databases and statements have completion.

## The inputrc file

usql reads the inputrc file of GNU Readline. The file is `~/.inputrc`, or the file that the `INPUTRC` environment variable names. Lines between `$if usql` and `$endif` apply only to usql.

To edit with the keys of vi, add this line:

```txt
set editing-mode vi
```

To make Up and Down show only the entries that start with what you typed, add these lines:

```txt
$if usql
"\e[A": history-search-backward
"\e[B": history-search-forward
$endif
```

The new lines apply the next time that you start usql.
