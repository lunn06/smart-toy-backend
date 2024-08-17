package requests

// RefreshTokensRequest swagger::model
type RefreshTokensRequest struct {
	RefreshToken        string `json:"refreshToken"`
	SmartToyFingerPrint string `json:"fingerprint"`
}

// RefreshTokensResponse swagger::model
type RefreshTokensResponse struct {
	Message      string `json:"message"`
	Error        string `json:"error"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
