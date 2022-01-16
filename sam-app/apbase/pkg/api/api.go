package api

import (
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func OkResponse(result string) (events.APIGatewayProxyResponse, error) {
	return response(
		http.StatusOK,
		result), nil
}

func ErrorResponse(err error) (events.APIGatewayProxyResponse, error) {
	return response(
		http.StatusBadRequest,
		errorResponseBody(err.Error())), nil
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
