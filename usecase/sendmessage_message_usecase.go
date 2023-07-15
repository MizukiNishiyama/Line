package usecase

import (
	"github.com/oklog/ulid/v2"
	"line/dao"
	"line/model"
)

type SendMessageUseCase struct {
	MessageDao *dao.MessageDao
}

func (uc *SendMessageUseCase) Handle(message model.Message) (model.Message, error) {
	id := ulid.Make().String()
	messageToInsert := model.Message{
		MessageId:      id,
		MessageContent: message.MessageContent,
		MessageTime:    message.MessageTime,
		UserId:         message.UserId,
		UserName:       message.UserName,
		RoomId:         message.RoomId,
	}

	err := uc.MessageDao.Insert(messageToInsert)
	if err != nil {
		return model.Message{}, err
	}

	return messageToInsert, nil
}
