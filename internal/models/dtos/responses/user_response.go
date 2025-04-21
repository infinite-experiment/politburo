package responses

type UserRegisterResponse struct {
	RegistrationStatus bool   `json:"status"`
	Error              string `json:"error,omitempty"`
}
