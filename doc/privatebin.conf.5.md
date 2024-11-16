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
**open-discussion** *bool* (default: false)
: The default value of open discussion for a paste.

**burn-after-reading** *bool* (default: false)
: The default value of burn after reading for a paste.

**formatter** *string* (default: "plaintext")
: The default formatter for a paste.

**expire** *string* (default: "1day")
: The default time to live for a paste.

**gzip** *bool* (default: false)
: Enable GZip the paste data.

**extra-header-fields** *object<string, string>*
: The extra HTTP header fields to include in the request sent.

**bin** *array\<bin\>*
: The list of bin instances.

## The bin object format:
**name** *string*
: The name of the bin instance.

**host** *string*
: The url of the bin instance.

**auth** *auth*
: The basic auth configuration of the bin instance.

**expire** *string*
: The default time to live for a paste.

**open-discussion** *bool*
: The default value of open discussion for a paste.

**burn-after-reading** *bool*
: The default value of burn after reading for a paste.

**formatter** *string*
: The formatter for the paste.

**gzip** *bool*
: GZip the paste data.

**extra-header-fields** *object<string, string>*
: The extra HTTP header fields to include in the request sent.

## The auth object format:
**username** *string*
: The basic auth username.

**password** *string*
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

# FILES
*~/.config/privatebin/config.json*
: Default location of the privatebin configuration. The file has to be
  created manually as it is not installed with a standard installation.

# AUTHORS
Bryan Frimin.
