package workers

import (
	"fmt"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/dtos"
	"time"
)

type LogbookRequest struct {
	FlightId  string
	Flight    dtos.UserFlightEntry
	SessionId string
}

var LogbookQueue = make(chan LogbookRequest, 100)

func LogbookWorker(c *common.CacheService, liveApiService *common.LiveAPIService) {
	for req := range LogbookQueue {

		cacheKey := string(constants.CachePrefixFlightHistory) + req.FlightId

		if req.Flight.DestinationAirport == "" || req.Flight.OriginAirport == "" {
			fmt.Println("❌ Skipping: missing origin or destination")
			continue
		}
		if val, found := c.Get(cacheKey); found && val != nil {
			fmt.Printf("⚠️  Flight %s already cached, skipping\n", req.FlightId)
			continue
		}

		session := req.SessionId
		// Call the Live API
		data, _, err := liveApiService.GetFlightRoute(req.FlightId, session)
		if err != nil {
			fmt.Printf("❌ Error fetching flight path: %v\n", err)
			continue
		}

		var waypoints []dtos.RouteWaypoint
		for _, pos := range data.Result {
			waypoints = append(waypoints, dtos.RouteWaypoint{
				Lat:       fmt.Sprintf("%f", pos.Latitude),
				Long:      fmt.Sprintf("%f", pos.Longitude),
				Altitude:  int(pos.Altitude),
				Timestamp: pos.Date,
			})
		}

		originNode := dtos.RouteNode{
			Name: req.Flight.OriginAirport,
			Lat:  "",
			Long: "",
		}
		destNode := dtos.RouteNode{
			Name: req.Flight.DestinationAirport,
			Lat:  "",
			Long: "",
		}

		if len(waypoints) > 0 {
			originNode.Lat = waypoints[0].Lat
			originNode.Long = waypoints[0].Long
			destNode.Lat = waypoints[len(waypoints)-1].Lat
			destNode.Long = waypoints[len(waypoints)-1].Long
		}
		hours := int(req.Flight.TotalTime) / 60
		minutes := int(req.Flight.TotalTime) % 60

		eqpmnt := common.GetAircraftLivery(req.Flight.LiveryID, c)

		aircraftName := ""
		liveryName := ""

		if eqpmnt != nil {
			aircraftName = eqpmnt.AircraftName
			liveryName = eqpmnt.LiveryName
		}

		flightInfo := dtos.FlightInfo{
			Meta: dtos.FlightMeta{
				Aircraft:   aircraftName,
				Livery:     liveryName,
				MaxSpeed:   0,
				Violations: len(req.Flight.Violations),
				Landings:   req.Flight.LandingCount,
				Duration:   hours*60 + minutes,
				StartedAt:  req.Flight.Created,
			},
			Route:  waypoints,
			Origin: originNode,
			Dest:   destNode,
		}

		c.Set(cacheKey, flightInfo, 600000*time.Second)
		fmt.Printf("✅ Cached %d points for flight %s\nCACHE_KEY=%s", len(data.Result), req.FlightId, cacheKey)

	}
}
