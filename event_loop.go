package statemachine

import "log"

func (sm *StateMachine) eventLoop() {
	log.Println("starting event loop...")
	for {
		event := <-sm.eventChann
		if event.Name == EventAbort {
			log.Println("event loop aborted")
			return
		}

		state, ok := sm.states[event.Name]
		if !ok {
			log.Printf("Error: unregistered event %s", event.Name)
			log.Println("event loop stoped")
			return
		}

		go sm.handleFunc(state.EventHandler, event)
	}
}

func (sm *StateMachine) handleFunc(fn EventHandler, currentEvent *Event) {
	nextEvent, err := fn(currentEvent)
	if err != nil {
		sm.defaultErrorHandler(currentEvent, err)
		sm.Dispatch(&Event{Name: EventAbort})
		return
	}
	sm.Dispatch(nextEvent)
}
