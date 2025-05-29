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
	CachePrefixLiveries      CachePrefix = "LIVERY_"
	CachePrefixExpertServer  CachePrefix = "EXPERT_SERVER_ID"
	CachePrefixWorldDetails  CachePrefix = "IF_WORLD_DETAILS"
)
