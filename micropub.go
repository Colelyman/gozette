package main

import (
	"errors"
	"fmt"

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

	// check the content-type
	contentType, err := GetContentType(req.Headers)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, err
	}
	entry, err := CreateEntry(contentType, req.Body)
	fmt.Printf("req.Headers: %v\n", req.Headers)
	fmt.Printf("req.Headers[authorization]: %s\n", req.Headers["authorization"])
	if CheckAuthorization(entry, req.Headers) {
		location, err := WriteEntry(entry)
		if err != nil {
			return &events.APIGatewayProxyResponse{
				StatusCode: 400,
				Body:       "There was an error committing the entry to the repository",
			}, errors.New("Error committing the entry to the repository")
		} else {
			return &events.APIGatewayProxyResponse{
				StatusCode: 202,
				Headers:    map[string]string{"Location": location},
			}, nil
		}
	} else {
		return &events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "Forbidden, there was a problem with the provided access token",
		}, errors.New("The provided access token does not grant permission")
	}
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
