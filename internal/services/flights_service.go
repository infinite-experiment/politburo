package services

import (
	"fmt"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/workers"
	"log"
	"strings"
	"time"
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
		}
		if rec.OriginAirport != "" && rec.DestinationAirport != "" && rec.TotalTime > 0 && time.Since(rec.Created) <= 72*time.Hour {
			select {
			case workers.LogbookQueue <- workers.LogbookRequest{FlightId: rec.ID, Flight: rec}:
				dto.MapUrl = fmt.Sprintf("https://%s%s", "comrade.cc?i=", rec.ID)
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
