package main

import (
	"context"
	"encoding/json"
	"net/http"

	"example.com/apbase/pkg/api"
	"example.com/apbase/pkg/config"
	"example.com/apbase/pkg/logging"
	"go.uber.org/zap"

	"ap/internal/entity"
	"ap/internal/repository"
	"ap/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

var (
	// Service
	userService service.UserService
	// Repository
	userRepository repository.UserRepository
	// Logger
	log logging.Logger
	// Config
	cfg *config.Config
)

//リクエストデータ
type Request struct {
	Name string `json:"name"`
}

//コードルドスタート時の初期化処理
func init() {
	var err error
	z, _ := zap.NewProduction()
	log = logging.ZapLogger{Log: z.Sugar()}
	cfg, err = config.LoadConfig()

	userRepository = repository.NewUserRepository()
	userService = service.UserService{Repository: &userRepository, Log: log, Config: cfg}

	if err != nil {
		//TODO: エラーハンドリング
		log.Fatal("初期化処理エラー:%s", err.Error())
		panic(err.Error())
	}

}

//ハンドラメソッド
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//TODO: dynamoDBのAP基盤機能側でContext格納するようにリファクタ
	//ctxの格納
	userRepository.Context = ctx

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
		return "", errors.Errorf("use GET request")
	}
	userId := req.PathParameters["user_id"]
	return userId, nil
}

//Postリクエストデータの解析
func parsePostRequest(req events.APIGatewayProxyRequest) (*Request, error) {
	var r Request
	if req.HTTPMethod != http.MethodPost {
		return &r, errors.Errorf("use POST request")
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
