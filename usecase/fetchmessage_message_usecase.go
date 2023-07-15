package usecase

import (
	"line/dao"
	"line/model"
)

type FetchMessageUseCase struct {
	MessageDao *dao.MessageDao
}

func (uc *FetchMessageUseCase) Handle(RoomId string) ([]model.Message, error) {
	return uc.MessageDao.FetchMessage(RoomId)
}
