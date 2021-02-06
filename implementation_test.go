package statemachine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StateMachine_Implementation(t *testing.T) {
	maxCycles := 10
	imp := &implementation{maxCycles: maxCycles}

	sm := New(
		&Event{Name: eventOne},
		map[string]State{
			eventOne: State{
				EventHandler: imp.eventOneHandler,
				Destination:  []string{eventTwo},
			},
			eventTwo: State{
				EventHandler: imp.eventTwoHandler,
				Destination:  []string{eventThree},
			},
			eventThree: State{
				EventHandler: imp.eventThreeHandler,
				Destination:  []string{eventOne},
			},
		},
	)

	sm.Run()

	require.Equal(t, imp.cycles, len(imp.eventOneCalls))
	require.Equal(t, imp.cycles, len(imp.eventTwoCalls))
	require.Equal(t, imp.cycles, len(imp.eventThreeCalls))

	for i := 0; i < maxCycles; i++ {
		require.Equal(t, imp.eventOneCalls[i], imp.eventTwoCalls[i])
		require.Equal(t, imp.eventOneCalls[i], imp.eventThreeCalls[i])
		require.Equal(t, imp.eventTwoCalls[i], imp.eventThreeCalls[i])
	}
}
