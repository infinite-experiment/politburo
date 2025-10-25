package workers

import (
	"context"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/db/repositories"
	gormModels "infinite-experiment/politburo/internal/models/gorm"
	"log"
	"time"
)

func StartCacheFiller(
	c *common.CacheService,
	api *common.LiveAPIService,
	liveryRepo *repositories.AircraftLiveryRepository,
	liverySvc *common.AircraftLiveryService,
) {
	ticker := time.NewTicker(6 * time.Hour) // 4x daily sync
	defer ticker.Stop()

	// Initial sync on startup
	syncAircraftLiveriesTask(c, api, liveryRepo, liverySvc)
	refillWorldStatus(c, api)

	for range ticker.C {
		refillWorldStatus(c, api)
		syncAircraftLiveriesTask(c, api, liveryRepo, liverySvc)
	}
}

// syncAircraftLiveriesTask syncs aircraft/livery data from IF API to database with change detection
func syncAircraftLiveriesTask(
	c *common.CacheService,
	api *common.LiveAPIService,
	liveryRepo *repositories.AircraftLiveryRepository,
	liverySvc *common.AircraftLiveryService,
) {
	ctx := context.Background()
	startTime := time.Now()

	// Fetch liveries from Infinite Flight API
	resp, _, err := api.GetAircraftLiveries()
	if err != nil {
		log.Printf("Error while fetching liveries from IF API: %s", err.Error())
		return
	}

	// Load existing liveries from database into map for change detection
	existingLiveries, err := liveryRepo.GetLiveryMap(ctx)
	if err != nil {
		log.Printf("Error while loading existing liveries from database: %s", err.Error())
		return
	}

	// Track changes
	var toUpsert []gormModels.AircraftLivery
	apiLiveryIDs := make(map[string]bool)
	addedCount := 0
	updatedCount := 0

	// Process each API livery
	for _, apiLivery := range resp.Liveries {
		apiLiveryIDs[apiLivery.LiveryId] = true

		if existingLivery, exists := existingLiveries[apiLivery.LiveryId]; exists {
			// Check if fields changed
			if existingLivery.AircraftName != apiLivery.AircraftName ||
				existingLivery.LiveryName != apiLivery.LiveryName ||
				existingLivery.AircraftID != apiLivery.AircraftID ||
				!existingLivery.IsActive {
				// Update needed
				toUpsert = append(toUpsert, common.ConvertAPILiveryToGORM(apiLivery))
				updatedCount++
			}
		} else {
			// New livery
			toUpsert = append(toUpsert, common.ConvertAPILiveryToGORM(apiLivery))
			addedCount++
		}
	}

	// Find removed liveries (in DB but not in API response)
	var removedIDs []string
	for liveryID := range existingLiveries {
		if !apiLiveryIDs[liveryID] {
			removedIDs = append(removedIDs, liveryID)
		}
	}

	// Execute database updates if changes detected
	hasChanges := len(toUpsert) > 0 || len(removedIDs) > 0

	if len(toUpsert) > 0 {
		if err := liveryRepo.UpsertBatch(ctx, toUpsert); err != nil {
			log.Printf("Error while upserting liveries: %s", err.Error())
			return
		}
	}

	if len(removedIDs) > 0 {
		if err := liveryRepo.MarkInactive(ctx, removedIDs); err != nil {
			log.Printf("Error while marking liveries inactive: %s", err.Error())
			return
		}
	}

	// Warm cache if changes detected (as per user requirement)
	if hasChanges {
		if err := liverySvc.WarmCache(ctx); err != nil {
			log.Printf("Error while warming livery cache: %s", err.Error())
		}
	}

	elapsed := time.Since(startTime)
	log.Printf(
		"Livery sync completed in %v: %d added, %d updated, %d removed (total in API: %d, in DB: %d)",
		elapsed,
		addedCount,
		updatedCount,
		len(removedIDs),
		len(resp.Liveries),
		len(existingLiveries),
	)
}

func refillWorldStatus(c *common.CacheService, api *common.LiveAPIService) {
	resp, err := api.GetSessions()

	if err != nil {
		return
	}

	c.Set(string(constants.CachePrefixWorldDetails), resp.Result, 60000*time.Minute)
	for _, world := range resp.Result {
		// Get expert server
		if world.WorldType == 3 {
			c.Set(string(constants.CachePrefixExpertServer), world.ID, 60000*time.Minute)
			break
		}
	}
}
