package model

type User struct {
	UserId       string `json:"UserId"`
	UserName     string `json:"UserName"`
	UserPassword string `json:"UserPassword"`
}

type UserResForHTTPPost struct {
	UserId string `json:"UserId"`
}
