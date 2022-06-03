// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package manage

import "log"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {

	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		// 모든 클라이언트에게 보낼 데이터 저장 공간
		Broadcast:  make(chan []byte),
		// 새로 참여하는 클라이언트
		Register:   make(chan *Client),
		// 종료하는 클라이언트
		Unregister: make(chan *Client),
		// 현재 허브에 참여중인 클라이언트들
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// log.Println(client.conn)
			h.Clients[client] = true
			log.Println(len(h.Clients),"명의 클라이언트 연결 상태..")
			for client:= range h.Clients {
				log.Println(client.Conn)
			}
			// 등록하게되면 허브의 clients[client]를 true로 등록한다.
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				// 해당 클라이언트 퇴장
				delete(h.Clients, client)
				close(client.Weight)
				log.Println(client.Conn, "클라이언트 연결 종료")
			}
		case message := <-h.Broadcast:
			log.Println(message)
			
			// 모든 클라이언트에게 전송
			for client := range h.Clients {
				(*client.Conn).Write(message)
			}
			
		}
	}
}
