package services

import (
	stdContext "context"
	"database/sql"
	"errors"
	"fmt"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/context"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/models/entities"
	"log"
	"net/http"
	"sync"
)

type RegistrationService struct {
	Cache          common.CacheService
	LiveAPI        *common.LiveAPIService
	UserRepository repositories.UserRepository
}

func NewRegistrationService(liveAPI *common.LiveAPIService, cache common.CacheService, userRepo repositories.UserRepository) *RegistrationService {
	return &RegistrationService{
		LiveAPI:        liveAPI,
		Cache:          cache,
		UserRepository: userRepo,
	}
}

type InitRegistrationValidation struct {
	UserDB    *entities.User
	IFProfile *dtos.UserStats
}

func (svc *RegistrationService) InitUserRegistration(ctx stdContext.Context, ifcId string, lastFlight string) (*dtos.InitApiResponse, string, error) {

	var steps []dtos.RegistrationStep
	claims := context.GetUserClaims(ctx)

	// STEP - 1: API fetch & DB Fetch

	steps = append(steps, dtos.RegistrationStep{
		Name: "if_api&duplicacy_check", Status: true, Message: "Validated at Live API. User not duplicate.",
	})
	data, err := svc.UserValidation(ctx, ifcId)

	if err != nil {
		log.Printf("\n Registration Init Error: %v", err)
		steps[0].Message = err.Error()
		steps[0].Status = false
		return &dtos.InitApiResponse{
			IfcId:   ifcId,
			Status:  false,
			Message: "Unable to initialise registration",
			Steps:   steps,
		}, "Error", nil
	}

	if data.IFProfile == nil {
		steps[0].Message = "User not found at Live API"
		steps[0].Status = false
		return &dtos.InitApiResponse{
			IfcId:   ifcId,
			Status:  false,
			Message: "Infinite Flight Details not found for IFC Id",
			Steps:   steps,
		}, "Error", nil
	}

	// STEP 2 - Fetch User Flights
	steps = append(steps, dtos.RegistrationStep{
		Name: "if_flight_history", Status: true, Message: "Fetched flight history",
	})

	page := 1
	routeStr := ""
outer:
	for {
		fltResp, _, err := svc.LiveAPI.GetUserFlights(data.IFProfile.UserID, page)

		if err != nil {
			log.Printf("%v", err)
			steps[1].Message = err.Error()
			steps[1].Status = false
			return &dtos.InitApiResponse{
				IfcId:   ifcId,
				Status:  false,
				Message: "Failed to fetch user flights",
				Steps:   steps,
			}, "Unable to find user flight history", nil
		}

		for c := 0; c < len(fltResp.Flights); c++ {
			or := fltResp.Flights[c].OriginAirport
			de := fltResp.Flights[c].DestinationAirport
			if or != "" && de != "" {
				routeStr = fmt.Sprintf("%s-%s", or, de)
				break outer
			}
			log.Printf("Flight Route: %s-%s", fltResp.Flights[c].OriginAirport, fltResp.Flights[c].DestinationAirport)
		}

		// Max 3 searches
		if page > 2 {
			log.Printf("%v", err)
			steps[1].Message = "No recent flight found"
			steps[1].Status = false
			return &dtos.InitApiResponse{
				IfcId:   ifcId,
				Status:  false,
				Message: "Failed to fetch user flights",
				Steps:   steps,
			}, "Unable to find user flight history", nil
		}

		page++
	}

	steps = append(steps, dtos.RegistrationStep{
		Name: "user_check", Status: true, Message: "Last Flight validated.",
	})

	if routeStr == "" || routeStr != lastFlight {
		if routeStr == "" {
			steps[2].Message = "No flights found"
		} else {
			steps[2].Message = "Last flight did not match! Please check logbook"
		}
		steps[2].Status = false
		return &dtos.InitApiResponse{
			IfcId:   ifcId,
			Status:  false,
			Message: "Logbook flight did not match",
			Steps:   steps,
		}, "Unable to find user flight history", nil
	}

	if data.UserDB != nil {
		steps[1].Status = false
		steps[1].Message = "Duplicate request for user registration!"
		return &dtos.InitApiResponse{
			IfcId:   ifcId,
			Status:  false,
			Message: "User already registered",
			Steps:   steps,
		}, "User already Present", nil
	}

	insData := &entities.User{
		DiscordID:     claims.UserID(),
		IsActive:      true,
		IFCommunityID: ifcId,
		IFApiID:       &data.IFProfile.UserID,
	}
	log.Printf("Initiating insert: \n %v", *insData)

	if err := svc.UserRepository.InsertUser(ctx, insData); err != nil {
		log.Printf("Error: %v", err)
		return &dtos.InitApiResponse{
			IfcId:   ifcId,
			Status:  false,
			Message: "Unable to insert",
		}, "Unable to Insert", nil
	}

	return &dtos.InitApiResponse{
		IfcId:   ifcId,
		Status:  true,
		Message: "User has been registered",
	}, "", nil

}

func (svc *RegistrationService) UserValidation(ctx stdContext.Context, ifcId string) (*InitRegistrationValidation, error) {

	claims := context.GetUserClaims(ctx)
	var (
		user        *entities.User
		dbErr       error
		apiErr      error
		statusToken int
		statsResp   *dtos.UserStatsResponse
		wg          sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		user, dbErr = svc.UserRepository.FindUserByDiscordId(ctx, claims.UserID())
	}()

	// 2) external API call
	go func() {
		defer wg.Done()
		statsResp, statusToken, apiErr = svc.LiveAPI.GetUserByIfcId(ifcId)
	}()

	wg.Wait()

	// 1) DB: ignore ErrNoRows
	if dbErr != nil {
		if !errors.Is(dbErr, sql.ErrNoRows) {
			return nil, dbErr
		}
		user = nil
	}

	// 2) API: ignore 404
	if apiErr != nil {
		if statusToken == http.StatusNotFound {
			statsResp = nil
		} else {
			return nil, apiErr
		}
	}

	if len(statsResp.Result) == 0 {
		return nil, errors.New("no user found")
	}

	log.Printf("API Called successfully. Status: %d. Response %v", statusToken, *statsResp)
	log.Printf("DB Result: %v", user)
	// 3) return whatever you got (nil if missing)
	return &InitRegistrationValidation{
		UserDB:    user,
		IFProfile: &statsResp.Result[0],
	}, nil
}
