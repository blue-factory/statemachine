package statemachine

import (
	"fmt"
	"log"

	internallogger "github.com/blue-factory/cryptobot/internal/logger"
)

type EventHandler func(e *Event) (*Event, error)

type State struct {
	EventHandler EventHandler
	Destination  []string
}

var (
	EventAbort    = "abort"
	PristineState = "pristine"
)

type StateMachine struct {
	initialEvent *Event
	current      string
	previous     string
	eventChann   chan *Event
	states       map[string]State
	Error        error
	logger       Logger
}

func New(initialEvent *Event, states map[string]State, logger Logger) *StateMachine {
	if logger == nil {
		logger = internallogger.New()
	}

	states[PristineState] = State{
		EventHandler: func(e *Event) (*Event, error) {
			return initialEvent, nil
		},
		Destination: []string{initialEvent.Name},
	}
	return &StateMachine{
		initialEvent: initialEvent,
		current:      PristineState,
		previous:     PristineState,
		eventChann:   make(chan *Event),
		states:       states,
		logger:       logger,
	}
}

func (sm *StateMachine) Run() {
	go sm.Dispatch(sm.initialEvent)
	sm.eventLoop()
}

func (sm *StateMachine) Dispatch(e *Event) chan error {
	e.done = make(chan error, 1)
	sm.eventChann <- e

	return e.done
}

func (sm *StateMachine) Render() string {
	str := "digraph {\n"
	str += "\trankdir=LR;\n"
	str += "\tsize=\"8\"\n"
	str += "\tnode [shape = circle];\n"

	for current, s := range sm.states {
		// TODO(ca): Add label value to state struct, eg. [label = "label"]
		for _, dest := range s.Destination {
			str += fmt.Sprintf("\t%s -> %s;\n", current, dest)
		}
	}

	str += "}"

	return str
}

func (sm *StateMachine) defaultErrorHandler(e *Event, err error) {
	log.Printf("Error\nevent: %s\nerror: %s", e.Name, err.Error())
}
