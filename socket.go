// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/gorilla/websocket"
	"github.com/nathanborror/skit/users"
	"log"
	"net/http"
	"time"
)

var userRepo = users.NewSqlUserRepository("db.sqlite3")

const (
	writeWait      = 10 * time.Second    // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second    // Time allowed to read the next pong message from the peer.
	pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 512                 // Maximum message size allowed from peer.
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	ws     *websocket.Conn // The websocket connection.
	send   chan []byte     // Buffered channel of outbound messages.
	User   *users.User
	Cookie *http.Cookie
	Cursor string
}

// wsRequest
type wsRequest struct {
	connection *connection
	message    []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		h.incoming <- wsRequest{c, message}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serverWs handles webocket requests from the peer.
func socketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	// Upgrade request to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	// Check for authenticated user
	session, _ := store.Get(r, "authenticated-user")
	if session.Values["hash"] == nil {
		return
	}

	// Grab user making request
	userHash := session.Values["hash"].(string)
	u, err := userRepo.Load(userHash)
	if err != nil {
		log.Println(err)
	}

	// Grab cookie
	cookie, _ := r.Cookie("authenticated-user")

	// Create connection
	c := &connection{send: make(chan []byte, 256), ws: ws, User: u, Cookie: cookie}
	h.register <- c
	go c.writePump()
	c.readPump()
}
