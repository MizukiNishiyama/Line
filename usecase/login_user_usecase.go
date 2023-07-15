package usecase

import (
	"line/dao"
	"line/model"
)

type LoginUserUseCase struct {
	UserDao *dao.UserDao
}

func (u *LoginUserUseCase) Login(user *model.User) (*model.User, error) {
	return u.UserDao.Login(user)
}
