---
title: PRIVATEBIN-INIT
header: Privatebin Manual
footer: 1.0.0
date: Feb 14, 2026
section: 1
---
# NAME
**privatebin-init** – generate a configuration file

# SYNOPSIS
**privatebin init** [-h | -help] [-\-force] [-\-host=\<url\>]

# DESCRIPTION
Generate a configuration file with sensible defaults and write it to
the preferred default location at _~/.config/privatebin/config.json_.
If the **-\-config** flag is set, the file is written to that path
instead.

The generated configuration contains a single bin entry pointing to the
host specified by **-\-host** (defaulting to _https://privatebin.net_),
with the following defaults: expire set to _1day_, formatter set to
_plaintext_, and gzip enabled.

If a configuration file already exists at the target path, the command
will fail unless **-\-force** is specified.

# OPTIONS
**-h, -\-help**
: Show help message.

**-\-force**
: Overwrite existing configuration file.

**-\-host** \<url\>
: The host URL of the default privatebin instance (default:
  _https://privatebin.net_).

# EXAMPLES
Generate a default configuration file:

    $ privatebin init

Generate a configuration file for a custom instance:

    $ privatebin init --host https://bin.example.com

Overwrite an existing configuration file:

    $ privatebin init --force

Generate a configuration file at a custom path:

    $ privatebin init --config /path/to/config.json

# SEE ALSO
**privatebin.conf**(5)

# AUTHORS
Bryan Frimin.
