package controller

import (
	"encoding/json"
	"line/model"
	"line/usecase"
	"log"
	"net/http"
)

type SearchUserController struct {
	SearchUserUseCase *usecase.SearchUserUseCase
}

func (c *SearchUserController) Handle(w http.ResponseWriter, r *http.Request) {
	UserName := r.URL.Query().Get("id")
	if UserName == "" {
		log.Println("fail: name is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := c.SearchUserUseCase.Handle(UserName)
	if err != nil {
		log.Printf("fail: SearchUserUseCase.Handle, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	usersRes := make([]model.User, len(users))
	for i, u := range users {
		usersRes[i] = model.User{UserId: u.UserId, UserName: u.UserName, UserPassword: u.UserPassword}
	}

	bytes, err := json.Marshal(usersRes)
	if err != nil {
		log.Printf("fail: json.Marshal, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(bytes)
}
