package types

// Instance represents a running database instance
type Instance struct {
	ID        string
	Name      string
	Engine    string
	Port      int
	DataDir   string
	PID       int
	Status    string
	CreatedAt int64
	Persist   bool
	Username  string
	Password  string
}
