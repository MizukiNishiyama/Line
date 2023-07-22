package controller

import (
	"encoding/json"
	"io"
	"line/model"
	"line/usecase"
	"log"
	"net/http"
)

type MakeRoomController struct {
	MakeRoomUseCase *usecase.MakeRoomUseCase
}

func (c *MakeRoomController) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("fail: io.ReadAll, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var userReq model.Follow
	if err := json.Unmarshal(body, &userReq); err != nil {
		log.Printf("fail: json.Unmarshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	room, err := c.MakeRoomUseCase.Handle(userReq)
	if err != nil {
		log.Printf("fail: MakeRoomUseCase.Handle, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userRes := model.RoomResForHTTPPost{RoomId: room.OpponentUserName}
	bytes, err := json.Marshal(userRes)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.Write(bytes)
}
