package usecase

import (
	"line/dao"
	"line/model"
)

type MakeRoomUseCase struct {
	RoomDao *dao.RoomDao
}

func (uc *MakeRoomUseCase) Handle(follow model.Follow) (model.Follow, error) {
	userToInsert := model.Follow{
		UserId:           follow.UserId,
		UserName:         follow.UserName,
		OpponentUserName: follow.OpponentUserName,
	}

	err := uc.RoomDao.MakeRoom(userToInsert)
	if err != nil {
		return model.Follow{}, err
	}

	return userToInsert, nil
}
