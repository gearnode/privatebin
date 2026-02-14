---
title: PRIVATEBIN.CONF
header: Privatebin Manual
footer: 1.0.0
date: Jan 20, 2022
section: 5
---

# NAME

**privatebin.conf** – privatebin CLI configuration file.

# DESCRIPTION

The privatebin(1) command line interface create paste to an PrivateBin
instance configured in the **config.json**.

# FORMAT

## Top level object keys:

**open-discussion** _bool_ (default: false)
: The default value of open discussion for a paste.

**burn-after-reading** _bool_ (default: false)
: The default value of burn after reading for a paste.

**formatter** _string_ (default: "plaintext")
: The default formatter for a paste.

**expire** _string_ (default: "1day")
: The default time to live for a paste.

**gzip** _bool_ (default: false)
: Enable GZip the paste data.

**skip-tls-verify** _bool_ (default: false)
: Skip TLS certificate verification when connecting to the privatebin instance.

**proxy** _string_
: Proxy URL to use for all requests. Supports HTTP, HTTPS, and SOCKS5
  schemes (e.g. "socks5://127.0.0.1:9050" for TOR). When set, overrides
  the **HTTP_PROXY**, **HTTPS_PROXY**, and **ALL_PROXY** environment
  variables. Can be overridden per-bin or by the **-\-proxy** CLI flag.

**extra-header-fields** _object<string, string>_
: The extra HTTP header fields to include in the request sent.

**bin** _array\<bin\>_
: The list of bin instances.

## The bin object format:

**name** _string_
: The name of the bin instance.

**host** _string_
: The url of the bin instance.

**auth** _auth_
: The basic auth configuration of the bin instance.

**expire** _string_
: The default time to live for a paste.

**open-discussion** _bool_
: The default value of open discussion for a paste.

**burn-after-reading** _bool_
: The default value of burn after reading for a paste.

**formatter** _string_
: The formatter for the paste.

**gzip** _bool_
: GZip the paste data.

**skip-tls-verify** _bool_
: Skip TLS certificate verification when connecting to the privatebin instance.

**proxy** _string_
: Proxy URL to use for requests to this bin instance. Supports HTTP,
  HTTPS, and SOCKS5 schemes. Overrides the top-level **proxy** value
  and the proxy environment variables. Can be overridden by the
  **-\-proxy** CLI flag.

**extra-header-fields** _object<string, string>_
: The extra HTTP header fields to include in the request sent.

## The auth object format:

**username** _string_
: The basic auth username.

**password** _string_
: The basic auth password.

# EXAMPLES

Minimal privatebin configuration file:

    {
        "bin": [
            {
                "name": "", // default
                "host": "https://privatebin.net"
            }
        ]
    }

A bit more complete configuration file:

    {
        "bin": [
            {
                "name": "example",
                "host": "bin.example.com",
                "auth": {
                    "username": "john.doe",
                    "password": "s$cr$t"
                },
                "formatter": "markdown",
                "burn-after-reading": false
            },
            {
                "name": "",
                "host": "https://privatebin.net"
    			"extra-header-fields": {
    				"Foo": "Bar",
    			},
            },
        ],
        "burn-after-reading": true
    }

Configuration using a SOCKS5 proxy (e.g. TOR):

    {
        "proxy": "socks5://127.0.0.1:9050",
        "bin": [
            {
                "name": "",
                "host": "https://privatebin.net"
            }
        ]
    }

# FILES

_~/.config/privatebin/config.json_
: Default location of the privatebin configuration. The file has to be
created manually as it is not installed with a standard installation.
The configuration file is optional; when absent, the CLI uses sensible
defaults (expire: 1day, formatter: plaintext, gzip: enabled).

# AUTHORS

Bryan Frimin.
