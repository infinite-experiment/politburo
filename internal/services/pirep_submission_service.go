package services

import (
	"context"
	"encoding/json"
	"log"

	"infinite-experiment/politburo/internal/models/dtos"
)

// PirepSubmissionService handles PIREP submission logic
type PirepSubmissionService struct {
	// Dependencies will be added here as we expand the service
}

// NewPirepSubmissionService creates a new PirepSubmissionService
func NewPirepSubmissionService() *PirepSubmissionService {
	return &PirepSubmissionService{}
}

// SubmitPirep processes an incoming PIREP submission request
// For now, this logs the DTO and returns a success response
func (s *PirepSubmissionService) SubmitPirep(ctx context.Context, request *dtos.PirepSubmitRequest) (*dtos.PirepSubmitResponse, error) {
	// Log the incoming request as JSON for debugging
	requestJSON, _ := json.MarshalIndent(request, "", "  ")
	log.Printf("[PirepSubmissionService] Received PIREP submission:\n%s\n", string(requestJSON))

	// Log individual fields for easy debugging
	log.Printf("[PirepSubmissionService] Mode: %s", request.Mode)
	log.Printf("[PirepSubmissionService] Route ID: %s", request.RouteID)
	log.Printf("[PirepSubmissionService] Flight Time: %s", request.FlightTime)
	log.Printf("[PirepSubmissionService] Pilot Remarks: %s", request.PilotRemarks)
	log.Printf("[PirepSubmissionService] Fuel: %v", request.FuelKg)
	log.Printf("[PirepSubmissionService] Cargo: %v", request.CargoKg)
	log.Printf("[PirepSubmissionService] Passengers: %v", request.Passengers)

	// Return success response for now
	response := &dtos.PirepSubmitResponse{
		Success: true,
		Message: "PIREP submission received and logged successfully",
		PirepID: "TEMP-" + request.Mode, // Placeholder PIREP ID
	}

	return response, nil
}
