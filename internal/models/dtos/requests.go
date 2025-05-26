package dtos

type InitUserRegistrationReq struct {
	IfcId      string `json:"ifc_id"`
	LastFlight string `json:"last_flight"`
}

type LiveApiUserStatsReq struct {
	DiscourseNames []string `json:"discourseNames"`
}
