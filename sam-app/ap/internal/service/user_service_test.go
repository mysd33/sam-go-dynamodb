package service

import (
	"ap/internal/entity"
	"ap/internal/repository"
	"testing"

	"example.com/apbase/pkg/config"
	"example.com/apbase/pkg/id"
	"example.com/apbase/pkg/logging"
)

//TODO: とりあえずの手作りのMockをtestfyに置き換え
type MockUserRepository struct {
}

func (d *MockUserRepository) GetUser(userId string) (*entity.User, error) {
	return &entity.User{ID: userId, Name: "dummy"}, nil
}

func (d *MockUserRepository) PutUser(user *entity.User) (*entity.User, error) {
	user.ID = id.GenerateId()
	return user, nil
}

func mockRepository() repository.UserRepository {
	repository := MockUserRepository{}
	return &repository
}

func TestRegist(t *testing.T) {
	log := logging.NewLogger()
	cfg := &config.Config{Hoge: config.Hoge{Name: "hoge"}}
	repository := mockRepository()
	sut := NewUserService(log, cfg, &repository)
	userName := "fuga"
	actual, _ := sut.Regist(userName)
	println(actual)
	//TODO: testifyでAssert文を追加
}
