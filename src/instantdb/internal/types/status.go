package types

// Status represents the current status of an instance
type Status struct {
	Running bool
	Healthy bool
	Message string
}
