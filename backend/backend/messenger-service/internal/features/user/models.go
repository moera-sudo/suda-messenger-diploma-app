package user

type InitData struct {
	User         InitDataUser `json:"user"`
	AuthDate     int64        `json:"auth_date"`
	ChatType     string       `json:"chat_type,omitempty"`
	ChatInstance string       `json:"chat_instance,omitempty"`
	StartParam   string       `json:"start_param,omitempty"`
}

type InitDataUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Username    string `json:"username,omitempty"`
	PhotoURL    string `json:"photo_url,omitempty"`
}
