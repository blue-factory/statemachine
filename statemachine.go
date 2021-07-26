package statemachine

import (
	log "github.com/sirupsen/logrus"
)

// EventHandler is the event handler type definition.
type EventHandler func(e *Event) (*Event, error)

// OnStateChangeHandler is the handler fun executed every time the state changes.
type OnStateChangeHandler func(state string) error

// State represents the combination according an EventHandler and
// its valid destinations. if the EventHandler dispatches an event
// which is not defined into Destination slice, this will break
// the event loop due to an invalid state transition.
type State struct {
	EventHandler EventHandler
	Destination  []string
}

var (
	// EventAbort is the abort state-name used to break the event loop
	EventAbort = "abort"
	// PristineState is the pristine state used to start the statemachine
	PristineState = "pristine"
)

// StateMachine is the functional struct which runs the event loop, trigger the state
// transition and also executes the correct EventHandler for every event dispatched
type StateMachine struct {
	// Unexported properties
	initialEvent  *Event
	current       string
	previous      string
	eventChann    chan *Event
	states        map[string]State
	Error         error
	logger        Logger
	onStateChange OnStateChangeHandler
}

// New constructs a new statemachine. This function needs you to specify the initialEvent
// the complete map[string-state-name]State
func New(initialEvent *Event, states map[string]State, logger Logger) *StateMachine {
	if logger == nil {
		logger = log.New()
	}

	states[PristineState] = State{
		EventHandler: func(e *Event) (*Event, error) {
			return initialEvent, nil
		},
		Destination: []string{initialEvent.Name},
	}
	states[EventAbort] = State{
		EventHandler: func(e *Event) (*Event, error) {
			return initialEvent, nil
		},
		Destination: []string{initialEvent.Name},
	}

	return &StateMachine{
		initialEvent:  initialEvent,
		current:       PristineState,
		previous:      PristineState,
		states:        states,
		logger:        logger,
		onStateChange: noopOnStateChange,
	}
}

func noopOnStateChange(state string) error { return nil }

// OnStateChange will replace the noopOnStateChange function by the given func.
func (sm *StateMachine) OnStateChange(fn OnStateChangeHandler) {
	sm.onStateChange = fn
}

// Run starts the event loop and dispatches the initial event
func (sm *StateMachine) Run() {
	sm.eventChann = make(chan *Event)
	go sm.Dispatch(sm.initialEvent)
	sm.eventLoop()
}

// Stop will immediately dispatch an EventAbort. Is a syntax sugar to cal Dispatch
// with EventAbort as the event name
func (sm *StateMachine) Stop() chan error {
	return sm.Dispatch(&Event{Name: EventAbort})
}

// Dispatch dispatches the given event and returns a chan error with cap=1
func (sm *StateMachine) Dispatch(e *Event) chan error {
	e.done = make(chan error, 1)
	sm.eventChann <- e

	return e.done
}

func (sm *StateMachine) defaultErrorHandler(e *Event, err error) {
	sm.logger.Infof("Error\nevent: %s\nerror: %s", e.Name, err.Error())
}
