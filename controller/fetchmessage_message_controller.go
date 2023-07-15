package controller

import (
	"encoding/json"
	"line/model"
	"line/usecase"
	"log"
	"net/http"
)

type FetchMessageController struct {
	FetchMessageUseCase *usecase.FetchMessageUseCase
}

func (c *FetchMessageController) Handle(w http.ResponseWriter, r *http.Request) {
	RoomId := r.URL.Query().Get("roomid")
	if RoomId == "" {
		log.Println("fail: userid is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messages, err := c.FetchMessageUseCase.Handle(RoomId)
	if err != nil {
		log.Printf("fail: FetchMessageUseCase.Handle, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messagesRes := make([]model.Message, len(messages))
	for i, u := range messages {
		messagesRes[i] = model.Message{MessageId: u.MessageId, MessageContent: u.MessageContent, MessageTime: u.MessageTime, UserId: u.UserId, RoomId: u.RoomId, UserName: u.UserName}
	}

	bytes, err := json.Marshal(messagesRes)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(bytes)
}
