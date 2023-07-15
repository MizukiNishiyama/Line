package controller

import (
	"encoding/json"
	"line/model"
	"line/usecase"
	"net/http"
)

type LoginUserController struct {
	LoginUserUseCase *usecase.LoginUserUseCase
}

func (u *LoginUserController) Handle(w http.ResponseWriter, r *http.Request) {
	user := &model.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err = u.LoginUserUseCase.Login(user)
	if err != nil {
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	// レスポンスをJSON形式で送る
	response := map[string]interface{}{
		"message": "Login successful",
		"UserId":  user.UserId,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
