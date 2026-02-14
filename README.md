<p align="center">
  <img src="doc/logo.png" alt="PrivateBin CLI Logo" height="200"/>
  <br>
  <i>A command-line tool for creating and managing PrivateBin pastes.</i>
</p>

<p align="center">
  <img alt="GitHub License" src="https://img.shields.io/github/license/gearnode/privatebin">
  <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/gearnode/privatebin?label=version">
</p>

## Overview

A CLI tool for interacting with [PrivateBin](https://privatebin.info/),
the secure and anonymous paste service. It integrates directly into
your terminal workflow, letting you create, retrieve, and manage
pastes without leaving the command line.

## Installation

`privatebin` can be installed via a package manager, from a prebuilt
binary, or from source.

### macOS

```
brew install privatebin-cli
```

### Arch Linux

[![privatebin-cli on AUR](https://img.shields.io/aur/version/privatebin-cli?label=privatebin-cli)](https://aur.archlinux.org/packages/privatebin-cli/)
[![privatebin-cli-bin on AUR](https://img.shields.io/aur/version/privatebin-cli-bin?label=privatebin-cli-bin)](https://aur.archlinux.org/packages/privatebin-cli-bin/)

Available on the Arch User Repository (AUR). Install using your
preferred AUR helper:

- [privatebin-cli](https://aur.archlinux.org/packages/privatebin-cli/) - Release package
- [privatebin-cli-bin](https://aur.archlinux.org/packages/privatebin-cli-bin) - Binary package

#### Example

```console
yay -Sy privatebin-cli
```

### Ubuntu / Debian

```
apt-get install privatebin-cli
```

### Prebuilt Binary

Prebuilt binaries are available for a variety of operating systems and
architectures. Visit the [latest release](https://github.com/gearnode/privatebin/releases/latest)
page and scroll down to the Assets section.

1. Download the archive for your operating system and architecture
2. Extract the archive
3. Move the executable to a directory in your `PATH`
4. Verify that you have execute permission on the file

### Build from Source

1.  Clone the repository:

        git clone https://github.com/gearnode/privatebin.git

2.  Navigate to the project directory:

        cd privatebin

3.  Build the binary and man pages:

        make

4.  Install them on your system:

        make install

## Usage

Create a paste from a file:

    cat resume.txt | privatebin create

Display a paste:

    privatebin show https://privatebin.net/?420fc9597328c72f#EezApNVTTRUuEkt1jj7r9vSfewLBvUohDSXWuvPEs1bF

Create a paste through a SOCKS5 proxy (e.g., Tor):

    cat resume.txt | privatebin --proxy socks5://127.0.0.1:9050 create

## Documentation

For detailed information on all commands and options, see the
[handbook](doc/handbook.md).

## Support

Found a bug or have a question? Open a
[GitHub issue](https://github.com/gearnode/privatebin/issues) or
contact me via [email](mailto:bryan@frimin.fr).

## License

This project is released under the ISC license. See
[LICENSE.txt](LICENSE.txt) for details.
