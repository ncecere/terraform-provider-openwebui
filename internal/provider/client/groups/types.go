package groups

type Group struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Permissions *GroupPermissions      `json:"permissions"`
	Data        map[string]interface{} `json:"data"`
	Meta        map[string]interface{} `json:"meta"`
	UserIDs     []string               `json:"user_ids"`
	CreatedAt   int64                  `json:"created_at"`
	UpdatedAt   int64                  `json:"updated_at"`
}

type GroupPermissions struct {
	Workspace WorkspacePermissions `json:"workspace"`
	Chat      ChatPermissions      `json:"chat"`
}

type WorkspacePermissions struct {
	Models    bool `json:"models"`
	Knowledge bool `json:"knowledge"`
	Prompts   bool `json:"prompts"`
	Tools     bool `json:"tools"`
}

type ChatPermissions struct {
	FileUpload bool `json:"file_upload"`
	Delete     bool `json:"delete"`
	Edit       bool `json:"edit"`
	Temporary  bool `json:"temporary"`
}
