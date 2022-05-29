// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
)

// Main goroutine에서 테스트용 함수
func checkClient(hub *Hub){
	log.Println(len(hub.clients),"명의 클라이언트 연결 상태..")
	for client:= range hub.clients {
		log.Println(client.conn)
	}
}

// 각각의 connection handle
//   @ 커넥션 종료를 감지하면 notify 채널에 에러 전달
//   @ hub의 unregister 채널에 종료한 클라이언트를 전달하고 핸들러 고루틴 종료
func handleConnection(conn net.Conn, hub *Hub) {
	defer conn.Close()
	notify := make(chan error) // detect connection state
	weights := make(chan []byte) // receive weights
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				notify <- err
				return
			}
			if n > 0 {
				fmt.Println("unexpected data: %s", buf[:n])
			}
			weights <- buf[:n]
		}
	}()
	client := &Client{hub: hub, conn: &conn, send: make(chan []byte, 256)}
	client.hub.register <- client
	// 새로운 클라이언트 연결 시 연결된 모든 클라이언트에게 테스트 브로드캐스트
	hub.broadcast<-([]byte("테스트"))
	
	// connection 생애 주기 동안 반복
	for {
		select {
		// 연결 해제 감지
		case err := <-notify:
			if io.EOF == err {
				fmt.Println("connection dropped message", err)
				hub.unregister <- client
				return
			}
		// 로컬 모델의 가중치 수신 감지
		case receive := <- weights:
			fmt.Println("receive",string(receive))
		}
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
		// 소켓 접속 클라이언트가 생기면 handleConnection goroutine에 위임
		conn, _ := server.Accept()
		go handleConnection(conn, hub)
	}	
}
