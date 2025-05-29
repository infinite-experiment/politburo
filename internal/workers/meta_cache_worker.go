package workers

import (
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"time"
)

func StartCacheFiller(c *common.CacheService, api *common.LiveAPIService) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	refillAirframeCacheTask(c, api)
	refillWorldStatus(c, api)

	for range ticker.C {
		refillWorldStatus(c, api)
		refillAirframeCacheTask(c, api)
	}
}

func refillAirframeCacheTask(c *common.CacheService, api *common.LiveAPIService) {
	resp, _, err := api.GetAircraftLiveries()

	if err != nil {
		return
	}
	for _, liv := range resp.Liveries {
		c.Set(string(constants.CachePrefixLiveries)+liv.LiveryId, liv, 2400)

	}
}

func refillWorldStatus(c *common.CacheService, api *common.LiveAPIService) {
	resp, err := api.GetSessions()

	if err != nil {
		return
	}

	c.Set(string(constants.CachePrefixWorldDetails), resp.Result, 60000)
	for _, world := range resp.Result {
		// Get expert server
		if world.WorldType == 3 {
			c.Set(string(constants.CachePrefixExpertServer), world.ID, 60000)
			break
		}
	}
}
