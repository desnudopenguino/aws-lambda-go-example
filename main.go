package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := request.QueryStringParameters["name"]
	response := fmt.Sprintf("Hello %s", name)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       response,
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
