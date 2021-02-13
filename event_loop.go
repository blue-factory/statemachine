package statemachine

import (
	"fmt"
	"strings"
)

func (sm *StateMachine) eventLoop() {
	sm.logger.Info("starting event loop...")
	for {
		nextEvent := <-sm.eventChann
		if nextEvent.Name == EventAbort {
			sm.logger.Info("event loop aborted")
			return
		}

		nextState, ok := sm.states[nextEvent.Name]
		if !ok {
			sm.Error = fmt.Errorf("Error: unregistered event %s", nextEvent.Name)
			sm.logger.Info(sm.Error)
			sm.logger.Info("event loop stoped")
			return
		}

		err := sm.validateTransition(nextEvent)
		if err != nil {
			sm.Error = err
			sm.logger.Info(sm.Error)
			sm.logger.Info("event loop stoped")
			return
		}

		eventToDispatch, err := sm.handleFunc(nextState.EventHandler, nextEvent)
		if err != nil {
			eventToDispatch = &Event{Name: EventAbort}
			sm.defaultErrorHandler(nextEvent, err)
		}

		go sm.Dispatch(eventToDispatch)
	}
}

func (sm *StateMachine) handleFunc(fn EventHandler, event *Event) (*Event, error) {
	nextEvent, err := fn(event)
	if err != nil {
		return nil, err
	}
	sm.previous = sm.current
	sm.current = event.Name

	return nextEvent, nil
}

func (sm *StateMachine) validateTransition(event *Event) error {
	currentState, ok := sm.states[sm.current]
	if !ok {
		return fmt.Errorf("Error: current state does not exists. %q", sm.current)
	}

	if sm.current != PristineState && !strings.Contains(strings.Join(currentState.Destination, " $ "), event.Name) {
		return fmt.Errorf("Error: cannot go to next state, wrong destination. From state %q to state %q", sm.current, event.Name)
	}

	return nil
}
