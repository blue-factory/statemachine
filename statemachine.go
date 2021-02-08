package statemachine

import (
	"fmt"
	"log"
)

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
	go sm.Dispatch(sm.initialEvent)
	sm.eventLoop()
}

func (sm *StateMachine) Dispatch(e *Event) {
	sm.eventChann <- e
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
