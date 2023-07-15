package controller

import (
	"encoding/json"
	"io"
	"line/model"
	"line/usecase"
	"log"
	"net/http"
)

type RegisterUserController struct {
	RegisterUserUseCase *usecase.RegisterUserUseCase
}

func (c *RegisterUserController) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("fail: io.ReadAll, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var userReq model.User
	if err := json.Unmarshal(body, &userReq); err != nil {
		log.Printf("fail: json.Unmarshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := c.RegisterUserUseCase.Handle(userReq)
	if err != nil {
		log.Printf("fail: RegisterUserUseCase.Handle, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userRes := model.UserResForHTTPPost{UserId: user.UserId}
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
