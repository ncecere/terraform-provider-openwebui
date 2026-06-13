# OpenWebUI Terraform Provider v3.0.0 Plan

This plan targets Open WebUI v0.9.6 and a major provider release. The provider is already on `terraform-plugin-framework` (`v1.16.0`), so a migration from SDKv2 is not required; v3 should focus on API coverage, typed schemas, state migration, and acceptance coverage.

## Current state

Validated locally:

```bash
go test ./...
```

Result: all current unit/provider tests pass.

The live server in `.env` exposes `/openapi.json` with 430 paths. The repository `openapi.json` is currently malformed (`Extra data: line 39837 column 2`) and should be replaced with the fetched v0.9.6 spec before using it as a review artifact.

Currently implemented resources:

- `openwebui_knowledge`
- `openwebui_knowledge_file`
- `openwebui_model`
- `openwebui_prompt`
- `openwebui_group`
- `openwebui_tool`
- `openwebui_tool_valves`
- `openwebui_pipeline`
- `openwebui_pipeline_valves`
- `openwebui_file`
- `openwebui_config_import`
- `openwebui_connections_config`
- `openwebui_tool_servers_config`
- `openwebui_code_execution_config`
- `openwebui_models_config`
- `openwebui_suggestions_config`
- `openwebui_banners_config`
- `openwebui_oauth_client`

Currently implemented data sources:

- `openwebui_model`
- `openwebui_knowledge`
- `openwebui_group`
- `openwebui_prompt`
- `openwebui_tool`
- `openwebui_pipeline`
- `openwebui_file`
- `openwebui_files`
- `openwebui_config_export`
- `openwebui_user`
- `openwebui_tool_server_verify`

## High priority v3 additions

### 1. Refresh OpenAPI and coverage tooling

- Replace malformed `openapi.json` with the v0.9.6 spec fetched from the configured server.
- Add a small script that compares `openapi.json` paths against provider client/resource coverage.
- Regenerate `docs/coverage.md` from that comparison or at least make it spec-versioned.

### 2. Admin/config coverage

Add resources/data sources for v0.9.6 config endpoints not currently represented:

- `openwebui_terminal_servers_config`
  - `GET/POST /api/v1/configs/terminal_servers`
  - Optional verify data source for `/verify`
  - Optional policy resource for `/policy`
- `openwebui_auth_config`
  - `GET/POST /api/v1/auths/admin/config`
- `openwebui_ldap_config`
  - `GET/POST /api/v1/auths/admin/config/ldap`
- `openwebui_ldap_server_config`
  - `GET/POST /api/v1/auths/admin/config/ldap/server`
- `openwebui_retrieval_config`
  - `GET /api/v1/retrieval/config`
  - `POST /api/v1/retrieval/config/update`
- `openwebui_evaluations_config`
  - `GET/POST /api/v1/evaluations/config`
- `openwebui_default_user_permissions`
  - `GET/POST /api/v1/users/default/permissions`

### 3. Workspace resources

Add first-class resources/data sources for API objects with CRUD-style coverage:

- `openwebui_function`
  - create/update/delete/read/list/export/load-url/sync
  - valves and user valves as separate resources or nested blocks
- `openwebui_function_valves`
- `openwebui_folder`
  - manage chat/workspace folders, parent, expanded state
- `openwebui_channel`
  - manage channels and members
- `openwebui_channel_webhook`
  - manage channel webhooks
- `openwebui_memory`
  - manage explicit memories where suitable

Recommended v3 scope: implement `function`, `function_valves`, and `folder` first; treat channels/memories as stretch items because they are closer to user-generated runtime content.

### 4. Improve existing resources

- Knowledge:
  - Add directory resources (`/knowledge/{id}/dirs/create`, update, delete).
  - Add batch file attachment support.
  - Add optional reindex/sync operations as explicit resources or action-like resources.
- Prompts:
  - Add version/history data sources.
  - Add metadata update support if missing from current schema.
  - Confirm toggle/active semantics match v0.9.6.
- Models:
  - Add import/export/sync data sources/actions if safe.
  - Add tags data source.
  - Confirm access update payloads still match v0.9.6.
- Tools:
  - Add load-from-url support.
  - Add user valves data source/resource if useful.
- Files:
  - Add rename support.
  - Add content update support.
  - Add process status data source.

### 5. Type formerly opaque JSON

Major release opportunity: reduce `*_json` string fields where stable schemas exist.

Prioritize typed blocks for:

- connections config
- code execution config
- models config
- banners config
- tool/terminal server config
- retrieval config
- auth/LDAP config
- default user permissions

Keep raw JSON escape hatches as `advanced_json` or `raw_json` where Open WebUI is unstable.

### 6. State upgrades and breaking changes

Because this is v3.0.0:

- Add state upgrade handlers for renamed attributes/resources where possible.
- Document any resource address migrations.
- Mark sensitive config fields as `Sensitive` consistently.
- Normalize IDs and command/path inputs consistently.
- Review all computed/optional fields for plan drift.

### 7. Acceptance testing

Expand live acceptance tests using the server from `.env` via environment variables:

```bash
export OPENWEBUI_ENDPOINT="https://test.chat.ai.it.ufl.edu/api/v1"
export OPENWEBUI_TOKEN="..."
export TF_ACC=1
go test ./internal/provider -run TestAcc -v
```

Add tests for:

- functions + valves
- terminal servers config + verify
- auth/LDAP config read/update/read
- retrieval config read/update/read
- default user permissions
- folder CRUD
- file rename/content/status
- knowledge directories and batch file attachment

## Suggested implementation order

1. Replace/fix `openapi.json`; add spec coverage script.
2. Add terminal servers config because it mirrors existing tool servers config.
3. Add retrieval config and evaluations config.
4. Add auth/LDAP/default permissions configs.
5. Add `function` + `function_valves` resources/data sources.
6. Add `folder` resource/data source.
7. Improve existing knowledge/files/prompts/models/tools with missing v0.9.6 endpoints.
8. Update docs/examples and cut `v3.0.0` changelog.

## Release checklist

- [ ] `go test ./...`
- [ ] `gofmt ./...`
- [ ] Acceptance tests against v0.9.6 server
- [ ] Updated `README.md`, `docs/index.md`, `docs/coverage.md`, examples
- [ ] Upgrade guide from v2.x to v3.0.0
- [ ] Changelog entry for breaking changes and new resources
- [ ] Tag `v3.0.0`
