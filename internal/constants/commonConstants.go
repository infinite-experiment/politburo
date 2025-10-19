package constants

type (
	RequestSource        string
	APIStatus            string
	CachePrefix          string
	LogbookRequestFilter string
)

const (
	RequestSourceAPI       RequestSource = "API"
	RequestSourceWebClient RequestSource = "WEB_CLIENT"

	APIStatusOk    APIStatus = "ok"
	APIStatusError APIStatus = "error"

	CachePrefixFlightHistory CachePrefix = "FH_"
	CachePrefixLiveries      CachePrefix = "LIVERY_"
	CachePrefixVAConfig      CachePrefix = "VA_CFG_"
	CachePrefixExpertServer  CachePrefix = "EXPERT_SERVER_ID"
	CachePrefixWorldDetails  CachePrefix = "IF_WORLD_DETAILS"
	CacheKeyServers          CachePrefix = "LIVE_SERVERS"
	CachePrefixLiveFlights   CachePrefix = "LIVE_FLIGHTS_"
	CachePrefixFPL           CachePrefix = "LIVE_FPL_"
	CachePrefixUserFlights   CachePrefix = "UFH_"

	FilterUser      LogbookRequestFilter = "USER"
	FilterDiscordId LogbookRequestFilter = "DISCORD_ID"
)
