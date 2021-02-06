package statemachine

import "log"

type EventHandler func(e *Event) (*Event, error)

type State struct {
	EventHandler EventHandler
	Destination  []string
}

var (
	EventAbort = "abort"
)

type StateMachine struct {
	initialEvent *Event
	eventChann   chan *Event
	states       map[string]State
}

func New(initialEvent *Event, states map[string]State) *StateMachine {
	return &StateMachine{
		initialEvent: initialEvent,
		eventChann:   make(chan *Event),
		states:       states,
	}
}

func (sm *StateMachine) Run() {
	go sm.SendEvent(sm.initialEvent)
	sm.eventLoop()
}

func (sm *StateMachine) SendEvent(e *Event) {
	sm.eventChann <- e
}

func (sm *StateMachine) defaultErrorHandler(e *Event, err error) {
	log.Printf("Error\nevent: %s\nerror: %s", e.Name, err.Error())
}
