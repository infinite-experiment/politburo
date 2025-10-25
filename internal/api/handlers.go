package api

type Handlers struct {
	deps *Dependencies
}

// NewHandlers creates a new handlers instance with injected dependencies
func NewHandlers(deps *Dependencies) *Handlers {
	return &Handlers{
		deps: deps,
	}
}
