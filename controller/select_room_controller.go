package controller

import (
	"encoding/json"
	"line/model"
	"line/usecase"
	"log"
	"net/http"
)

type SelectRoomController struct {
	SelectRoomUseCase *usecase.SelectRoomUseCase
}

func (c *SelectRoomController) Handle(w http.ResponseWriter, r *http.Request) {
	UserId := r.URL.Query().Get("userid")
	if UserId == "" {
		log.Println("fail: userid is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rooms, err := c.SelectRoomUseCase.Handle(UserId)
	if err != nil {
		log.Printf("fail: SelectRoomUseCase.Handle, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	roomsRes := make([]model.Room, len(rooms))
	for i, u := range rooms {
		roomsRes[i] = model.Room{RoomId: u.RoomId, UserId1: u.UserId1, UserId2: u.UserId2, UserName1: u.UserName1, UserName2: u.UserName2}
	}

	bytes, err := json.Marshal(roomsRes)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(bytes)
}
