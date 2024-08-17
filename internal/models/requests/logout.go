package requests

// LogoutRequest swagger::model
type LogoutRequest struct {
	RefreshToken string
}

// LogoutResponse swagger::model
type LogoutResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
