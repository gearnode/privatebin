---
title: PRIVATEBIN-SHOW
header: Privatebin Manual
footer: 1.0.0
date: Jan 20, 2022
section: 1
---
# NAME
**privatebin-show** – show a paste

# SYNOPSIS
**privatebin show** [-h | -\-help] [-\-confirm-burn] [-\-insecure]\
\ \ \ \ \ \ \ \ \ \ \ \ \ \ \ \ [-\-password] \<url\>

# DESCRIPTION
Show paste.

# OPTIONS
**-h, -\-help**
: Show help message.

**-\-confirm-burn**
: Confirm paste opening. It will be deleted immediately afterwards.

**-\-insecure**
: Allow reading paste from untrusted instance.

**-\-password**
: The paste password when paste has a password.

# EXAMPLES
Show a paste on the default privatebin instance:

    $ privatebin show https://example.com/foobar#mk

# SEE ALSO
**privatebin.conf**(5)

# AUTHORS
Bryan Frimin.
