# Changelog

## v3.0.0 - Open WebUI v0.9.6 compatibility release

### Added

- GLWTPL `LICENSE` file for registry/license metadata (addresses #4).
- Safe live E2E harness under `dev_testing/e2e` with default non-fixture cases.
- Folder management via `openwebui_folder` and `data.openwebui_folder`.
- Knowledge directory and file-directory support via `openwebui_knowledge_dir`, `openwebui_knowledge_files`, and `directory_id` on `openwebui_knowledge_file`.
- File enhancements: inline `content`, `target_filename`, file content data source, and process-status data source.
- Terminal server config resource and terminal server verification data source.
- Raw catalog data sources for models, base models, model tags, model export, prompts, prompt tags, tools, and tool export.
- Prompt history data sources for history list, history entry, and history diff.
- Function management via `openwebui_function`, `openwebui_function_valves`, `data.openwebui_function`, and `data.openwebui_functions`.
- Raw admin config resources/data sources for retrieval config, evaluations config, and default user permissions.

### Changed

- Refreshed `openapi.json` from a live Open WebUI v0.9.6 server.
- Prompt resource now supports `data_json`, `meta_json`, `tags`, `is_active`, and `commit_message`.
- Prompt API compatibility updated for v0.9.6 `name` payloads and ID-based read/update/delete routes (addresses #7 prompt creation failure).
- Model config behavior now preserves intentionally empty `model_order_list` state.
- Documentation, coverage notes, and index updated for v3 resources and data sources.

### Fixed

- Model data source schema mismatches.
- Model deletion request body compatibility with v0.9.6.
- Files list parsing for v0.9.6 `{ items, total }` responses.
- Group data source config decoding and stale schema mismatches.
- Group creation with no users/empty users now remains stable after apply (addresses #3).
- Model creation no longer requires explicit `capabilities` when omitted (addresses #7 capabilities failure).
- Model metadata merge/null handling now preserves dedicated fields correctly (addresses PR #5).
- Connections configuration management is available through `openwebui_connections_config` (addresses #2).

### Deferred / Known Gaps

- Fixture-backed E2E cases remain opt-in/deferred for pipelines, pipeline valves, tool/terminal server verification, OAuth client registration, and config import.
- MCP/tool server add/remove item-level testing is deferred.
- Some Open WebUI endpoints expose loose schemas; several provider surfaces intentionally use raw JSON to preserve API compatibility.
