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
	Prefix string `json:"prefix,omitempty" validate:"max=10"`
	Suffix string `json:"suffix,omitempty" validate:"max=10"`
	VAName string `json:"name" validate:"required"`
}
