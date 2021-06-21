package main

import (
	"fmt"

	"github.com/blue-factory/statemachine"
)

var (
	eventOne   = "one"
	eventTwo   = "two"
	eventThree = "three"
)

func main() {
	basicSM := &BasicSM{maxCycles: 15}

	sm := statemachine.New(
		&statemachine.Event{Name: eventOne},
		map[string]statemachine.State{
			eventOne: {
				EventHandler: basicSM.eventOneHandler,
				Destination:  []string{eventTwo, statemachine.EventAbort},
			},
			eventTwo: {
				EventHandler: basicSM.eventTwoHandler,
				Destination:  []string{eventThree},
			},
			eventThree: {
				EventHandler: basicSM.eventThreeHandler,
				Destination:  []string{eventOne},
			},
		},
		nil,
	)

	fmt.Println(sm.RenderMermaid())
	fmt.Println("----------------------")

	sm.Run()

	fmt.Printf("Expected max cicles %d\t Actual max cicles: %d\n", basicSM.maxCycles, basicSM.cycles)
	fmt.Printf("call stack:\nOne: %v\nTwo: %v\nThree: %v\n", basicSM.eventOneCalls, basicSM.eventTwoCalls, basicSM.eventThreeCalls)
}

type BasicSM struct {
	maxCycles       int
	cycles          int
	eventOneCalls   []int
	eventTwoCalls   []int
	eventThreeCalls []int
}

func (s *BasicSM) eventOneHandler(e *statemachine.Event) (*statemachine.Event, error) {
	if s.cycles == s.maxCycles {
		return &statemachine.Event{Name: statemachine.EventAbort}, nil
	}

	s.eventOneCalls = append(s.eventOneCalls, len(s.eventOneCalls)+1)
	s.cycles = s.cycles + 1
	return &statemachine.Event{Name: eventTwo}, nil
}

func (s *BasicSM) eventTwoHandler(e *statemachine.Event) (*statemachine.Event, error) {
	s.eventTwoCalls = append(s.eventTwoCalls, len(s.eventTwoCalls)+1)
	return &statemachine.Event{Name: eventThree}, nil
}

func (s *BasicSM) eventThreeHandler(e *statemachine.Event) (*statemachine.Event, error) {
	s.eventThreeCalls = append(s.eventThreeCalls, len(s.eventThreeCalls)+1)
	return &statemachine.Event{Name: eventOne}, nil
}
