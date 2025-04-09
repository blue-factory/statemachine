package main

import (
	"log"
	"net/http"
	"time"

	"github.com/blue-factory/statemachine"
	"github.com/blue-factory/statemachine/websocket"
)

func main() {
	// Creamos una máquina de estados simple
	sm := statemachine.New(
		&statemachine.Event{Name: "start"},
		map[string]statemachine.State{
			"start": {
				EventHandler: func(e *statemachine.Event) (*statemachine.Event, error) {
					// Esperamos 2 segundos antes de la transición
					time.Sleep(2 * time.Second)
					return &statemachine.Event{Name: "state-mid"}, nil
				},
				Destination: []string{"state-mid"},
			},
			"state-mid": {
				EventHandler: func(e *statemachine.Event) (*statemachine.Event, error) {
					// Esperamos 2 segundos antes de la transición
					time.Sleep(2 * time.Second)
					return &statemachine.Event{Name: "state-end"}, nil
				},
				Destination: []string{"state-end"},
			},
			"state-end": {
				EventHandler: func(e *statemachine.Event) (*statemachine.Event, error) {
					// Esperamos 2 segundos antes de la transición
					time.Sleep(2 * time.Second)
					return &statemachine.Event{Name: "start"}, nil
				},
				Destination: []string{"start"},
			},
		},
		nil,
	)

	// Creamos el WebSocket Manager
	wsManager := websocket.NewWebSocketManager(sm)

	// Configuramos el handler para el WebSocket
	http.HandleFunc("/ws", wsManager.HandleWebSocket)

	// Conectamos el WebSocket Manager con los cambios de estado
	sm.OnStateChange(func(state string) error {
		wsManager.BroadcastStateChange(sm.GetPreviousState(), state)
		return nil
	})

	// Iniciamos el WebSocket Manager en una goroutine
	go wsManager.Start()

	// Iniciamos la máquina de estados en una goroutine
	go sm.Run()

	// Iniciamos el servidor HTTP
	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
} 