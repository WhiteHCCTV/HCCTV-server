// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "log"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		// 모든 클라이언트에게 보낼 데이터 저장 공간
		broadcast:  make(chan []byte),
		// 새로 참여하는 클라이언트
		register:   make(chan *Client),
		// 종료하는 클라이언트
		unregister: make(chan *Client),
		// 현재 허브에 참여중인 클라이언트들
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// log.Println(client.conn)
			h.clients[client] = true
			
			// 등록하게되면 허브의 clients[client]를 true로 등록한다.
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				// 해당 클라이언트 퇴장
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Println(message)

			for client := range h.clients {
				client.send <- message
				// select {
				// case client.send <- message:
				// default:
				// 	close(client.send)
				// 	delete(h.clients, client)
				// }
			}
		}
	}
}