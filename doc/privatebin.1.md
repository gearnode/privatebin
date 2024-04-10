---
title: PRIVATEBIN
header: Privatebin Manual
footer: 1.0.0
date: Jan 20, 2022
section: 1
---
# NAME
**privatebin** – manage privatebin pastes with simple shell command

# SYNOPSIS
**privatebin** [-h | -\-help] [-v | -\-version] [-\-bin=\<name\>]\
\ \ \ \ \ \ \ \ \ \ \ [-\-config=\<filename\>] [-\-header=\<key=value\>]\
\ \ \ \ \ \ \ \ \ \ \ [-\-output=\<format\>] \<command\> [\<args\>]

# DESCRIPTION
A minimalist, open source command line interface for **PrivateBin**
instances.

# OPTIONS
**-h, -\-help**
: Show help message.

**-v, --version**
: Prints the privatebin cli version.

**-b, -\-bin** \<name\>
: The privatebin instance name.

**-c, -\-config** \<path\>
: The path of the configuration file (default
  "~/.config/privatebin/config.json").

**-H, -\-header** \<key=value\>
: The extra HTTP header fields to include in the request sent.

**-o, -\-output** \<format\>
: The output format can be \"\" or \"json\" (default \"\").

# COMMANDS

**privatebin-create(1)**
: Create a paste

**privatebin-show(1)**
: Show a paste

# EXIT STATUS
The **privatebin** utility exits 0 on success, and >0 if an error
occurs.

# EXAMPLES
Create a paste on the default privatebin instance:

    $ cat example.txt | privatebin create

# SEE ALSO
**privatebin.conf**(5)

# AUTHORS
Bryan Frimin.
