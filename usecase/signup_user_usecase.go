package usecase

import (
	"github.com/oklog/ulid/v2"
	"line/dao"
	"line/model"
)

type RegisterUserUseCase struct {
	UserDao *dao.UserDao
}

func (uc *RegisterUserUseCase) Handle(user model.User) (model.User, error) {
	UserId := ulid.Make().String()
	userToInsert := model.User{
		UserId:       UserId,
		UserName:     user.UserName,
		UserPassword: user.UserPassword,
	}

	err := uc.UserDao.Signup(userToInsert)
	if err != nil {
		return model.User{}, err
	}

	return userToInsert, nil
}
