---
title: Terminal graphics
description: How usql finds out which images a terminal draws, how to choose the graphics yourself, and the terminals that draw them.
sidebar:
  order: 8
---

usql draws images in a terminal that can show them. `\chart` draws a chart, and usql shows its logo when an interactive session starts. usql draws images only in an interactive session at a terminal. The [README of usql](https://github.com/xo/usql#readme) says how to draw them from a script.

A terminal draws images with a graphics protocol. A graphics protocol is a set of escape codes that carry an image to the terminal. usql knows three of them:

| Protocol | The name in usql |
| --- | --- |
| [The Kitty graphics protocol](https://sw.kovidgoyal.net/kitty/graphics-protocol.html) | `kitty` |
| [The iTerm image protocol](https://iterm2.com/documentation-images.html) | `iterm` |
| [Sixel](https://saitoha.github.io/libsixel/) | `sixel` |

## How usql finds out

When usql starts at a terminal, it asks the terminal which protocol it draws. It waits up to 250 ms for the answer. In tmux, it asks tmux about the terminal of the client. If the terminal does not answer, usql reads `TERM`, `TERM_PROGRAM` and the other environment variables that a terminal sets.

Keys that you type while usql waits go to the prompt. usql does not lose them.

## Choose the protocol yourself

The `USQL_TERM_GRAPHICS` and `TERM_GRAPHICS` environment variables choose the protocol, and usql does not ask the terminal. Set one of them to `kitty`, `iterm`, `sixel` or `none`. `none` turns the images off. An empty value lets usql find out.

If usql draws nothing, or draws garbage, name the protocol of your terminal:

```sh
USQL_TERM_GRAPHICS=kitty usql
```

To turn the images off:

```sh
USQL_TERM_GRAPHICS=none usql
```

If both variables are set, `USQL_TERM_GRAPHICS` wins over `TERM_GRAPHICS`.

## The term_graphics setting

The setting `term_graphics` holds the protocol of the session. `\uset term_graphics` shows it:

```txt
(not connected)=> \uset term_graphics
term_graphics is kitty
```

`\uset term_graphics` with a name changes the protocol for the session:

```sql
\uset term_graphics sixel
```

usql uses no images when the setting holds a name that is not a protocol.

## Terminals that draw images

These terminals work with usql:

| Terminal | Systems | Protocol |
| --- | --- | --- |
| [WezTerm](https://wezfurlong.org/wezterm/) | Windows, macOS, Linux and others | `iterm` |
| [iTerm2](https://iterm2.com) | macOS | `iterm` |
| [kitty](https://sw.kovidgoyal.net/kitty/) | Linux, macOS and BSD | `kitty` |
| [foot](https://codeberg.org/dnkl/foot) | Linux and other Wayland systems | `sixel` |

[Are We Sixel Yet?](https://www.arewesixelyet.com) lists more terminals that draw Sixel images.

[Your first chart](/charts/first-chart/) shows how to draw a chart.
