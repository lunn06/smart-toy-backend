package requests

// LoginRequest swagger::model
type LoginRequest struct {
	Email               string `json:"email"`
	Password            string `json:"password"`
	SmartToyFingerPrint string `json:"fingerprint"`
}
