package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var AllowedOrigins []string

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		for _, allowedOrigin := range AllowedOrigins {
			if allowedOrigin == origin {
				return true
			}
		}
		return false
	},
}

// Stream handles the websocket connection
type Stream struct {
	WsConn    *websocket.Conn
	WsChannel chan []byte
}

// NewStream creates a new WebsocketStream
func NewStream(wsConn *websocket.Conn, wsChan chan []byte) *Stream {
	return &Stream{
		WsConn:    wsConn,
		WsChannel: wsChan,
	}
}

// Handler handles the data received on the channel
func (w *Stream) Handler(ch <-chan []byte) {
	defer func() {
		fmt.Println("closing the websocket connection")
		err := w.close()
		if err != nil {
			_ = fmt.Errorf("error closing the websocket: %v", err)
			return
		}
	}()

	for {
		select {
		case data := <-ch:
			w.write(data)
		}
	}
}

// Close closes the websocket connection
func (w *Stream) close() error {
	err := w.WsConn.Close()
	//close(w.WsChannel)

	if err != nil {
		return err
	}
	return nil
}

// Write writes data to the websocket connection
func (w *Stream) write(data []byte) {
	if err := w.WsConn.WriteMessage(websocket.TextMessage, data); err != nil {
		_ = fmt.Errorf("error writing the message to the websocket: %v", err)
		return
	}
}
