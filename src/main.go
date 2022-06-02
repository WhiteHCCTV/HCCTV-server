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
var now uint32 = 0
// Main goroutine에서 테스트용 함수
func checkClient(hub *Hub){
	log.Println(len(hub.clients),"명의 클라이언트 연결 상태..")
	for client:= range hub.clients {
		log.Println(client.conn)
	}
}
func fed_avg(hub *Hub){
	for {}
}
func aggregationTimer(hub *Hub, c chan bool){
	for {
		if (hub.count == now){
			fed_avg(hub)
			now = 0
		}
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
			weights <- buf[:n]
		}
	}()
	client := &Client{hub: hub, conn: &conn, send: make(chan []byte, 4096)}
	client.hub.register <- client
	hub.count++
	// 새로운 클라이언트 연결 시 연결된 모든 클라이언트에게 테스트 브로드캐스트
	hub.broadcast<-([]byte("테스트"))
	fmt.Println("![Client add]",hub.count, " and " , now , "!amd",len(hub.clients))
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
			// fmt.Println(receive)
			// fmt.Println(string(receive))
			fmt.Println("received : ",len(receive))
			now++
			fmt.Println("![Client send to me]",hub.count, " and " , now , "!amd",len(hub.clients))

		}
	}
}


func main() {
	flag.Parse()
	hub := newHub()
	c := make(chan bool)
	go aggregationTimer(hub, c)
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
		fmt.Println(conn.RemoteAddr().String())
		go handleConnection(conn, hub)
		
	}	
}
