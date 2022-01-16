package service

import "ap/internal/db"

type UserService struct {
	Repository *db.UserRepository
}

func (us UserService) Regist(userName string) (*db.User, error) {
	user := db.User{}
	user.Name = userName
	return us.Repository.PutUser(&user)
}

func (us UserService) Find(userId string) (*db.User, error) {
	return us.Repository.GetUser(userId)
}
