package client

import (
	"context"
	"net/http"
	"net/url"
)

// ConnectionsConfigForm captures connection configuration.
type ConnectionsConfigForm struct {
	EnableDirectConnections bool `json:"ENABLE_DIRECT_CONNECTIONS"`
	EnableBaseModelsCache   bool `json:"ENABLE_BASE_MODELS_CACHE"`
}

// ImportConfigForm captures the config import payload.
type ImportConfigForm struct {
	Config map[string]any `json:"config"`
}

// OAuthClientRegistrationForm registers an OAuth client.
type OAuthClientRegistrationForm struct {
	URL        string  `json:"url"`
	ClientID   string  `json:"client_id"`
	ClientName *string `json:"client_name,omitempty"`
}

// ToolServerConnection captures a tool server connection entry.
type ToolServerConnection struct {
	URL      string         `json:"url"`
	Path     string         `json:"path"`
	Type     *string        `json:"type,omitempty"`
	AuthType *string        `json:"auth_type"`
	Headers  any            `json:"headers,omitempty"`
	Key      *string        `json:"key"`
	Config   map[string]any `json:"config"`
}

// ToolServersConfigForm captures tool server configuration.
type ToolServersConfigForm struct {
	Connections []ToolServerConnection `json:"TOOL_SERVER_CONNECTIONS"`
}

// TerminalServerConnection captures a terminal server connection entry.
type TerminalServerConnection struct {
	ID         *string        `json:"id,omitempty"`
	Name       *string        `json:"name,omitempty"`
	Enabled    *bool          `json:"enabled,omitempty"`
	URL        string         `json:"url"`
	Path       *string        `json:"path,omitempty"`
	Key        *string        `json:"key,omitempty"`
	AuthType   *string        `json:"auth_type,omitempty"`
	Config     map[string]any `json:"config,omitempty"`
	ServerType *string        `json:"server_type,omitempty"`
	PolicyID   *string        `json:"policy_id,omitempty"`
	Policy     map[string]any `json:"policy,omitempty"`
}

// TerminalServersConfigForm captures terminal server configuration.
type TerminalServersConfigForm struct {
	Connections []TerminalServerConnection `json:"TERMINAL_SERVER_CONNECTIONS"`
}

// CodeInterpreterConfigForm captures code execution configuration.
type CodeInterpreterConfigForm struct {
	EnableCodeExecution                bool    `json:"ENABLE_CODE_EXECUTION"`
	CodeExecutionEngine                string  `json:"CODE_EXECUTION_ENGINE"`
	CodeExecutionJupyterURL            *string `json:"CODE_EXECUTION_JUPYTER_URL"`
	CodeExecutionJupyterAuth           *string `json:"CODE_EXECUTION_JUPYTER_AUTH"`
	CodeExecutionJupyterAuthToken      *string `json:"CODE_EXECUTION_JUPYTER_AUTH_TOKEN"`
	CodeExecutionJupyterAuthPassword   *string `json:"CODE_EXECUTION_JUPYTER_AUTH_PASSWORD"`
	CodeExecutionJupyterTimeout        *int64  `json:"CODE_EXECUTION_JUPYTER_TIMEOUT"`
	EnableCodeInterpreter              bool    `json:"ENABLE_CODE_INTERPRETER"`
	CodeInterpreterEngine              string  `json:"CODE_INTERPRETER_ENGINE"`
	CodeInterpreterPromptTemplate      *string `json:"CODE_INTERPRETER_PROMPT_TEMPLATE"`
	CodeInterpreterJupyterURL          *string `json:"CODE_INTERPRETER_JUPYTER_URL"`
	CodeInterpreterJupyterAuth         *string `json:"CODE_INTERPRETER_JUPYTER_AUTH"`
	CodeInterpreterJupyterAuthToken    *string `json:"CODE_INTERPRETER_JUPYTER_AUTH_TOKEN"`
	CodeInterpreterJupyterAuthPassword *string `json:"CODE_INTERPRETER_JUPYTER_AUTH_PASSWORD"`
	CodeInterpreterJupyterTimeout      *int64  `json:"CODE_INTERPRETER_JUPYTER_TIMEOUT"`
}

