# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)

## [unreleased]
### Added
- filter for `forward_min_amount` and `invoice_min_amount` to skip events with
low values.
- new template variables for `forward_event`: `{{.AmountIn}}`, `{{.AmountOut}}`
, `{{.FeeRate}}`. `{{.Amount}}` was removed. See [TEMPLATES.md](TEMPLATES.md) 
for details. 
