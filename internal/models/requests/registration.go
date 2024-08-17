package requests

// RegisterRequest swagger::model
type RegisterRequest struct {
	Email               string `json:"email"`
	Password            string `json:"password"`
	SmartToyFingerPrint string `json:"fingerprint"`
}

// RegisterRequest swagger::model
type RegisterResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
