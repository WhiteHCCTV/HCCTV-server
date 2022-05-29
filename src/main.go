// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)
func broadcast(hub *Hub){
	message := []byte("hi im server")
	for client := range hub.clients {
		client.send <- message
		// net.Conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // 2分钟无接收信息超时
		// _,err := (net.Conn).Write(message)
		_, err := (*client.conn).Write(message)
		if err != nil {
			fmt.Println("Failed to write data : ", err)
			break;
		}
		// time.Sleep(1 * time.Second)
		log.Println("Send to cli")
	}
	
	
}
func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()

	// 소켓 서버 
	serverAddr := "localhost:4000"
	server, err := net.Listen("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	for {
        conn, err := server.Accept()
        if err != nil {
            log.Println("Failed to Accept : ", err)
            continue
        }
		log.Println(conn," is conneted !!")
		// 클라이언트가 서버에 접속하면 client객체 생성하여 허브에 추가
		client := &Client{hub: hub, conn: &conn, send: make(chan []byte, 256)}
		client.hub.register <- client
		// go client.writePump()
        go broadcast(hub)
    }
}
