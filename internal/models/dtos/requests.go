package dtos

type InitUserRegistrationReq struct {
	IfcId      string `json:"ifc_id"`
	LastFlight string `json:"last_flight"`
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
	VACode string `json:"va_code"  validate:"required,min=3,max=4"`
	VAName string `json:"name" validate:"required"`
}

type VAConfig map[string]string

type VAConfigKeys struct {
	ConfigKeys []string `json:"config_keys"`
	ConfVals   []string `json:"samp"`
}
