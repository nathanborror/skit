// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	connections map[*connection]bool // Registered connections
	incoming    chan wsRequest       // Inbound messages from the connections
	register    chan *connection     // Register requests from the connections
	unregister  chan *connection     // Unregister requests from connections
}

var h = hub{
	incoming:    make(chan wsRequest),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

/*
type Op int

const (
	Request Op = iota
	Broadcast
)
*/

type Message struct {
	Url string
	//Op Op
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
			fmt.Println(obj.Url)

			req, err := http.NewRequest("GET", "http://localhost:8080"+obj.Url, nil)
			req.Header.Add("X-Requested-With", "XMLHttpRequest")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}

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
