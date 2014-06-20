// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	connections map[*connection]bool // Registered connections
	incoming    chan wsRequest       // Inbound messages from the connections
	register    chan *connection     // Register requests from the connections
	unregister  chan *connection     // Unregister requests from connections
}

// h represents a WebSocket hub
var h = hub{
	incoming:    make(chan wsRequest),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

// Op represents 'broadcast' or 'request'
type Op int

const (
	// Request means we only update the client making the request
	Request Op = iota
	// Broadcast means we update all connected clients
	Broadcast
)

// Message expects a url and an operation (request or broadcast)
type Message struct {
	URL string
	Op  Op
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		case r := <-h.incoming:
			var obj Message
			json.Unmarshal(r.message, &obj)

			req, err := http.NewRequest("GET", "http://localhost:8080"+obj.URL, nil)

			// Update cursor
			r.connection.Cursor = obj.URL

			// Add cookie
			req.AddCookie(r.connection.Cookie)

			// TODO: Replace this with "the better way" referenced in
			// render.go and update here.
			req.Header.Add("X-Requested-With", "XMLHttpRequest")

			// Make request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println(err)
			}

			// Get the body of the response
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}

			if obj.Op == Request  {
				// Only send to clients with cursors on the request URL
				for c := range h.connections {
					// if r.connection.User == c.User {
					if c.Cursor == obj.URL {
						if obj.URL != "/" {
							select {
							case c.send <- body:
							default:
								close(c.send)
								delete(h.connections, c)
							}
						} else {
							if c.User == r.connection.User {
								select {
								case c.send <- body:
								default:
									close(c.send)
									delete(h.connections, c)
								}
							}
						}
					}
				}
			} else {
				// Broadcast to all connected clients
				for c := range h.connections {
					select {
					case c.send <- body:
					default:
						close(c.send)
						delete(h.connections, c)
					}
				}
			}
		}
	}
}
