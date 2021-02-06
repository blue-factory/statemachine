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

		handler, ok := sm.eventHandlers[event.Name]
		if !ok {
			log.Printf("Error: unregiestered event %s", event.Name)
			log.Println("event loop stoped")
			return
		}

		go sm.handleFunc(handler, event)
	}
}

func (sm *StateMachine) handleFunc(fn EventHandler, e *Event) {
	// before
	event, err := fn(e)
	if err != nil {
		sm.defaultErrorHandler(e, err)
		sm.SendEvent(&Event{Name: EventAbort})
		return
	}
	sm.SendEvent(event)
	// after
}
