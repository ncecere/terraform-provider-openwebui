## 2.0.0 - 2025-09-20

### Added
- GitHub Actions workflow that builds and publishes tagged releases to the Terraform Registry.
- CHANGELOG tracking notable provider updates.

### Changed
- Provider examples now target version `~> 2.0` and demonstrate the structured `params`/`capabilities` attributes.
- Default build version in the `Makefile` set to `2.0.0` for local compilation.

### Removed
- `data_json` / `meta_json` arguments from `openwebui_group` to align with current Open WebUI APIs.

### Fixed
- Prompt commands are normalised with a leading `/` for all API calls, preventing mismatches when creating, reading, or deleting prompts.
- Group member emails are sorted to ensure stable Terraform plans.
- Tests and documentation updated to reflect the newer schema expectations.

## 0.1.0 - 2025-09-20
- Initial experimental release of the Open WebUI Terraform provider.
