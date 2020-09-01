package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// WebsocketStream handles the websocket connection
type WebsocketStream struct {
	Ws        *websocket.Conn
	WsChannel chan []byte
}

// Close closes the websocket connection
func (w *WebsocketStream) Close() error {
	err := w.Ws.Close()
	//close(w.WsChannel)

	if err != nil {
		return err
	}
	return nil
}

//
func (w *WebsocketStream) Handler(ch <-chan []byte) {
	defer func() {
		fmt.Println("closing the websocket connection")
		err := w.Close()
		if err != nil {
			_ = fmt.Errorf("error closing the websocket: %v", err)
			return
		}
	}()

	for {
		select {
		case data := <-ch:
			w.Write(data)
		}
	}
}

// Write writes data to the websocket connection
func (w *WebsocketStream) Write(data []byte) {
	if err := w.Ws.WriteMessage(websocket.TextMessage, data); err != nil {
		_ = fmt.Errorf("error writing the message to the websocket: %v", err)
		return
	}
}
