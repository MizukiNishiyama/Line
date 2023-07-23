package model

type Room struct {
	RoomId    string `json:"RoomId"`
	UserId1   string `json:"UserId1"`
	UserId2   string `json:"UserId2"`
	UserName1 string `json:"UserName1"`
	UserName2 string `json:"UserName2"`
}

type Follow struct {
	UserId           string `json:"UserId"`
	UserName         string `json:"UserName"`
	OpponentUserName string `json:"OpponentUserName"`
}

type RoomResForHTTPPost struct {
	RoomId string `json:"RoomId"`
}

type UserId struct {
	UserId string `json:"UserId"`
}
