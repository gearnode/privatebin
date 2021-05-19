# Introduction
This repository contains a CLI for privatebin server.

# Abstract
I am a big privatebin user, and I never found a clean CLI to deal with
it. It why I decided to built this project to simplify my day to day
workflow.

# Install
You can install the command line interface with:

	go get -u github.com/gearnode/privatebin
	
# Usage
You can create paste from file with:

	cat resume.txt | privatebin -bin demo

# Build
You can build the command line interface with:

	make build

# Documentation
The [handbook](doc/handbook.md) contains informations about various
aspects of the command line interface.

You can also use the standard Go documentation tool to read code
documentation, for example:

	go doc -all github.com/gearnode/privatebin


# Contact
If you find a bug or have any question, feel free to open a Github issue
or to contact me [by email](mailto:bryan@frimin.fr).

Please note that I do not currently review or accept any contribution.

# Licence
Released under the ISC license.

Copyright (c) 2020-2021 Bryan Frimin <bryan@frimin.fr>.

Permission to use, copy, modify, and/or distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
