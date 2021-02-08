package statemachine

import (
	"github.com/blue-factory/cryptobot/internal/statemachine"
)

var (
	eventOne   = "one"
	eventTwo   = "two"
	eventThree = "three"
)

type implementation struct {
	maxCycles       int
	cycles          int
	eventOneCalls   []int
	eventTwoCalls   []int
	eventThreeCalls []int
}

func (s *implementation) eventOneHandler(e *statemachine.Event) (*statemachine.Event, error) {
	if s.cycles == s.maxCycles {
		return &statemachine.Event{Name: statemachine.EventAbort}, nil
	}

	s.eventOneCalls = append(s.eventOneCalls, len(s.eventOneCalls)+1)
	s.cycles = s.cycles + 1
	return &statemachine.Event{Name: eventTwo}, nil
}

func (s *implementation) eventTwoHandler(e *statemachine.Event) (*statemachine.Event, error) {
	s.eventTwoCalls = append(s.eventTwoCalls, len(s.eventTwoCalls)+1)
	return &statemachine.Event{Name: eventThree}, nil
}

func (s *implementation) eventThreeHandler(e *statemachine.Event) (*statemachine.Event, error) {
	s.eventThreeCalls = append(s.eventThreeCalls, len(s.eventThreeCalls)+1)
	return &statemachine.Event{Name: eventOne}, nil
}
