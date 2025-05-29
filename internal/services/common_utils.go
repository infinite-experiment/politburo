package services

import (
	"fmt"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/dtos"
	"log"
	"strings"
	"time"
)

func GetResponseTime(init time.Time) string {
	timeDiff := time.Since(init).Milliseconds()
	return fmt.Sprintf("%dms", timeDiff)
}

func GetAircraftLivery(livId string, c *common.CacheService) *dtos.AircraftLivery {
	val, res := c.Get("LIVERY_" + livId)

	log.Printf("Querying for %s\n%v", "LIVERY_"+livId, val)
	if !res {
		return nil
	}

	eqpmnt, ok := val.(dtos.AircraftLivery)
	if !ok {
		return nil
	}

	return &eqpmnt
}

func GetShortAircraftName(fullName string) string {
	if short, ok := constants.AircraftShortNames[fullName]; ok {
		return short
	}
	// fallback to first 4 uppercase characters
	runes := []rune(fullName)
	if len(runes) > 4 {
		return strings.ToUpper(string(runes[:4]))
	}
	return strings.ToUpper(fullName)
}

func GetShortLiveryName(name string) string {
	if code, ok := constants.LiveryShortNames[name]; ok {
		return code
	}
	// fallback
	runes := []rune(name)
	if len(runes) > 4 {
		return strings.ToUpper(string(runes[:4]))
	}
	return strings.ToUpper(name)
}
