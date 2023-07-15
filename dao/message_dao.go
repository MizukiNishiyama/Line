package dao

import (
	"database/sql"
	"line/model"
)

type MessageDao struct {
	DB *sql.DB
}

// 新しいメッセージをDBに挿入
func (dao *MessageDao) Insert(message model.Message) error {
	_, err := dao.DB.Exec("INSERT into message VALUES(?, ?, ?, ? ,?, ?)", message.MessageId, message.MessageContent, message.MessageTime, message.UserId, message.RoomId, message.UserName)
	return err
}

// 特定のメッセージを取得
func (dao *MessageDao) FetchMessage(RoomId string) ([]model.Message, error) {
	rows, err := dao.DB.Query("SELECT * FROM message WHERE RoomId = ?", RoomId)
	if err != nil {
		return nil, err
	}

	messages := make([]model.Message, 0)
	for rows.Next() {
		var u model.Message
		if err := rows.Scan(&u.MessageId, &u.MessageContent, &u.MessageTime, &u.UserId, &u.RoomId, &u.UserName); err != nil {
			return nil, err
		}
		messages = append(messages, u)
	}

	return messages, nil
}
