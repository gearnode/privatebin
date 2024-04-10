---
title: PRIVATEBIN-CREATE
header: Privatebin Manual
footer: 1.0.0
date: Jan 20, 2022
section: 1
---
# NAME
**privatebin-create** – create a paste

# SYNOPSIS
**privatebin create** [-h | -help]  [-\-burn-after-reading] [-\-expire=\<time\>]\
\ \ \ \ \ \ \ \ \ \ \ \ \ \ \ \ \ \ [-\-formatter=\<format\>] [-\-open-discussion]\
\ \ \ \ \ \ \ \ \ \ \ \ \ \ \ \ \ \ [-\-password=\<password\>] [-\-gzip] [-\-attachment] \
\ \ \ \ \ \ \ \ \ \ \ \ \ \ \ \ \ \ [-\-filename=\<filename\>] *STDIN*

# DESCRIPTION
Create paste.

# OPTIONS
**-h, -\-help**
: Show help message.

**-\-burn-after-reading**
: Delete the paste after reading.

**-\-expire** \<time\>
: The time to live of the paste.

**-\-formatter** \<format\>
: The text formatter to use, can be plaintext, markdown or
  syntaxhighlighting.

**-\-open-discussion**
: Enable discussion on the paste.

**-\-password**
: Add password on the paste.

**-\-attachment**
: Create the paste as an attachment.

**-\-filename**
: Open and read filename instead of `stdin`.

**-\-gzip**
: GZip the paste data.

# EXAMPLES
Create a paste on the default privatebin instance:

    $ cat example.txt | privatebin create

# SEE ALSO
**privatebin.conf**(5)

# AUTHORS
Bryan Frimin.
