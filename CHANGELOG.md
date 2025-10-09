# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

## [unreleased]
### Added
- add new event: keysend notifications
- filter for `forward_min_amount` and `invoice_min_amount` to skip events with
low values.
- new template variables `{{.FeeRate}}` for `forward_event`. See
[TEMPLATES.md](TEMPLATES.md) for details. 
- add new event: payment succeeded and rebalancing succeeded notifications
