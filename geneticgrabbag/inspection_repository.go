package inspect

// InspectionRepository provides access to health inspection records.
type InspectionRepository interface {

	// Name for this repository instance.
	Name() string

	// GetAll returns a list of all health inspections.
	GetAll() (Inspections, error)
}
