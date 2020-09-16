package main

import (
	"os"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type BodyRequest struct {
	ContactName string `json:"name"`
	ContactEmail string `json:"email"`
	ContactMessage string `json:"message"`
}

// BodyResponse is our self-made struct to build response for Client
type BodyResponse struct {
	ResponseMsg string `json:"output"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bodyRequest := BodyRequest{
		ContactName: "",
		ContactEmail: "",
		ContactMessage: "",
	}

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	// We will build the BodyResponse and send it back in json form
	bodyResponse := BodyResponse{
		ResponseMsg: bodyRequest.ContactName +" ("+ bodyRequest.ContactEmail +") says: "+ bodyRequest.ContactMessage +" "+ os.Getenv("LAST_NAME"),
	}

	// Marshal the response into json bytes, if error return 404
	response, err := json.Marshal(&bodyResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
