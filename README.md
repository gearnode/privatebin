<p align="center">
  <img src="doc/logo.png" alt="PrivateBin CLI Logo" height="200"/>
  <br>
  <i>A powerful CLI for creating and managing PrivateBin pastes with ease.</i>
</p>

<p align="center">
  <img alt="GitHub License" src="https://img.shields.io/github/license/gearnode/privatebin">
  <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/gearnode/privatebin?label=version">
</p>

## Overview

PrivateBin's secure and anonymous paste service is indispensable for
many developers and privacy enthusiasts. Recognizing the need for a
more efficient way to interact with PrivateBin from the terminal, I
developed this CLI tool. It's designed to seamlessly integrate with
your workflow, enabling swift creation and management of pastes.

## Installation

`privatebin` can be installed using a prebuilt binary, through package
managers, or from source. Follow the instructions below for your
preferred method.

### macOS

```
brew install privatebin-cli
```

### Arch Linux

[![privatebin-cli on AUR](https://img.shields.io/aur/version/privatebin-cli?label=privatebin-cli)](https://aur.archlinux.org/packages/privatebin-cli/)
[![privatebin-cli-bin on AUR](https://img.shields.io/aur/version/privatebin-cli-bin?label=privatebin-cli-bin)](https://aur.archlinux.org/packages/privatebin-cli-bin/)

Available on the Arch User Repository (AUR). Install using your
favorite AUR helper:

- [privatebin-cli](https://aur.archlinux.org/packages/privatebin-cli/) - Release package
- [privatebin-cli-bin](https://aur.archlinux.org/packages/privatebin-cli-bin) - Binary package

#### Example Installation:

```console
yay -Sy privatebin-cli
```

### Ubuntu / Debian

```
apt-get install privatebin-cli
```

### Prebuilt binary

Prebuilt binaries are available for a variety of operating systems and
architectures. Visit the latest release page, and scroll down to the
Assets section.

1. Download the archive for the desired edition, operating system, and architecture
2. Extract the archive
3. Move the executable to the desired directory
4. Add this directory to the PATH environment variable
5. Verify that you have execute permission on the file

### Build from Source

1.  Clone the repository:

    git clone https://github.com/gearnode/privatebin.git

2.  Navigate to the project directory:

        cd privatebin

3.  Build the project (binary and man pages):

        make

4.  Install the binary and man pages on your system:

        make install

## Usage

Create a paste from a file:

    cat resume.txt | privatebin create

Display a paste:

    privatebin show https://privatebin.net/?420fc9597328c72f#EezApNVTTRUuEkt1jj7r9vSfewLBvUohDSXWuvPEs1bF

## Documentation

For detailed information on all CLI commands and features, check out
the [handbook](doc/handbook.md).

## Support

Encountered a bug or have questions? Feel free to open a GitHub issue
or contact me directly via [email](mailto:bryan@frimin.fr).

## License

This project is released under the ISC license. See the
[LICENSE.txt](LICENSE.txt) file for details. It's designed with both
openness and freedom of use in mind, but with no warranty as per the
ISC standard disclaimer.
