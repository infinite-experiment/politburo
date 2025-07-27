package dtos

type InitUserRegistrationReq struct {
	IfcId      string `json:"ifc_id"`
	LastFlight string `json:"last_flight"`
}

type LiveApiUserStatsReq struct {
	DiscourseNames []string `json:"discourseNames"`
}

type InitServerRequest struct {
	VACode string `json:"va_code"  validate:"required,min=3,max=4"`
	VAName string `json:"name" validate:"required"`
}

type VAConfig struct {
	ConfigKey string `json:"config_key" validate:"required"`
	ConfigVal string `json:"config_value" validate:"required"`
}

type VAConfigKeys struct {
	ConfigKeys []string `json:"config_keys"`
	ConfVals   []string `json:"samp"`
}
