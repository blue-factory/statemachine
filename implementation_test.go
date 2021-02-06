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
		map[string]EventHandler{
			eventOne:   imp.eventOneHandler,
			eventTwo:   imp.eventTwoHandler,
			eventThree: imp.eventThreeHandler,
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
