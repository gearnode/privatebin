# Introduction

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
