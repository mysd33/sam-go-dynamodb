package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"users/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

var (
	// ErrResponse
	ErrResponse = errors.New("Error")
)
var UserRepository db.UserRepository

type request struct {
	Name string `json:"name"`
}

func init() {
	UserRepository = db.New()
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//ctxの格納
	UserRepository.Context = ctx

	//Get
	if request.HTTPMethod == http.MethodGet {
		return getHandler(ctx, request)
	}
	//Post
	return postHandler(ctx, request)
}

func getHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userId, err := parseGetRequest(request)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}
	//result, err := UserRepository.GetUser(userId, ctx)
	result, err := UserRepository.GetUser(userId)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}
	return response(
		http.StatusOK,
		fmt.Sprintf("UserName: %v", result.Name),
	), nil
}

func postHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//parse Request
	p, err := parsePostRequest(request)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}

	user := db.User{}
	user.Name = p.Name
	//result, err := UserRepository.PutUser(&user, ctx)
	result, err := UserRepository.PutUser(&user)
	if err != nil {
		return response(
			http.StatusBadRequest,
			errorResponseBody(err.Error()),
		), nil
	}

	return response(
		http.StatusOK,
		fmt.Sprintf("UserId: %v", result.ID),
	), nil
}

func parseGetRequest(req events.APIGatewayProxyRequest) (string, error) {
	if req.HTTPMethod != http.MethodGet {
		return "", fmt.Errorf("use GET request")
	}
	userId := req.PathParameters["user_id"]
	return userId, nil
}

func parsePostRequest(req events.APIGatewayProxyRequest) (*request, error) {
	var r request
	if req.HTTPMethod != http.MethodPost {
		return &r, fmt.Errorf("use POST request")
	}

	err := json.Unmarshal([]byte(req.Body), &r)
	if err != nil {
		return &r, errors.Wrapf(err, "failed to parse request")
	}

	if err != nil {
		return &r, errors.Wrapf(err, "invalid URL")
	}

	return &r, nil
}

func response(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func errorResponseBody(msg string) string {
	return fmt.Sprintf("{\"message\":\"%s\"}", msg)
}

func main() {
	lambda.Start(handler)
}
