package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/blue-factory/statemachine"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // En producción, implementar validación de origen
	},
}

// WebSocketManager maneja todas las conexiones WebSocket
type WebSocketManager struct {
	clients          map[*Client]bool
	broadcast        chan []byte
	register         chan *Client
	unregister       chan *Client
	stateMachine     *statemachine.StateMachine
	stateMachineInfo []byte
	mu               sync.RWMutex
}

// NewWebSocketManager crea una nueva instancia del WebSocket Manager
func NewWebSocketManager(sm *statemachine.StateMachine) *WebSocketManager {
	manager := &WebSocketManager{
		clients:      make(map[*Client]bool),
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		stateMachine: sm,
	}

	// Pre-generamos el JSON con la información de la máquina de estados
	manager.stateMachineInfo = manager.generateStateMachineInfo()

	return manager
}

// Start inicia el WebSocket Manager
func (m *WebSocketManager) Start() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true
			m.mu.Unlock()
			// Enviamos la información de la máquina de estados al nuevo cliente
			client.send <- m.stateMachineInfo

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client.send)
			}
			m.mu.Unlock()

		case message := <-m.broadcast:
			m.mu.RLock()
			for client := range m.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(m.clients, client)
				}
			}
			m.mu.RUnlock()
		}
	}
}

// HandleWebSocket maneja las conexiones WebSocket entrantes
func (m *WebSocketManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error al actualizar conexión: %v", err)
		return
	}

	client := &Client{
		manager: m,
		conn:    conn,
		send:    make(chan []byte, 256),
	}

	m.register <- client

	go client.writePump()
}

// BroadcastStateChange envía un mensaje de cambio de estado a todos los clientes
func (m *WebSocketManager) BroadcastStateChange(prevState, currentState string) {
	message := map[string]interface{}{
		"type": "state_change",
		"data": map[string]interface{}{
			"previous_state": prevState,
			"current_state":  currentState,
		},
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error al marshal mensaje de estado: %v", err)
		return
	}

	m.broadcast <- jsonMessage
}

// generateStateMachineInfo genera el JSON con la información de la máquina de estados
func (m *WebSocketManager) generateStateMachineInfo() []byte {
	info := map[string]interface{}{
		"type": "state_machine_info",
		"data": map[string]interface{}{
			"states":      make([]string, 0),
			"transitions": make([]map[string]string, 0),
		},
	}

	states := m.stateMachine.GetStates()
	for stateName, state := range states {
		info["data"].(map[string]interface{})["states"] = append(
			info["data"].(map[string]interface{})["states"].([]string),
			stateName,
		)

		for _, dest := range state.Destination {
			transition := map[string]string{
				"from": stateName,
				"to":   dest,
			}
			info["data"].(map[string]interface{})["transitions"] = append(
				info["data"].(map[string]interface{})["transitions"].([]map[string]string),
				transition,
			)
		}
	}

	jsonData, err := json.Marshal(info)
	if err != nil {
		log.Printf("Error al generar info de máquina de estados: %v", err)
		return []byte("{}")
	}

	return jsonData
}
