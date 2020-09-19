package main

import (
        "context"
        "time"
	"os"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
        "github.com/mailgun/mailgun-go"
	"gopkg.in/ezzarghili/recaptcha-go.v4"
)

var my_domain string = os.Getenv("MAIL_DOMAIN")

var privateAPIKey string = os.Getenv("API_KEY")

var recaptchaSecret string = os.Getenv("RECAPTCHA_SECRET")

type BodyRequest struct {
	RecaptchaResponse string `json:"g-recaptcha-response"`
	ContactName string `json:"name"`
	ContactEmail string `json:"email"`
	ContactMessage string `json:"message"`
	BotCheck string `json:"comment"`
}

// BodyResponse is our self-made struct to build response for Client
type BodyResponse struct {
	ResponseMsg string `json:"output"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	bodyRequest := BodyRequest{
		RecaptchaResponse: "",
		ContactName: "",
		ContactEmail: "",
		ContactMessage: "",
		BotCheck: "",
	}

	err := json.Unmarshal([]byte(request.Body), &bodyRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	if len(bodyRequest.BotCheck) > 0 {
		return events.APIGatewayProxyResponse{Body: false, StatusCode: 404}, nil
	}

	captcha, _ := recaptcha.NewReCAPTCHA(recaptchaSecret, recaptcha.V3, 5 * time.Second)

	captcha_err := captcha.Verify(bodyRequest.RecaptchaResponse)

	if captcha_err != nil {
		return events.APIGatewayProxyResponse{Body: captcha_err.Error(), StatusCode: 404}, nil
	}

	sendResult := sendMail(bodyRequest.ContactEmail, bodyRequest.ContactName, bodyRequest.ContactMessage)

	// We will build the BodyResponse and send it back in json form
	bodyResponse := BodyResponse{
		ResponseMsg: sendResult,
	}

	// Marshal the response into json bytes, if error return 404
	response, err := json.Marshal(&bodyResponse)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 404}, nil
	}

	//Returning response with AWS Lambda Proxy Response
	return events.APIGatewayProxyResponse{Body: string(response), StatusCode: 200}, nil
}

// send function
func sendMail(from_email string, from_name string, message_body string) string {
    // Create an instance of the Mailgun Client
    mg := mailgun.NewMailgun(my_domain, privateAPIKey)

    sender := os.Getenv("SENDER")
    subject := os.Getenv("SUBJECT")
    body := message_body
    recipient := os.Getenv("RECIPIENT")

    // The message object allows you to add attachments and Bcc recipients
    message := mg.NewMessage(sender, subject, body, recipient)

    // Add entered reply to
    message.SetReplyTo(from_name + " <" + from_email +">")

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
    defer cancel()

    // Send the message with a 10 second timeout
    _, _, err := mg.Send(ctx, message)

    if err != nil {
        return err.Error()
    }
    return "success"
}

func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(handler)
}
