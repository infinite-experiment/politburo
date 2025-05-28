package services

import (
	"fmt"
	"infinite-experiment/politburo/internal/models/dtos"
)

type FlightsService struct {
	Cache      CacheService
	ApiService *LiveAPIService
}

func NewFlightsService(cache CacheService, liveApi *LiveAPIService) *FlightsService {

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
		response.Records = append(response.Records,
			dtos.HistoryRecord{
				Origin:    rec.OriginAirport,
				Dest:      rec.DestinationAirport,
				TimeStamp: rec.Created.UTC(),
				Landings:  rec.LandingCount,
				Server:    rec.Server,
				Aircraft:  rec.AircraftID,
				Livery:    rec.LiveryID,
				MapUrl:    fmt.Sprintf("https://%s%s", "vizburo.infinite-flight.com/user_id=", rec.ID),
			})
	}

	return response, nil
}
