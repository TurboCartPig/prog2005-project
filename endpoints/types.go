package endpoints

type WebhookData struct {
	EventType string `json:"event_type"`
	User      User   `json:"user"`
	Project   struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		Description   string `json:"description"`
		WebURL        string `json:"web_url"`
		DefaultBranch string `json:"default_branch"`
	} `json:"project"`
	ObjectAttributes struct {
		Title               string        `json:"title"`
		Description         string        `json:"description"`
		AuthorID            int           `json:"author_id"`
		DueDate             string        `json:"due_date"`
		ClosedAt            interface{}   `json:"closed_at"`
		Confidential        bool          `json:"confidential"`
		CreatedAt           string        `json:"created_at"`
		DiscussionLocked    interface{}   `json:"discussion_locked"`
		ID                  int           `json:"id"`
		Iid                 int           `json:"iid"`
		LastEditedAt        interface{}   `json:"last_edited_at"`
		LastEditedByID      interface{}   `json:"last_edited_by_id"`
		MilestoneID         interface{}   `json:"milestone_id"`
		MovedToID           interface{}   `json:"moved_to_id"`
		DuplicatedToID      interface{}   `json:"duplicated_to_id"`
		ProjectID           int           `json:"project_id"`
		RelativePosition    int           `json:"relative_position"`
		StateID             int           `json:"state_id"`
		TimeEstimate        int           `json:"time_estimate"`
		TotalTimeSpent      int           `json:"total_time_spent"`
		UpdatedAt           string        `json:"updated_at"`
		UpdatedByID         int           `json:"updated_by_id"`
		Weight              interface{}   `json:"weight"`
		URL                 string        `json:"url"`
		AssigneeIds         []interface{} `json:"assignee_ids"`
		Labels              []Labels      `json:"labels"`
		State               string        `json:"state"`
		Action              string        `json:"action"`
	} `json:"object_attributes"`
	Labels     []Labels    `json:"labels"`
	Changes    interface{} `json:"changes"`
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
	Color       string      `json:"color"`
	ProjectID   int         `json:"project_id"`
	Description interface{} `json:"description"`
	Type        string      `json:"type"`
	GroupID     interface{} `json:"group_id"`
}
