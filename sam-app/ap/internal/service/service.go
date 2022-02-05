package service

import (
	"ap/internal/entity"
	"ap/internal/repository"

	"example.com/apbase/pkg/config"
	"example.com/apbase/pkg/logging"
)

type UserService struct {
	Log        logging.Logger
	Config     *config.Config
	Repository *repository.UserRepository
}

func (us UserService) Regist(userName string) (*entity.User, error) {
	//TODO: Viperによる設定ファイルの読み込みのとりあえずの確認
	us.Log.Info("hoge.name=%s", us.Config.Hoge.Name)

	//Zapによるログ出力の例
	us.Log.Info("UserName=%s", userName)

	user := entity.User{}
	user.Name = userName
	return us.Repository.PutUser(&user)
}

func (us UserService) Find(userId string) (*entity.User, error) {
	return us.Repository.GetUser(userId)
}
