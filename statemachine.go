package statemachine

import (
	"bytes"
	"fmt"
	"log"

	internallogger "github.com/blue-factory/cryptobot/internal/logger"
)

type EventHandler func(e *Event) (*Event, error)
type OnStateChangeHandler func(state string) error

type State struct {
	EventHandler EventHandler
	Destination  []string
}

var (
	EventAbort    = "abort"
	PristineState = "pristine"
)

type StateMachine struct {
	initialEvent  *Event
	current       string
	previous      string
	eventChann    chan *Event
	states        map[string]State
	Error         error
	logger        Logger
	onStateChange OnStateChangeHandler
}

func noopOnStateChange(state string) error { return nil }

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
		initialEvent:  initialEvent,
		current:       PristineState,
		previous:      PristineState,
		eventChann:    make(chan *Event),
		states:        states,
		logger:        logger,
		onStateChange: noopOnStateChange,
	}
}

func (sm *StateMachine) OnStateChange(fn OnStateChangeHandler) {
	sm.onStateChange = fn
}

func (sm *StateMachine) Run() {
	go sm.Dispatch(sm.initialEvent)
	sm.eventLoop()
}

func (sm *StateMachine) Dispatch(e *Event) {
	sm.eventChann <- e
}

func (sm *StateMachine) RenderGraphviz() string {
	b := bytes.NewBufferString("")
	b.WriteString("digraph {\n")
	b.WriteString("\trankdir=LR;\n")
	b.WriteString("\tsize=\"8\"\n")
	b.WriteString("\tnode [shape = circle];\n")

	for current, s := range sm.states {
		// TODO(ca): Add label value to state struct, eg. [label = "label"]
		for _, dest := range s.Destination {
			b.WriteString(fmt.Sprintf("\t%s -> %s;\n", current, dest))
		}
	}

	b.WriteString("}")

	return b.String()
}

func (sm *StateMachine) RenderMermaid() string {
	b := bytes.NewBufferString("")
	b.WriteString("stateDiagram-v2\n")

	for current, s := range sm.states {
		for _, dest := range s.Destination {
			c := current
			if current == PristineState {
				c = "[*]"
			}

			d := dest
			if dest == EventAbort {
				d = "[*]"
			}

			b.WriteString(fmt.Sprintf("\t%s --> %s\n", c, d))
		}
	}

	return b.String()
}

func (sm *StateMachine) defaultErrorHandler(e *Event, err error) {
	log.Printf("Error\nevent: %s\nerror: %s", e.Name, err.Error())
}
