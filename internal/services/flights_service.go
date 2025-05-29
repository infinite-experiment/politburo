package services

import (
	"fmt"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/workers"
)

type FlightsService struct {
	Cache      *common.CacheService
	ApiService *common.LiveAPIService
}

func NewFlightsService(cache *common.CacheService, liveApi *common.LiveAPIService) *FlightsService {

	return &FlightsService{
		Cache:      cache,
		ApiService: liveApi,
	}
}

func (svc *FlightsService) GetUserFlights(userId string, page int) (*dtos.FlightHistoryDto, error) {

	response := &dtos.FlightHistoryDto{
		PageNo:  page,
		Error:   "",
		Records: nil,
	}
	flt, _, err := svc.ApiService.GetUserByIfcId(userId)
	if err != nil || len(flt.Result) < 1 {
		response.Error = "Unable to fetch user"
		return response, err
	}

	uId := flt.Result[0].UserID
	flts, _, err := svc.ApiService.GetUserFlights(uId, page)

	if err != nil {
		response.Error = "Unable to fetch flights from Live API"
		return response, err
	}
	if len(flts.Flights) < 1 {
		response.Error = "No flights"
		return response, fmt.Errorf("empty result")
	}

	for _, rec := range flts.Flights {

		eqpmnt := GetAircraftLivery(rec.LiveryID, svc.Cache)

		aircraftName := ""
		liveryName := ""

		if eqpmnt != nil {
			aircraftName = eqpmnt.AircraftName
			liveryName = eqpmnt.LiveryName
		}

		dto := dtos.HistoryRecord{
			Origin:     rec.OriginAirport,
			Dest:       rec.DestinationAirport,
			TimeStamp:  rec.Created.UTC(),
			Landings:   rec.LandingCount,
			Server:     rec.Server,
			Equipment:  fmt.Sprintf("%s %s", GetShortAircraftName(aircraftName), GetShortLiveryName(liveryName)),
			Livery:     liveryName,
			Callsign:   rec.Callsign,
			Violations: len(rec.Violations),
		}
		select {
		case workers.LogbookQueue <- workers.LogbookRequest{FlightId: rec.ID}:
			dto.MapUrl = fmt.Sprintf("https://%s%s", "comrade.cc/i=", rec.ID)
			dto.MapUrl = ""
		default:
			dto.MapUrl = ""

		}
		response.Records = append(response.Records, dto)
	}

	return response, nil
}
