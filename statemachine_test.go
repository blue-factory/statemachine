package statemachine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StateMachine_SendEvent(t *testing.T) {
	sm := &StateMachine{eventChann: make(chan *Event)}

	firstEventSent := &Event{Name: "first"}
	go sm.SendEvent(firstEventSent)
	firstEvent := <-sm.eventChann

	secondEventSent := &Event{Name: "second"}
	go sm.SendEvent(secondEventSent)
	secondEvent := <-sm.eventChann

	require.NotEqual(t, firstEvent.Name, secondEvent.Name)
	require.Equal(t, firstEvent, firstEventSent)
	require.Equal(t, secondEvent, secondEventSent)
}
