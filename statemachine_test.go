package statemachine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StateMachine_Dispatch(t *testing.T) {
	sm := &StateMachine{eventChann: make(chan *Event)}

	firstEventSent := &Event{Name: "first"}
	go sm.Dispatch(firstEventSent)
	firstEvent := <-sm.eventChann

	secondEventSent := &Event{Name: "second"}
	go sm.Dispatch(secondEventSent)
	secondEvent := <-sm.eventChann

	require.NotEqual(t, firstEvent.Name, secondEvent.Name)
	require.Equal(t, firstEvent, firstEventSent)
	require.Equal(t, secondEvent, secondEventSent)
}

func Test_StateMachine_Implementation(t *testing.T) {
	maxCycles := 10
	imp := &implementation{maxCycles: maxCycles}

	sm := New(
		&Event{Name: eventOne},
		map[string]State{
			eventOne: {
				EventHandler: imp.eventOneHandler,
				Destination:  []string{eventTwo, EventAbort},
			},
			eventTwo: {
				EventHandler: imp.eventTwoHandler,
				Destination:  []string{eventThree},
			},
			eventThree: {
				EventHandler: imp.eventThreeHandler,
				Destination:  []string{eventOne},
			},
		},
		nil,
	)

	sm.Run()
	fmt.Println("----------------------")
	fmt.Println("----------------------")
	fmt.Println(sm.RenderMermaid())
	fmt.Println("----------------------")
	fmt.Println("----------------------")

	require.NoError(t, sm.Error)
	require.Equal(t, imp.cycles, len(imp.eventOneCalls))
	require.Equal(t, imp.cycles, len(imp.eventTwoCalls))
	require.Equal(t, imp.cycles, len(imp.eventThreeCalls))

	for i := 0; i < maxCycles; i++ {
		require.Equal(t, imp.eventOneCalls[i], imp.eventTwoCalls[i])
		require.Equal(t, imp.eventOneCalls[i], imp.eventThreeCalls[i])
		require.Equal(t, imp.eventTwoCalls[i], imp.eventThreeCalls[i])
	}
}

func Test_StateMachine_Implementation_Wrong_Destination(t *testing.T) {
	maxCycles := 10
	imp := &implementation{maxCycles: maxCycles}

	sm := New(
		&Event{Name: eventOne},
		map[string]State{
			eventOne: State{
				EventHandler: imp.eventOneHandler,
				Destination:  []string{eventTwo, EventAbort},
			},
			eventTwo: State{
				EventHandler: func(e *Event) (*Event, error) {
					return &Event{Name: eventOne}, nil
				},
				Destination: []string{eventThree},
			},
			eventThree: State{
				EventHandler: imp.eventThreeHandler,
				Destination:  []string{eventOne},
			},
		},
		nil,
	)

	sm.Run()

	require.Equal(t, sm.Error.Error(), `Error: cannot go to next state, wrong destination. From state "two" to state "one"`)
}

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

func (s *implementation) eventOneHandler(e *Event) (*Event, error) {
	if s.cycles == s.maxCycles {
		return &Event{Name: EventAbort}, nil
	}

	s.eventOneCalls = append(s.eventOneCalls, len(s.eventOneCalls)+1)
	s.cycles = s.cycles + 1
	return &Event{Name: eventTwo}, nil
}

func (s *implementation) eventTwoHandler(e *Event) (*Event, error) {
	s.eventTwoCalls = append(s.eventTwoCalls, len(s.eventTwoCalls)+1)
	return &Event{Name: eventThree}, nil
}

func (s *implementation) eventThreeHandler(e *Event) (*Event, error) {
	s.eventThreeCalls = append(s.eventThreeCalls, len(s.eventThreeCalls)+1)
	return &Event{Name: eventOne}, nil
}
