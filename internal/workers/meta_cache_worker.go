package workers

import (
	"infinite-experiment/politburo/internal/common"
	"time"
)

func StartCacheFiller(c *common.CacheService, api *common.LiveAPIService) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	refillCacheTask(c, api)

	for range ticker.C {
		refillCacheTask(c, api)
	}
}

func refillCacheTask(c *common.CacheService, api *common.LiveAPIService) {
	resp, _, err := api.GetAircraftLiveries()

	if err != nil {
		return
	}
	// log.Printf("=========\n%v\n=======%v\n\n", resp, *resp)
	for _, liv := range resp.Liveries {
		c.Set("LIVERY_"+liv.LiveryId, liv, 2400)

	}
}
