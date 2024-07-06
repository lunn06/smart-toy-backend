package models

// RegisterRequest swagger::model
type RegisterRequest struct {
	Email               string `json:"email"`
	Password            string `json:"password"`
	SmartToyFingerPrint string `json:"fingerprint"`
}
