---
layout: page
page_title: "OpenWebUI Provider Coverage"
description: |-
  Coverage summary for Open WebUI API resources and data sources.
---

# Coverage Summary

This page outlines which Open WebUI API areas are covered by the provider.

## Resources

* Knowledge bases (`openwebui_knowledge`)
* Knowledge file attachments (`openwebui_knowledge_file`)
* Knowledge directories (`openwebui_knowledge_dir`)
* Models (`openwebui_model`)
* Prompts (`openwebui_prompt`)
* Groups (`openwebui_group`)
* Tools (`openwebui_tool`)
* Tool valves (`openwebui_tool_valves`)
* Functions (`openwebui_function`)
* Function valves (`openwebui_function_valves`)
* Pipelines (`openwebui_pipeline`)
* Pipeline valves (`openwebui_pipeline_valves`)
* Files (`openwebui_file`)
* Folders (`openwebui_folder`)
* Config import (`openwebui_config_import`)
* Admin configs:
  * Connections (`openwebui_connections_config`)
  * Tool servers (`openwebui_tool_servers_config`)
  * Terminal servers (`openwebui_terminal_servers_config`)
  * Code execution (`openwebui_code_execution_config`)
  * Models (`openwebui_models_config`)
  * Suggestions (`openwebui_suggestions_config`)
  * Banners (`openwebui_banners_config`)
  * Retrieval (`openwebui_retrieval_config`)
  * Evaluations (`openwebui_evaluations_config`)
  * Default user permissions (`openwebui_default_user_permissions`)
* OAuth clients (`openwebui_oauth_client`)

## Data Sources

* Knowledge (`openwebui_knowledge`, `openwebui_knowledge_files`)
* Models (`openwebui_model`, `openwebui_models`, `openwebui_base_models`, `openwebui_model_tags`, `openwebui_models_export`)
* Prompts (`openwebui_prompt`, `openwebui_prompts`, `openwebui_prompt_tags`, `openwebui_prompt_history`, `openwebui_prompt_history_entry`, `openwebui_prompt_history_diff`)
* Groups (`openwebui_group`)
* Tools (`openwebui_tool`, `openwebui_tools`, `openwebui_tools_export`)
* Functions (`openwebui_function`, `openwebui_functions`)
* Pipelines (`openwebui_pipeline`)
* Files (`openwebui_file`, `openwebui_files`, `openwebui_file_content`, `openwebui_file_process_status`)
* Folders (`openwebui_folder`)
* Config export (`openwebui_config_export`)
* Users (`openwebui_user`)
* Admin config data sources (`openwebui_retrieval_config`, `openwebui_evaluations_config`, `openwebui_default_user_permissions`)
* Tool server verification (`openwebui_tool_server_verify`)
* Terminal server verification (`openwebui_terminal_server_verify`)

## Notes

* Pipeline list payloads are loosely typed in the OpenAPI spec; the provider preserves them as raw JSON.
* The suggestions config endpoint does not expose a read API, so state is maintained from the last apply.
* Config export/import payloads can contain secrets; treat `config_json` as sensitive.
* OAuth client registration is write-only; state is preserved from the last apply.

## Not Yet Covered

* User-generated content resources such as chats, notes, or evaluations.
* Task automation endpoints (autocomplete, titles, tags, etc.).
* Media generation endpoints (audio, images).
