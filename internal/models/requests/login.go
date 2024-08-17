package requests

// LoginRequest swagger::model
type LoginRequest struct {
	Email               string `json:"email"`
	Password            string `json:"password"`
	SmartToyFingerPrint string `json:"fingerprint"`
}

// LoginResponse swagger::model
type LoginResponse struct {
	Message      string `json:"message"`
	Error        string `json:"error"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
