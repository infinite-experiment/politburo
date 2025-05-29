package common

import (
	"fmt"
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

func GetAircraftLivery(livId string, c *CacheService) *dtos.AircraftLivery {
	val, res := c.Get(string(constants.CachePrefixLiveries) + livId)

	if !res {
		return nil
	}

	eqpmnt, ok := val.(dtos.AircraftLivery)
	if !ok {
		return nil
	}

	return &eqpmnt
}

func GetSessionId(c *CacheService, server string) *string {
	val, found := c.Get(string(constants.CachePrefixWorldDetails))
	if !found {
		return nil
	}

	if sessions, ok := val.([]dtos.Session); ok {
		for _, session := range sessions {
			if server == session.Name {
				return &session.ID
			}
		}
	}
	return nil
}

func GetFlightFromCache(c *CacheService, flightId string) *dtos.FlightInfo {

	log.Printf("\nFinding key: %s\n", string(constants.CachePrefixFlightHistory)+flightId)
	val, found := c.Get(string(constants.CachePrefixFlightHistory) + flightId)
	if !found {
		return nil
	}

	if flight, ok := val.(dtos.FlightInfo); ok {
		return &flight
	}
	return nil
}

func GetExpertServer(c *CacheService) *string {
	val, found := c.Get(string(constants.CachePrefixExpertServer))
	if !found {
		return nil
	}

	if serverID, ok := val.(string); ok {
		return &serverID
	}
	return nil
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
