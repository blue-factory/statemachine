package statemachine

import (
	"testing"

	"github.com/blue-factory/cryptobot/internal/statemachine"
	"github.com/stretchr/testify/require"
)

func Test_StateMachine_Implementation(t *testing.T) {
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
				EventHandler: imp.eventTwoHandler,
				Destination:  []string{eventThree},
			},
			eventThree: statemachine.State{
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