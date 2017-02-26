package dispatch

import (
	"github.com/gorilla/websocket"
	"main/utils"
	"net/http"
	"encoding/json"
)

type (
	webSocketMessage struct {
		Content []byte
		Type    int
	}
)

var (
	upgrader websocket.Upgrader
)

func init() {
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

// CreateDispatchHandler returns a REST API handler for a given dispatcher
func CreateDispatchHandler(dispatcher *Dispatcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade the HTTP request to a WebSocket
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			utils.WriteError(w, utils.InternalServerError(err))
			return
		}
		defer c.Close()

		// Create a client
		clnt := dispatcher.CreateClient()

		// Start a goroutine listening for incoming messages
		incomingMessages := make(chan webSocketMessage, 10)
		closed := make(chan bool)

		go func() {
			for {
				msgType, content, err := c.ReadMessage()

				if err != nil {
					if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
						closed <- true
						return
					}
					utils.LogError(err)
				}

				incomingMessages <- webSocketMessage{Content: content, Type: msgType}
			}
		}()

		// Listen for incoming and outgoing messages
		for {
			select {
			case message := <-incomingMessages:
				switch message.Type {
				default:
					var incoming incomingMessage
					err = json.Unmarshal(message.Content, &incoming)

					if err != nil {
						clnt.OutgoingMessages <- BadRequestError(err).OutgoingMessage(-1)
						continue
					}

					clnt.handleIncomingMessage(incoming)
				}
			case message := <-clnt.OutgoingMessages:
				err := c.WriteJSON(message)
				if err != nil {
					utils.LogError(err)
					continue
				}
			case _ = <-closed:
				return
			}
		}

		close(incomingMessages)
		close(closed)

		dispatcher.RemoveClient(clnt)
	}
}
