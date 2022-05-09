package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// var (
// 	upgrader = websocket.Upgrader{}
// )

func test(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Read
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
      var objmap map[string]interface{}
      _ = json.Unmarshal(msg, &objmap)
      event := objmap["event"].(string)
      sendData := map[string]interface{}{
         "event": "res",
         "data":  nil,
      }
      switch event {
      case "open":
         log.Printf("Received: %s\n", event)
      case "req":
         sendData["data"] = objmap["data"]
         log.Printf("Received: %s\n", event)
      }
      refineSendData, err := json.Marshal(sendData)
      err = ws.WriteMessage(mt, refineSendData)
      if err != nil {
         log.Println("write:", err)
         break
      }

		fmt.Printf("%s\n", msg)
	}
   return err
}

func main() {
   // 1. mysql 연결 -> storage.conn

   // 2. id generator 작성 : { id(pk), 글로벌 모델 배포 일자 , 글로벌 모델 버전 } 을 스키마로 갖는다.
   //      @ id      : 클라이언트를 구분하기위함, id generator에서 중복없이 생성하는 로직 작성
   //      @ 배포 일자 : aggregation 주기에 따라 라운드 스케쥴링에 관여하는 파라미터로 사용할 예정
   //                  ex ) 7일마다 aggregation -> 로컬에서 학습이 되어야하므로 글로벌 모델을 배포받은 뒤 최소 5일이 지난 클라이언트만 참여함
   //                   ㄴ> 해당 로직 없이 aggregation을 하면 글로벌 모델의 학습이 무뎌질 수 있음.
   
   // 3. aggregation 행렬 평균 연산 로직 작성
   
   // 4. 새로운 글로벌 모델을 평가하는 python program을 실행하는 로직 작성
   
   // 5. API 작성
   //     A. 클라이언트의 앱 init시 id 생성을 위한 http GET /client/id 작성
   //     B. 클라이언트가 앱 실행 시 항상 연결되어있을 WebSocket /model 작성

	e := echo.New()
   // read file to binary
   globalModel, err := ioutil.ReadFile("filename")
   if err != nil {
      panic(err)
   }
   fmt.Print(string(globalModel))
   // @Todo : send globalModel to all clients when arregation is done.
   //    goroutines : All recevice check ( true ) -> Arregation ( done ) -> Make tflite with python program ( done ) -> Send to all
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Logger.Fatal(e.Start(":3000"))
}