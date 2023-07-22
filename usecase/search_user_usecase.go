package usecase

import (
	"line/dao"
	"line/model"
)

type SearchUserUseCase struct {
	UserDao *dao.UserDao
}

func (uc *SearchUserUseCase) Handle(UserName string) ([]model.User, error) {
	return uc.UserDao.FindById(UserName)
}
