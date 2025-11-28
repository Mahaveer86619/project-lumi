package views

type ErrorResponse struct {
	Message  string   `json:"error"`
	Session  string   `json:"session,"`
	Status   string   `json:"status"`
	Expected []string `json:"expected"`
}
