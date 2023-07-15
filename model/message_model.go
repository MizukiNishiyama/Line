package model

type Message struct {
	MessageId      string `json:"MessageId"`
	MessageContent string `json:"MessageContent"`
	MessageTime    string `json:"MessageTime"`
	UserId         string `json:"UserId"`
	RoomId         string `json:"RoomId"`
	UserName       string `json:"UserName"`
}

type MessageResForHTTPPost struct {
	MessageId string `json:"MessageId"`
}
