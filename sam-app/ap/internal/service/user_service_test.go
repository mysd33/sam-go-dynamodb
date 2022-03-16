package service

import (
	"ap/internal/entity"
	"ap/internal/repository"
	"testing"

	"example.com/apbase/pkg/config"
	"example.com/apbase/pkg/id"
	"example.com/apbase/pkg/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//testfyによるMockの定義
type MockUserRepository struct {
	mock.Mock
}

func (d *MockUserRepository) GetUser(userId string) (*entity.User, error) {
	//Mockの設定
	ret := d.Called(userId)
	return ret.Get(0).(*entity.User), nil
}

func (d *MockUserRepository) PutUser(user *entity.User) (*entity.User, error) {
	//Mockの設定
	ret := d.Called(user)
	return ret.Get(0).(*entity.User), nil
}

func TestRegist(t *testing.T) {
	//入力値
	inputUserName := "fuga"
	//期待値
	expectedName := "fuga"
	log := logging.NewLogger()
	cfg := &config.Config{Hoge: config.Hoge{Name: "hoge"}}

	//Mockの戻り値の設定
	mock := new(MockUserRepository)
	mockReturnValue := entity.User{ID: id.GenerateId(), Name: expectedName}
	mock.On("PutUser", &entity.User{Name: inputUserName}).Return(&mockReturnValue)
	var repository repository.UserRepository = mock
	sut := NewUserService(log, cfg, &repository)

	actual, _ := sut.Regist(inputUserName)
	println(actual)
	//testifyによるAssert文
	assert.Equal(t, expectedName, actual.Name)
}
