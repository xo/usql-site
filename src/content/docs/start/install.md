---
title: Install usql
description: Install usql from a release, with Homebrew, the AUR, Go or a container.
sidebar:
  order: 1
---

usql is one program, `usql` (`usql.exe` on Windows). Every way to install it below gives you every database driver. Pick the tab for your system.

:::tabs{sync="install"}
::::tab[Release]
1. Download the file for your system from the [releases of usql](https://github.com/xo/usql/releases).
2. Extract `usql` or `usql.exe` from the `.tar.bz2` or `.zip` file.
3. Move the program to a directory in your `PATH` (Linux and macOS) or `%PATH%` (Windows).
::::
::::tab[Homebrew]
Homebrew works on macOS and Linux:

```sh
brew install xo/xo/usql
```

For the ODBC databases, add the tap and install with ODBC:

```sh
brew tap xo/xo
brew install --with-odbc usql
```
::::
::::tab[Arch Linux]
On Arch Linux, install usql from the AUR. `usql-bin` installs the released program:

```sh
yay -S usql-bin
```

`usql` builds it from source, which takes much longer because it compiles every driver:

```sh
yay -S usql
```
::::
::::tab[Go]
```sh
go install github.com/xo/usql@latest
```
::::
::::tab[Container]
The image `docker.io/usql/usql` works with Docker, Podman and the other container runtimes. This command opens the SQLite file `data/test.db` of the current directory:

```sh
docker run --rm -it --volume "$(pwd)/data:/data" docker.io/usql/usql:latest sqlite3://data/test.db
```
::::
:::

## Make sure that it works

```sh
usql --version
```

usql prints its version. Continue with [First steps](/start/first-steps/).
