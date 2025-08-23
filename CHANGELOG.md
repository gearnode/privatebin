# Introduction

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.1.0] - 2025-08-15

### Added

- Add `skip-tls-verify` configuration option to skip TLS certificate verification.
- Add `--skip-tls-verify` flag to `create` and `show` commands.

### Changed

- Update dependencies.
- Upgrade go version from 1.22 to 1.23.

### Fixed

- Properly handle and log errors from `rootCmd.Execute()` in the main
  function.
- Extra headers not apply to show request.

## [2.0.1] - 2024-04-15

### Fixed

- Top level flags are not handled.

## [2.0.0] - 2024-04-11

### Added

- Add `privatebin show` command.
- Add `privatebin create` command.

### Changed

- Minimal Golang version is now v1.22.
- Minimal PrivateBin instance version is now 1.7.
- Configuration use kebab-case instead of sake-case.

## [1.4.0] - 2023-01-08

### Added

- Add `-gzip` flag to compress data with gzip.

### Changed

- According to OWAP recommendation, increase the number of PBKDF2
  iterations.

## [1.3.0] - 2022-11-06

### Added

- Add `-filename` flag to read file instead of stdin.
- Add `-attachment` flag to update data as an attachment.

### Changed

- Upgrade to Go 1.19.
- Use `gearno.de` import url.

### Fixed

- Create request error not handled.

## [1.2.0] - 2022-09-04

### Added

- Add privatebin version through the `-version` flag.

### Fixed

- Add `User-Agent` request header to mitigate WAF (Cloudflare, etc.)
  blocking request from the CLI.

## [1.1.1] - 2022-07-20

Nothing.

## [1.1.0] - 2022-06-23

### Added

- Add privatebin paste password support. Via the optional `-password`
  flag.

## [1.0.1] - 2022-01-20

### Fixed

- Missing URL path on the returned URL.

## [1.0.0] - 2021-09-06

### Added

- Add privatebin(1) man page.
- Add privatebin.conf(5) man page.

### Changed

- Makefile is now BSD and GNU compatible.
- Configuration file is now stored in the
  `~/.config/privatebin/config.json`.

## [0.1.0] - 2021-05-19

- First release.
