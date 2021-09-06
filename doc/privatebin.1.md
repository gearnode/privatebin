---
title: PRIVATEBIN
header: Privatebin Manual
footer: 1.0.0
date: Sep 05, 2021
section: 1
---
# NAME
**privatebin** – create privatebin paste with simple shell command

# SYNOPSIS
**privatebin** [-help] [-bin=\<name\>] [-cfg-file=\<filename\>]\
\ \ \ \ \ \ \ \ \ \ \ \[-burn-after-reading] [-expire=\<time\>] [-formatter=\<format\>]\
\ \ \ \ \ \ \ \ \ \ \ \[-open-discussion]

# DESCRIPTION
A minimalist, open source command line interface for **PrivateBin**
instances.

# OPTIONS
**-help**
: Show help message.

**-bin** \<name\>
: The privatebin instance name.

**-burn-after-reading**
: Delete the paste after reading.

**-cfg-file** \<path\>
: The path of the configuration file (default
  "~/.config/privatebin/config.json").

**-expire** \<time\>
: The time to live of the paste.

**-formatter** \<format\>
: The text formatter to use, can be plaintext, markdown or
  syntaxhighlighting.

**-open-discussion**
: Enable discussion on the paste.

# EXIT STATUS
The **privatebin** utility exits 0 on success, and >0 if an error
occurs.

# EXAMPLES
Create a paste on the default privatebin instance:

    $ cat example.txt | privatebin

# SEE ALSO
**privatebin.conf**(5)
