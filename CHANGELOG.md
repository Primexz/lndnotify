# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

## [Unreleased]
### Added
- New event: opening channel notification (@Feelancer21)
- New event: closing channel notification (@Feelancer21)
- New template variables `{{.CloseInitiator}}` and `{{.CloseType}}` for `channel_close_event`. (@Feelancer21)

### Fixed
### Changed
### Removed

## [1.3.1] - 2025-10-11

### Fixed
- LND Connection attempts are now performed indefinitely (Previously, it was automatically abandoned after 15 minutes) (@Primexz)


## [1.3.0] - 2025-10-11

### Added
- New event: on chain transaction (@Primexz)


## [1.2.1] - 2025-10-09

### Added
- ``--version`` command to get the current version of the app (@Primexz)


## [1.2.0] - 2025-10-09

### Added
- New event: payment succeeded and rebalancing succeeded notifications (@feelancer21)
- New event: keysend notifications (@Primexz)
- Filter for `forward_min_amount` and `invoice_min_amount` to skip events with
low values (@feelancer21)
- New template variables `{{.FeeRate}}` for `forward_event`. See
[TEMPLATES.md](TEMPLATES.md) for details. (@feelancer21)

### Fixed

- Peer event: use pubkey as fallback (@Primexz)

### Changed

- Internal: new ppm rate formatting (@feelancer21)