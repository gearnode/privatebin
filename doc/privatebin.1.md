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
\ \ \ \ \ \ \ \ \ \ \ [-\-output=\<format\>] [-\-proxy=\<url\>]\
\ \ \ \ \ \ \ \ \ \ \ \<command\> [\<args\>]

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
: The path of the configuration file. When not set, the CLI searches
  *$HOME/.config/privatebin/config.json*, then
  *$XDG\_CONFIG\_HOME/privatebin/config.json*, the platform-native user
  configuration directory, and finally the directories listed in
  **$XDG\_CONFIG\_DIRS** (see **privatebin.conf**(5) for details).

**-H, -\-header** \<key=value\>
: The extra HTTP header fields to include in the request sent.

**-o, -\-output** \<format\>
: The output format can be \"\" or \"json\" (default \"\").

**-\-proxy** \<url\>
: Proxy URL to use for requests. Supports HTTP, HTTPS, and SOCKS5
  schemes (e.g. socks5://127.0.0.1:9050 for TOR). This flag overrides
  the proxy value from the configuration file and the **HTTP_PROXY**,
  **HTTPS_PROXY**, and **ALL_PROXY** environment variables.

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

Create a paste through a SOCKS5 proxy (e.g. TOR):

    $ cat example.txt | privatebin --proxy socks5://127.0.0.1:9050 create

# ENVIRONMENT

**XDG\_CONFIG\_HOME**
: User-specific configuration directory. When set, the CLI searches
  this directory for *privatebin/config.json* (after checking
  *$HOME/.config*). See **privatebin.conf**(5) for the full search
  order.

**XDG\_CONFIG\_DIRS**
: Colon-separated list of system-wide configuration directories.
  Defaults to */etc/xdg* when not set. Each directory is searched for
  *privatebin/config.json*.

**HTTP_PROXY**, **HTTPS_PROXY**, **ALL_PROXY**
: When no **-\-proxy** flag is provided and no **proxy** configuration
  value is set, the standard proxy environment variables are honored.

**NO_PROXY**
: A comma-separated list of host names or IP addresses for which the
  proxy should not be used.

# SEE ALSO
**privatebin.conf**(5)

# AUTHORS
Bryan Frimin.
