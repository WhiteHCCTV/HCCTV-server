// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"HCCTV/conf"
	. "HCCTV/manage"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
)


var (
	// round scheuler
	//   @ If currWeight > N ( hyperParameter, num of client )
	//        and currWeight == len or hub.clients
	//        -> federated averaging start
	currWeight uint32 = 0

	// HyperParameter, threshold of minimum num of clients
	N = 1
	filenum = 0

	Round = 1
)


func fed_avg(hub *Hub){
	log.Println("\nRound",Round,"종료")
	log.Println("Federated averaging & Measure accuracy\n")
	// @Todo : matrix average algorithm
	cmd  := exec.Command("python3","./aggregator/aggregator.py")
	cmd.Stdout = os.Stdout
	if err := cmd.Run() ; err != nil {
		log.Println(err)
	}
	// round parameter reinitialize
	currWeight = 0
	// fmt.Println("currWeight 초기화 : ", currWeight)
}
func aggregationTimer(hub *Hub, c chan bool){
	// This goroutine is always running state
	for {
		// fmt.Println(currWeight)
		// This will check the state is ready to fed_avg 
		// if (uint32(len(hub.Clients)) == currWeight && len(hub.Clients) > N){
		if (currWeight > 1 ){
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
	weights := make([]byte, 0) // receive local weights
	go func() {
		buf := make([]byte, 30) // this will check disconnection:EOF & "sizeNNNNNNNNNNNN"
		for {
			_, err := conn.Read(buf)
			if err != nil {
				notify <- err
				return
			}

			data := string(buf[:bytes.Index(buf, []byte("\x00"))])
			// Parameter file transfer process
			//   1. Client send me "sizeNN..NN" and binary of file.
			//   2. Parsing size and recevie byte until over size.
			//   3. Cut the weight [0:size] and write file.
			//   4. CurrWeight ++
			if (len(data)>4  && strings.Compare("size",data[:4])==0){
				size, err := strconv.Atoi(data[4:])
				if err != nil {
					log.Println(err)
				}
				fmt.Println("--------Local weights receiving start------")
				dump := 0
				for{
					buf = make([]byte,1024)
					if (size < dump){
						break
					}
					_,err = conn.Read(buf)
					if err!=nil{
						log.Println(err)
					}
					weights = append(weights, buf...)
					dump += 1024
				}
				weights = weights[:size]
				fmt.Println("--------Local weights receiving done-------")

				// now := time.Now().Format("2006-01-02#15:04:05")
				err = ioutil.WriteFile("./aggregator/locals/Client"+strconv.Itoa(filenum), weights, 0644)
				if err != nil {
					panic(err)
				}
				// fmt.Printf("--------Create file at ./weights/%s-------\n",now)
				currWeight++
			}
			
		}
	}()
	client := &Client{Hub: hub, Conn: &conn, Weight: make(chan []byte, 4096)}
	client.Hub.Register <- client	
	// connection 생애 주기 동안 반복
	for {
		select {
		// 연결 해제 감지
		case err := <-notify:
			if io.EOF == err {
				// fmt.Println(client.Conn," is disconnected : ", err)
				currWeight--
				hub.Unregister <- client
				return
			}
		}
	}
}


func main() {
	flag.Parse()
	hub := NewHub()
	c := make(chan bool)
	go aggregationTimer(hub, c)
	go hub.Run()
	// 소켓 서버 
	serverAddr := conf.GetAddr()
	server, err := net.Listen("tcp", serverAddr)
	if err != nil {
		panic(err)
	}
	log.Println("Aggregation server opend at", serverAddr)
	for {
		// 소켓 접속 클라이언트가 생기면 handleConnection goroutine에 위임
		conn, _ := server.Accept()
		go handleConnection(conn, hub)		
	}	
}
