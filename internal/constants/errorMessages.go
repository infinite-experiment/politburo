package constants

const (
	StatusError             = "Error"
	StatusAlreadyPresent    = "User Already Present"
	StatusFailedToFetch     = "Failed to fetch user flights"
	StatusLogbookMismatch   = "Logbook flight did not match"
	StatusUserNotFound      = "Infinite Flight Details not found for IFC Id"
	StatusAlreadyRegistered = "User already registered"
	StatusInsertFailed      = "Unable to insert"
	StatusRegistrationInit  = "Unable to initialise registration"
	StatusRegistered        = "User has been registered"
)

const (
	MsgUserNotFound          = "User not found at Live API"
	MsgNoFlightsFound        = "No flights found"
	MsgFlightLogbookMismatch = "Last flight did not match! Please check logbook"
	MsgDuplicateRequest      = "Duplicate request for user registration!"
)
