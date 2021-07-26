package statemachine

// Event is the basic functional unit of the statemachine.
// Name should never be an empty string
type Event struct {
	Name string
	Data map[string]interface{}
	done chan error
}
