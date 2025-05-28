package constants

type (
	RequestSource string
	APIStatus     string
	CachePrefix   string
)

const (
	RequestSourceAPI       RequestSource = "API"
	RequestSourceWebClient RequestSource = "WEB_CLIENT"

	APIStatusOk    APIStatus = "ok"
	APIStatusError APIStatus = "error"

	CachePrefixFlightHistory CachePrefix = "FH_"
)
