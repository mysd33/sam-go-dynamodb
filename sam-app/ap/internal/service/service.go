package service

import (
	"ap/internal/entity"
	"ap/internal/repository"
)

type UserService struct {
	Repository *repository.UserRepository
}

func (us UserService) Regist(userName string) (*entity.User, error) {
	user := entity.User{}
	user.Name = userName
	return us.Repository.PutUser(&user)
}

func (us UserService) Find(userId string) (*entity.User, error) {
	return us.Repository.GetUser(userId)
}
