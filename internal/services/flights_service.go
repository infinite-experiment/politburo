package services

import (
	"context"
	"errors"
	"fmt"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/workers"
	"log"
	"regexp"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type FlightsService struct {
	Cache      *common.CacheService
	ApiService *common.LiveAPIService
	Cfg        *common.VAConfigService
}

const maxRouteWorkers = 8

func SplitCallsign(raw string) (variable, prefix, suffix string) {
	cs := strings.TrimSpace(raw)

	//----------------------------------------------------------------------
	// 1. strip suffix
	//----------------------------------------------------------------------
	switch {
	case strings.HasSuffix(strings.ToUpper(cs), " HEAVY"):
		suffix = "Heavy"
		cs = strings.TrimSpace(cs[:len(cs)-6])
	case strings.HasSuffix(strings.ToUpper(cs), " SUPER"):
		suffix = "Super"
		cs = strings.TrimSpace(cs[:len(cs)-6])
	default:
		reFlight := regexp.MustCompile(`(?i)\s+FLIGHT\s+OF\s+(\d+)$`)
		if m := reFlight.FindStringSubmatch(cs); len(m) == 2 {
			suffix = "Flight of " + m[1]
			cs = strings.TrimSpace(reFlight.ReplaceAllString(cs, ""))
		}
	}

	//----------------------------------------------------------------------
	// 2. split the remaining string
	//----------------------------------------------------------------------
	tokens := strings.Fields(cs)
	if len(tokens) == 0 {
		return "", "", suffix
	}

	variable = tokens[len(tokens)-1]
	if len(tokens) > 1 {
		prefix = strings.Join(tokens[:len(tokens)-1], " ")
	}
	return variable, prefix, suffix
}

func NewFlightsService(cache *common.CacheService, liveApi *common.LiveAPIService, cfgSvc *common.VAConfigService) *FlightsService {

	return &FlightsService{
		Cache:      cache,
		ApiService: liveApi,
		Cfg:        cfgSvc,
	}
}

const userTTL = 15 * time.Minute
const fltTTL = 15 * time.Minute

// -----------------------------------------------------------------------------
// 1) User-lookup by IFC ID  (GET /users?ifcId=…)
// -----------------------------------------------------------------------------
func (svc *FlightsService) getUserByIfcIDCached(ifcID string) (*dtos.UserStatsResponse, error) {
	cacheKey := "LIVE_USER_" + ifcID

	val, err := svc.Cache.GetOrSet(cacheKey, userTTL, func() (any, error) {
		resp, _, err := svc.ApiService.GetUserByIfcId(ifcID)
		return resp, err // ← store *value* (resp is a struct)
	})
	if err != nil {
		return nil, err
	}

	lookup, ok := val.(*dtos.UserStatsResponse) // cached value is already *ptr*
	if !ok {
		return nil, fmt.Errorf("cache assertion failed for %s", cacheKey)
	}
	return lookup, nil
}

// -----------------------------------------------------------------------------
// 2) Paged flight list for a userID  (GET /users/{id}/flights?page=n)
// -----------------------------------------------------------------------------
func (svc *FlightsService) getUserFlightsCached(userID string, page int) (*dtos.UserFlightsResponse, error) {
	cacheKey := fmt.Sprintf("LIVE_FLIGHTS_%s_%d", userID, page)

	val, err := svc.Cache.GetOrSet(cacheKey, fltTTL, func() (any, error) {
		resp, _, err := svc.ApiService.GetUserFlights(userID, page)
		return resp, err
	})
	if err != nil {
		return nil, err
	}

	flts, ok := val.(*dtos.UserFlightsResponse)
	if !ok {
		return nil, fmt.Errorf("cache assertion failed for %s", cacheKey)
	}
	return flts, nil
}

func (svc *FlightsService) GetUserFlights(userId string, page int, sID string) (*dtos.FlightHistoryDto, error) {

	response := &dtos.FlightHistoryDto{
		PageNo:  page,
		Error:   "",
		Records: nil,
	}
	flt, err := svc.getUserByIfcIDCached(userId)
	if err != nil || len(flt.Result) < 1 {
		response.Error = "Unable to fetch user"
		return response, err
	}

	uId := ""
	userFound := false
	for _, res := range flt.Result {
		log.Printf("Matching %s - %s", *res.DiscourseUsername, userId)
		if strings.EqualFold(*res.DiscourseUsername, userId) {
			userFound = true
			uId = res.UserID
			break
		}
	}

	if !userFound {
		response.Error = "Unable to fetch user"
		return response, err
	}

	flts, err := svc.getUserFlightsCached(uId, page)

	if err != nil {
		response.Error = "Unable to fetch flights from Live API"
		return response, err
	}
	if len(flts.Flights) < 1 {
		response.Error = "No flights"
		return response, fmt.Errorf("empty result")
	}

	for _, rec := range flts.Flights {

		eqpmnt := common.GetAircraftLivery(rec.LiveryID, svc.Cache)

		aircraftName := ""
		liveryName := ""

		if eqpmnt != nil {
			aircraftName = eqpmnt.AircraftName
			liveryName = eqpmnt.LiveryName
		}
		totalMinutes := int(rec.TotalTime)
		hours := totalMinutes / 60
		minutes := totalMinutes % 60
		dur := fmt.Sprintf("%02d:%02d", hours, minutes)

		dto := dtos.HistoryRecord{
			Origin:     rec.OriginAirport,
			Dest:       rec.DestinationAirport,
			TimeStamp:  rec.Created.UTC(),
			Landings:   rec.LandingCount,
			Server:     rec.Server,
			Equipment:  fmt.Sprintf("%s %s", common.GetShortAircraftName(aircraftName), common.GetShortLiveryName(liveryName)),
			Livery:     liveryName,
			Callsign:   rec.Callsign,
			Violations: len(rec.Violations),
			Duration:   dur,
			Aircraft:   aircraftName,
		}
		if rec.OriginAirport != "" && rec.DestinationAirport != "" && rec.TotalTime > 0 && time.Since(rec.Created) <= 72*time.Hour {
			select {
			case workers.LogbookQueue <- workers.LogbookRequest{FlightId: rec.ID, Flight: rec, SessionId: sID}:
				dto.MapUrl = fmt.Sprintf("https://%s%s", "comradebot.cc?i=", rec.ID)
				//dto.MapUrl = ""
			default:
				dto.MapUrl = ""

			}
		} else {
			dto.MapUrl = ""
		}
		response.Records = append(response.Records, dto)
	}

	return response, nil
}

func (svc *FlightsService) mapToLiveFlight(resp *dtos.FlightsResponse, sId string) *[]dtos.LiveFlight {

	if resp == nil || len(resp.Flights) == 0 {
		return nil
	}

	out := make([]dtos.LiveFlight, len(resp.Flights))
	for i, flt := range resp.Flights {
		cVar, pfx, sfx := SplitCallsign(flt.Callsign)

		// Last report
		lastReport, err := common.ParseLiveAPITime(flt.LastReport)

		if err != nil {
			fmt.Printf("Couldn't parse to time: %s", flt.LastReport)
		}

		uname := flt.Username
		if uname == "" {
			uname = "<hidden>"
		}

		alt := (int(flt.Altitude) / 100) * 100
		spd := int(flt.Speed)

		eqpmnt := common.GetAircraftLivery(flt.LiveryID, svc.Cache)
		acft, liv := "", ""

		if eqpmnt != nil {
			acft = eqpmnt.AircraftName
			liv = eqpmnt.LiveryName
		}
		// Make DTO
		out[i] = dtos.LiveFlight{
			ReportTime:     lastReport,
			Callsign:       flt.Callsign,
			CallsignSuffix: sfx,
			SessionID:      sId,
			CallsignVar:    cVar,
			CallsignPrefix: pfx,
			IsConnected:    flt.IsConnected,
			AircraftId:     flt.AircraftID,
			LiveryId:       flt.LiveryID,
			FlightID:       flt.FlightID,
			Username:       uname,
			UserID:         flt.UserID,
			AltitudeFt:     alt,
			SpeedKts:       spd,
			Aircraft:       acft,
			Livery:         liv,
		}
	}

	return &out
}
func MatchCallsignVar(variable, startsWith, endsWith string) bool {
	v := strings.ToUpper(variable)

	if startsWith != "" && !strings.HasPrefix(v, strings.ToUpper(startsWith)) {
		return false
	}
	if endsWith != "" && !strings.HasSuffix(v, strings.ToUpper(endsWith)) {
		return false
	}
	return true
}

func FilterFlights(in []dtos.LiveFlight, pfx, sfx string) []dtos.LiveFlight {
	out := make([]dtos.LiveFlight, 0, len(in)) // fresh backing array

	for _, f := range in {
		if MatchCallsignVar(f.CallsignVar, pfx, sfx) {
			out = append(out, f) // copies struct into 'out'
		}
	}
	return out
}
func (svc *FlightsService) GetLiveFlights(sId string) (*[]dtos.LiveFlight, error) {
	cacheKey := string(constants.CachePrefixLiveFlights) + sId
	val, err := svc.Cache.GetOrSet(cacheKey, 2*time.Minute, func() (any, error) {
		f, _, err := svc.ApiService.GetFlights(sId)

		if err != nil {
			return nil, err
		}

		flights := svc.mapToLiveFlight(f, sId)

		return *flights, nil

	})

	if err != nil {
		fmt.Printf("\nError while fetching live flights: %v", err)
		return nil, err
	}

	flts, ok := val.([]dtos.LiveFlight)
	if !ok {
		fmt.Printf("\nError while parsing flights")
		return nil, errors.New("unable to fetch live flights")
	}

	return &flts, nil
}

func (svc *FlightsService) getFPLCacheKey(ifSid string, flightId string) string {
	return string(constants.CachePrefixFPL) + ifSid + "_" + flightId
}

func (svc *FlightsService) GetFlightPlan(ifSid string, flightId string) (*dtos.FlightPlanResponse, error) {
	cacheKey := svc.getFPLCacheKey(ifSid, flightId)
	// log.Printf("\n\nGet FPL called. cacheKey: %s", cacheKey)
	val, err := svc.Cache.GetOrSet(cacheKey, 5*time.Minute, func() (any, error) {
		// log.Printf("\nFetching FPL. cacheKey: %s", cacheKey)

		fpl, _, err := svc.ApiService.GetFlightPlan(ifSid, flightId)
		if err != nil {
			log.Printf("Failed to fetch FPL: %v", err)
			return nil, err
		}

		// log.Printf("\nFetched FPL. Waypoints: %d", len(fpl.Waypoints))

		return *fpl, nil
	})
	if err != nil {
		return nil, err
	}

	fpl, ok := val.(dtos.FlightPlanResponse)

	if !ok {
		return nil, fmt.Errorf("failed to unmarshal FPL")
	}
	// log.Printf("Fetched FPL with waypoints %d for %s", len(fpl.Waypoints), fpl.FlightID)
	return &fpl, nil
}

// Returns origin and Destination airports
func (svc *FlightsService) GetFlightRoute(ifSid string, flightId string) (string, string) {
	org, dest := "", ""

	fpl, err := svc.GetFlightPlan(ifSid, flightId)

	if err == nil {
		wayp := fpl.Waypoints

		wl := len(wayp)
		if wl > 1 {
			x := wayp[0]
			y := wayp[wl-1]
			// log.Printf("\nFiltering waypoints. Length: %d, First: %s, Last: %s", len(fpl.Waypoints), x, y)

			if len(x) == 4 {
				org = x
			}
			if len(y) == 4 {
				dest = y
			}
		}
	}

	return org, dest

}

func (svc *FlightsService) enrichFlightData(flts *[]dtos.LiveFlight) *[]dtos.LiveFlight {
	start := time.Now() // ← start timer
	defer func() {
		log.Printf("enriched %d flights in %v",
			len(*flts), time.Since(start))
	}()
	grp, sem := errgroup.Group{}, make(chan struct{}, maxRouteWorkers)

	for i := range *flts {
		i := i
		sem <- struct{}{}
		grp.Go(func() error {
			defer func() { <-sem }()

			f := &(*flts)[i]

			org, dest := svc.GetFlightRoute(f.SessionID, f.FlightID)
			f.Origin, f.Destination = org, dest
			return nil
		})
	}
	_ = grp.Wait()
	return flts
}

func (svc *FlightsService) GetVALiveFlights(ctx context.Context, vaId string) (*[]dtos.LiveFlight, error) {

	sId, ok := svc.Cfg.GetConfigVal(ctx, vaId, common.ConfigKeyIFServerID)

	if !ok {
		return nil, errors.New("no IF server configured for VA")
	}

	pfx, _ := svc.Cfg.GetConfigVal(ctx, vaId, common.ConfigKeyCallsignPrefix)
	sfx, _ := svc.Cfg.GetConfigVal(ctx, vaId, common.ConfigKeyCallsignSuffix)

	if pfx == "" && sfx == "" {
		return nil, errors.New("prefix and Suffix not configured for airline")
	}

	live_flt, err := svc.GetLiveFlights(sId)
	if err != nil {
		fmt.Printf("No live flights found with error: %v", err)
		return nil, err
	}

	va_flt := FilterFlights(*live_flt, pfx, sfx)

	va_flt = *svc.enrichFlightData(&va_flt)

	return &va_flt, nil

}

func (svc *FlightsService) GetLiveServers() (*[]dtos.Session, error) {
	const cacheKey = string(constants.CacheKeyServers)

	if val, found := svc.Cache.Get(cacheKey); found {
		if sessions, ok := val.(*[]dtos.Session); ok {
			return sessions, nil
		}
	}

	// Fetch fresh data
	data, err := svc.ApiService.GetSessions()
	if err != nil {
		return nil, err
	}

	sessions := &data.Result
	svc.Cache.Set(cacheKey, sessions, 5*time.Minute) // 5 minutes

	return sessions, nil
}
