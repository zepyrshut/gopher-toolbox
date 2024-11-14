package handlers

const (
	InvalidRequest  string = "invalid_request"
	InternalError   string = "internal_error"
	RequestID       string = "request_id"
	NotFound        string = "not_found"
	Created         string = "created"
	Updated         string = "updated"
	Deleted         string = "deleted"
	Enabled         string = "enabled"
	Disabled        string = "disabled"
	Retrieved       string = "retrieved"
	ErrorCreating   string = "error_creating"
	ErrorUpdating   string = "error_updating"
	ErrorEnabling   string = "error_enabling"
	ErrorDisabling  string = "error_disabling"
	ErrorGetting    string = "error_getting"
	ErrorGettingAll string = "error_getting_all"
	InvalidEntityID string = "invalid_entity_id"
	NotImplemented  string = "not_implemented"

	UserUsernameKey       string = "user_username_key"
	UserEmailKey          string = "user_email_key"
	UsernameAlReadyExists string = "username_already_exists"
	EmailAlreadyExists    string = "email_already_exists"
	IncorrectPassword     string = "incorrect_password"
	ErrorGeneratingToken  string = "error_generating_token"
	LoggedIn              string = "logged_in"

	CategoryNameKey       string = "category_name_key"
	CategoryAlreadyExists string = "category_already_exists"

	ItemsNameKey      string = "items_name_key"
	NameAlreadyExists string = "name_already_exists"
)
