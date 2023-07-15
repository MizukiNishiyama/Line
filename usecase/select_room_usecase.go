package usecase

import (
	"line/dao"
	"line/model"
)

type SelectRoomUseCase struct {
	RoomDao *dao.RoomDao
}

func (uc *SelectRoomUseCase) Handle(UserId string) ([]model.Room, error) {
	return uc.RoomDao.SearchRoom(UserId)
}
