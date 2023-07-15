package controller

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"line/model"
	"line/usecase"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type SendMessageController struct {
	SendMessageUseCase *usecase.SendMessageUseCase
}

//	func (c *SendMessageController) Handle(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json")
//		w.Header().Set("Access-Control-Allow-Origin", "*")
//		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
//		body, err := io.ReadAll(r.Body)
//		if err != nil {
//			log.Printf("fail: io.ReadAll, %v\n", err)
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//		var messageReq model.Message
//		if err := json.Unmarshal(body, &messageReq); err != nil {
//			log.Printf("fail: json.Unmarshal, %v\n", err)
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//
//		message, err := c.SendMessageUseCase.Handle(messageReq)
//		if err != nil {
//			log.Printf("fail: SendMessageUseCase.Handle, %v\n", err)
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//
//		messageRes := model.MessageResForHTTPPost{MessageId: message.MessageId}
//		bytes, err := json.Marshal(messageRes)
//		if err != nil {
//			log.Printf("fail: json.Marshal, %v\n", err)
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//
//		if r.Method == "OPTIONS" {
//			w.WriteHeader(http.StatusOK)
//			return
//		}
//		w.Write(bytes)
//	}
func parseMessage(rawMsg []byte) *model.Message {
	var msg model.Message
	err := json.Unmarshal(rawMsg, &msg)
	if err != nil {
		log.Println(err)
	}
	return &msg
}

func (c *SendMessageController) HandleWS(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		// Parse the message and save it to database
		msg := parseMessage(message) // This function would parse the received data into a model.Message
		newMsg, err := c.SendMessageUseCase.Handle(*msg)
		if err != nil {
			log.Println(err)
		}
		*msg = newMsg

		if err != nil {
			log.Println(err)
			return
		}

		if err = ws.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			return
		}
	}
}
