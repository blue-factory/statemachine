package statemachine

type Event struct {
	Name string
	Data map[string]interface{}
	done chan error
}
