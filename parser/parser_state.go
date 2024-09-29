package parser

// state is a type alias for int and used to create the available value states below
type state int

// Available states for each type used in parsing
const (
	// Group states
	GrpStart state = iota
	GrpOpen

	// Condition states
	ConType
	ConValue
	ConEnd
)
