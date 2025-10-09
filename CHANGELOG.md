# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

## [unreleased]
### Added
### Fixed
### Changed
### Removed

## [1.2.1] - 2025-10-09

### Added
- ``--version`` command to get the current version of the app


## [1.2.0] - 2025-10-09

### Added
- new event: payment succeeded and rebalancing succeeded notifications (@feelancer21)
- new event: keysend notifications (@Primexz)
- filter for `forward_min_amount` and `invoice_min_amount` to skip events with
low values (@feelancer21)
- new template variables `{{.FeeRate}}` for `forward_event`. See
[TEMPLATES.md](TEMPLATES.md) for details. (@feelancer21)

### Fixed

- peer event: use pubey as fallback

### Changed

- internal: new ppm rate formatting