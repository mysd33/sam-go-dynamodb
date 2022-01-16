package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/apbase/pkg/api"

	"ap/internal/db"
	apsvc "ap/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

var (
	// Service
	userService apsvc.UserService
	// Repository
	userRepository db.UserRepository
)

//リクエストデータ
type request struct {
	Name string `json:"name"`
}

//コードルドスタート時の初期化処理
func init() {
	userRepository = db.NewUserRepository()
	userService = apsvc.UserService{Repository: &userRepository}
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
func formatResponse(user *db.User) string {
	resp, _ := json.Marshal(user)
	return string(resp)
}

//Main関数
func main() {
	lambda.Start(handler)
}
