// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"HCCTV/conf"
	"flag"
	"fmt"
	"io"
	"net"
)


var (
	// round scheuler
	//   @ If currWeight > N ( hyperParameter, num of client )
	//        and currWeight == len or hub.clients
	//        -> federated averaging start
	currWeight uint32 = 0

	// HyperParameter, threshold of minimum num of clients
	N = 5
)


func fed_avg(hub *Hub){
	fmt.Println("가중치 평균 연산 시작")
	fmt.Println("현재 참여 중인 클라이언트")

	// @Todo : matrix average algorithm
	for client := range hub.clients{
		fmt.Println(client.conn)
	}

	// round parameter reinitialize
	currWeight = 0
	fmt.Println("currWeight 초기화 : ", currWeight)
}
func aggregationTimer(hub *Hub, c chan bool){
	// This goroutine is always running state
	for {
		// This will check the state is ready to fed_avg 
		if (uint32(len(hub.clients)) == currWeight && len(hub.clients) > N){
			fed_avg(hub)
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
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				notify <- err
				return
			}
			weights <- buf[:n]
		}
	}()
	client := &Client{hub: hub, conn: &conn, weight: make(chan []byte, 4096)}
	client.hub.register <- client	
	// connection 생애 주기 동안 반복
	for {
		select {
		// 연결 해제 감지
		case err := <-notify:
			if io.EOF == err {
				fmt.Println(client.conn," is disconnected : ", err)
				currWeight--
				hub.unregister <- client
				return
			}
		// 로컬 모델의 가중치 수신 감지
		case receive := <- weights:
			fmt.Println("received : ",receive)
			currWeight++
			fmt.Println("![Client send to me] recevied : ", currWeight , " and connected : ",len(hub.clients))
		}
	}
}


func main() {
	fmt.Println(conf.GetAddr())
	flag.Parse()
	hub := newHub()
	c := make(chan bool)
	go aggregationTimer(hub, c)
	go hub.run()
	// 소켓 서버 
	serverAddr := "localhost:8080"
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
