package main

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	// a handler for GET requests, used for troubleshooting
	if req.HTTPMethod == "GET" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       "Everything is working, this is the GET request body: " + req.Body,
		}, nil
	}
	// check if the request is a post
	if req.HTTPMethod != "POST" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 405,
			Body:       "The HTTP method is not allowed, make a POST request",
		}, errors.New("HTTP method is not valid")
	}
	// fmt.Println(req.Headers)
	// fmt.Println(req.Body)

	// check the content-type and proceed down the rabbit hole
	return CheckContentType(req)
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
