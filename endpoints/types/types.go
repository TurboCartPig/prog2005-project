package types

type WebhookData struct {
	EventType string `json:"event_type"`
	User      User   `json:"user"`
	Project   struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		WebURL      string `json:"web_url"`
	} `json:"project"`
	ObjectAttributes struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		AuthorID    int      `json:"author_id"`
		DueDate     string   `json:"due_date"`
		ProjectID   int      `json:"project_id"`
		Labels      []Labels `json:"labels"`
		URL         string   `json:"url"`
	} `json:"object_attributes"`
	Labels     []Labels `json:"labels"`
	Repository struct {
		Name        string `json:"name"`
		URL         string `json:"url"`
		Description string `json:"description"`
		Homepage    string `json:"homepage"`
	} `json:"repository"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Labels struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	ProjectID   int         `json:"project_id"`
	Description interface{} `json:"description"`
	Type        string      `json:"type"`
	GroupID     interface{} `json:"group_id"`
}

type ChannelRegistration struct {
	ChannelID  string `json:"channel_id"`
	RepoWebURL string `json:"repo_web_url"`
}

type Deadline struct {
	RepoWebURL  string `json:"repo_web_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IssueWebURL string `json:"issue_web_url"`
	DueDate     string `json:"due_date"`
}
