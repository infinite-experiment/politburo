package dtos

type InitUserRegistrationReq struct {
	IfcId      string  `json:"ifc_id"`
	LastFlight string  `json:"last_flight"`
	Callsign   *string `json:"callsign,omitempty"` // Optional: for VA servers, links user to VA with callsign
}

type LinkUserToVAReq struct {
	Callsign string `json:"callsign"` // Required: callsign number (1-5 digits)
}

type SyncUser struct {
	UserID   string `json:"user_id"`
	Callsign string `json:"callsign"`
}

type SetRole struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type LiveApiUserStatsReq struct {
	DiscourseNames []string `json:"discourseNames"`
}

type InitServerRequest struct {
	VACode         string `json:"va_code" validate:"required,min=3,max=5"`
	VAName         string `json:"name" validate:"required"`
	CallsignPrefix string `json:"callsign_prefix"`
	CallsignSuffix string `json:"callsign_suffix"`
}

type VAConfig map[string]string

type VAConfigKeys struct {
	ConfigKeys []string `json:"config_keys"`
	ConfVals   []string `json:"samp"`
}
