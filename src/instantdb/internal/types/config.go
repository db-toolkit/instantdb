package types

// Config holds configuration for starting a database instance
type Config struct {
	Name      string
	Port      int
	DataDir   string
	Persist   bool
	WithData  string
}
