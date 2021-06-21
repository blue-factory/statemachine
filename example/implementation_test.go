package statemachine

import (
	"fmt"
	"testing"

	"github.com/blue-factory/statemachine"
	"github.com/stretchr/testify/require"
)

func Test_StateMachine_Implementation(t *testing.T) {
	maxCycles := 10
	imp := &implementation{maxCycles: maxCycles}

	sm := statemachine.New(
		&statemachine.Event{Name: eventOne},
		map[string]statemachine.State{
			eventOne: {
				EventHandler: imp.eventOneHandler,
				Destination:  []string{eventTwo, statemachine.EventAbort},
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

	sm := statemachine.New(
		&statemachine.Event{Name: eventOne},
		map[string]statemachine.State{
			eventOne: statemachine.State{
				EventHandler: imp.eventOneHandler,
				Destination:  []string{eventTwo, statemachine.EventAbort},
			},
			eventTwo: statemachine.State{
				EventHandler: func(e *statemachine.Event) (*statemachine.Event, error) {
					return &statemachine.Event{Name: eventOne}, nil
				},
				Destination: []string{eventThree},
			},
			eventThree: statemachine.State{
				EventHandler: imp.eventThreeHandler,
				Destination:  []string{eventOne},
			},
		},
		nil,
	)

	sm.Run()

	require.Equal(t, sm.Error.Error(), `Error: cannot go to next state, wrong destination. From state "two" to state "one"`)
}