// ModelsConfigForm captures default model configuration.
type ModelsConfigForm struct {
	DefaultModels       *string  `json:"DEFAULT_MODELS"`
	DefaultPinnedModels *string  `json:"DEFAULT_PINNED_MODELS"`
	ModelOrderList      []string `json:"MODEL_ORDER_LIST"`
}

// PromptSuggestion captures a default prompt suggestion.
type PromptSuggestion struct {
	Title   []string `json:"title"`
	Content string   `json:"content"`
}

// SetDefaultSuggestionsForm updates default suggestions.
type SetDefaultSuggestionsForm struct {
	Suggestions []PromptSuggestion `json:"suggestions"`
}

// BannerModel captures banner configuration.
type BannerModel struct {
	ID          string  `json:"id"`
	Type        string  `json:"type"`
	Title       *string `json:"title,omitempty"`
	Content     string  `json:"content"`
	Dismissible bool    `json:"dismissible"`
	Timestamp   int64   `json:"timestamp"`
}

// SetBannersForm updates banners configuration.
type SetBannersForm struct {
	Banners []BannerModel `json:"banners"`
}

// GetConnectionsConfig retrieves connections config.
func (c *Client) GetConnectionsConfig(ctx context.Context) (*ConnectionsConfigForm, error) {
	var resp ConnectionsConfigForm
	if err := c.do(ctx, http.MethodGet, "configs/connections", nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SetConnectionsConfig updates connections config.
func (c *Client) SetConnectionsConfig(ctx context.Context, form ConnectionsConfigForm) (*ConnectionsConfigForm, error) {
	var resp ConnectionsConfigForm
	if err := c.do(ctx, http.MethodPost, "configs/connections", nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetToolServersConfig retrieves tool server configuration.
func (c *Client) GetToolServersConfig(ctx context.Context) (*ToolServersConfigForm, error) {
	var resp ToolServersConfigForm
	if err := c.do(ctx, http.MethodGet, "configs/tool_servers", nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SetToolServersConfig updates tool server configuration.
func (c *Client) SetToolServersConfig(ctx context.Context, form ToolServersConfigForm) (*ToolServersConfigForm, error) {
	var resp ToolServersConfigForm
	if err := c.do(ctx, http.MethodPost, "configs/tool_servers", nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetTerminalServersConfig retrieves terminal server configuration.
func (c *Client) GetTerminalServersConfig(ctx context.Context) (*TerminalServersConfigForm, error) {
	var resp TerminalServersConfigForm
	if err := c.do(ctx, http.MethodGet, "configs/terminal_servers", nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SetTerminalServersConfig updates terminal server configuration.
func (c *Client) SetTerminalServersConfig(ctx context.Context, form TerminalServersConfigForm) (*TerminalServersConfigForm, error) {
	var resp TerminalServersConfigForm
	if err := c.do(ctx, http.MethodPost, "configs/terminal_servers", nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// VerifyTerminalServer verifies terminal server connectivity.
func (c *Client) VerifyTerminalServer(ctx context.Context, connection TerminalServerConnection) error {
	if err := c.do(ctx, http.MethodPost, "configs/terminal_servers/verify", nil, connection, nil); err != nil {
		return err
	}

	return nil
}

// GetCodeExecutionConfig retrieves code execution configuration.
func (c *Client) GetCodeExecutionConfig(ctx context.Context) (*CodeInterpreterConfigForm, error) {
	var resp CodeInterpreterConfigForm
	if err := c.do(ctx, http.MethodGet, "configs/code_execution", nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SetCodeExecutionConfig updates code execution configuration.
func (c *Client) SetCodeExecutionConfig(ctx context.Context, form CodeInterpreterConfigForm) (*CodeInterpreterConfigForm, error) {
	var resp CodeInterpreterConfigForm
	if err := c.do(ctx, http.MethodPost, "configs/code_execution", nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetModelsConfig retrieves models configuration.
func (c *Client) GetModelsConfig(ctx context.Context) (*ModelsConfigForm, error) {
	var resp ModelsConfigForm
	if err := c.do(ctx, http.MethodGet, "configs/models", nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SetModelsConfig updates models configuration.
func (c *Client) SetModelsConfig(ctx context.Context, form ModelsConfigForm) (*ModelsConfigForm, error) {
	var resp ModelsConfigForm
	if err := c.do(ctx, http.MethodPost, "configs/models", nil, form, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetBanners retrieves banners configuration.
func (c *Client) GetBanners(ctx context.Context) ([]BannerModel, error) {
	var resp []BannerModel
	if err := c.do(ctx, http.MethodGet, "configs/banners", nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// SetBanners updates banners configuration.
func (c *Client) SetBanners(ctx context.Context, banners []BannerModel) ([]BannerModel, error) {
	var resp []BannerModel
	form := SetBannersForm{Banners: banners}
	if err := c.do(ctx, http.MethodPost, "configs/banners", nil, form, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// SetDefaultSuggestions updates default suggestions configuration.
func (c *Client) SetDefaultSuggestions(ctx context.Context, suggestions []PromptSuggestion) ([]PromptSuggestion, error) {
	var resp []PromptSuggestion
	form := SetDefaultSuggestionsForm{Suggestions: suggestions}
	if err := c.do(ctx, http.MethodPost, "configs/suggestions", nil, form, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ExportConfig retrieves the full configuration export.
func (c *Client) ExportConfig(ctx context.Context) (map[string]any, error) {
	var resp map[string]any
	if err := c.do(ctx, http.MethodGet, "configs/export", nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ImportConfig applies a configuration export payload.
func (c *Client) ImportConfig(ctx context.Context, config map[string]any) (map[string]any, error) {
	var resp map[string]any
	form := ImportConfigForm{Config: config}
	if err := c.do(ctx, http.MethodPost, "configs/import", nil, form, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// RegisterOAuthClient registers an OAuth client with optional type.
func (c *Client) RegisterOAuthClient(ctx context.Context, form OAuthClientRegistrationForm, clientType *string) (map[string]any, error) {
	var resp map[string]any
	var query url.Values
	if clientType != nil {
		query = url.Values{}
		query.Set("type", *clientType)
	}

	if err := c.do(ctx, http.MethodPost, "configs/oauth/clients/register", query, form, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// VerifyToolServer verifies tool server connectivity.
func (c *Client) VerifyToolServer(ctx context.Context, connection ToolServerConnection) error {
	if err := c.do(ctx, http.MethodPost, "configs/tool_servers/verify", nil, connection, nil); err != nil {
		return err
	}

	return nil
}

// GetRetrievalConfig retrieves raw retrieval configuration.
func (c *Client) GetRetrievalConfig(ctx context.Context) (map[string]any, error) {
	var resp map[string]any
	if err := c.do(ctx, http.MethodGet, "retrieval/config", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SetRetrievalConfig updates raw retrieval configuration.
func (c *Client) SetRetrievalConfig(ctx context.Context, config map[string]any) (map[string]any, error) {
	var resp map[string]any
	if err := c.do(ctx, http.MethodPost, "retrieval/config/update", nil, config, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetEvaluationsConfig retrieves raw evaluations configuration.
func (c *Client) GetEvaluationsConfig(ctx context.Context) (map[string]any, error) {
	var resp map[string]any
	if err := c.do(ctx, http.MethodGet, "evaluations/config", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SetEvaluationsConfig updates raw evaluations configuration.
func (c *Client) SetEvaluationsConfig(ctx context.Context, config map[string]any) (map[string]any, error) {
	var resp map[string]any
	if err := c.do(ctx, http.MethodPost, "evaluations/config", nil, config, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetDefaultUserPermissions retrieves raw default user permissions.
func (c *Client) GetDefaultUserPermissions(ctx context.Context) (map[string]any, error) {
	var resp map[string]any
	if err := c.do(ctx, http.MethodGet, "users/default/permissions", nil, nil, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SetDefaultUserPermissions updates raw default user permissions.
func (c *Client) SetDefaultUserPermissions(ctx context.Context, config map[string]any) (map[string]any, error) {
	var resp map[string]any
	if err := c.do(ctx, http.MethodPost, "users/default/permissions", nil, config, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
