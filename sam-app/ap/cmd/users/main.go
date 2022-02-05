package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"example.com/apbase/pkg/api"

	"ap/internal/entity"
	"ap/internal/repository"
	"ap/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	// Service
	userService service.UserService
	// Repository
	userRepository repository.UserRepository
	// Config
	config *Config
)

//設定ファイルの構造体(Viper)
type Config struct {
	Hoge Hoge `yaml:hoge`
}

//TODO: とりあえずのサンプル
type Hoge struct {
	Name string `yaml:name`
}

//リクエストデータ
//TODO: request → Request
type request struct {
	Name string `json:"name"`
}

//設定ファイルのロード
func loadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")
	// 環境変数がすでに指定されてる場合はそちらを優先させる
	viper.AutomaticEnv()
	// データ構造をキャメルケースに切り替える用の設定
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.Errorf("設定ファイル読み込みエラー")
	}
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, errors.Errorf("設定ファイルアンマーシャルエラー")
	}
	return &cfg, nil
}

//コードルドスタート時の初期化処理
func init() {
	userRepository = repository.NewUserRepository()
	userService = service.UserService{Repository: &userRepository}

	var err error
	config, err = loadConfig()
	if err != nil {
		//TODO: エラーハンドリング
		log.Fatalf("初期化処理エラー:%s", err.Error())
		panic(err.Error())
	}

}

//ハンドラメソッド
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//TODO: dynamoDBのAP基盤機能側でContext格納するようにリファクタ
	//ctxの格納
	userRepository.Context = ctx

	//TODO: とりあえずの設定ファイルの読み込み確認
	log.Printf("hoge.name=%s", config.Hoge.Name)

	//Getリクエストの処理
	if request.HTTPMethod == http.MethodGet {
		return getHandler(ctx, request)
	}
	//Postリクエストの処理
	return postHandler(ctx, request)
}

//Getリクエストの処理
func getHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//リクエストデータの解析
	userId, err := parseGetRequest(request)
	if err != nil {
		return api.ErrorResponse(err)
	}
	//サービスの実行
	result, err := userService.Find(userId)
	if err != nil {
		return api.ErrorResponse(err)
	}
	//レスポンスデータの返却
	resultString := formatResponse(result)
	return api.OkResponse(resultString)
}

//Postリクエストの処理
func postHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//リクエストデータの解析
	p, err := parsePostRequest(request)
	if err != nil {
		return api.ErrorResponse(err)
	}
	//サービスの実行
	result, err := userService.Regist(p.Name)
	if err != nil {
		return api.ErrorResponse(err)
	}
	//レスポンスデータの返却
	resultString := formatResponse(result)
	return api.OkResponse(resultString)
}

//Getリクエストデータの解析
func parseGetRequest(req events.APIGatewayProxyRequest) (string, error) {
	if req.HTTPMethod != http.MethodGet {
		return "", fmt.Errorf("use GET request")
	}
	userId := req.PathParameters["user_id"]
	return userId, nil
}

//Postリクエストデータの解析
func parsePostRequest(req events.APIGatewayProxyRequest) (*request, error) {
	var r request
	if req.HTTPMethod != http.MethodPost {
		return &r, fmt.Errorf("use POST request")
	}

	err := json.Unmarshal([]byte(req.Body), &r)
	if err != nil {
		return &r, errors.Wrapf(err, "failed to parse request")
	}

	return &r, nil
}

//レスポンスデータの生成
func formatResponse(user *entity.User) string {
	resp, _ := json.Marshal(user)
	return string(resp)
}

//Main関数
func main() {
	lambda.Start(handler)
}
