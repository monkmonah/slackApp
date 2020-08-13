package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/slack-go/slack"
)

var api = slack.New("YOUR_TOKEN")
var signingSecret = "YOUR_SIGNING_SECRET"

// You can open a dialog with a user interaction. (like pushing buttons, slash commands ...)
// https://api.slack.com/surfaces/modals
// https://api.slack.com/interactivity/entry-points
func main() {
	lambda.Start(handler)
}
//func handler(w http.ResponseWriter, r *http.Request) {
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	headers := make(http.Header)
	for key, values := range request.MultiValueHeaders {
		for _, value := range values {
			headers.Add(key, value)
		}
	}

	// Verify signing secret
	sv, err := slack.NewSecretsVerifier(headers, signingSecret)
	if err != nil {
		log.Printf("[ERROR] Fail to secrets verifier: %v", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Fail to secrets verifier",
		}, nil
	}
	body := []byte(request.Body)
	_, _ = sv.Write(body)

	if err := sv.Ensure(); err != nil {
		log.Printf("[ERROR] Fail to Ensure: %v", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Fail to Ensure",
		}, nil
	}

	// Parse request body
	str, _ := url.QueryUnescape(string(body))
	str = strings.Replace(str, "payload=", "", 1)
	var message slack.InteractionCallback
	if err := json.Unmarshal([]byte(str), &message); err != nil {
		log.Printf("[ERROR] Fail to unmarshal json: %v", err)
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Fail to unmarshal json",
		}, nil
	}

	switch message.Type {
	case slack.InteractionTypeInteractionMessage:
		// Make new dialog components and open a dialog.
		// Component-Text
		textInput := slack.NewTextInput("TextSample", "Sample label - Text", "Default value")

		// Component-TextArea
		textareaInput := slack.NewTextAreaInput("TexaAreaSample", "Sample label - TextArea", "Default value")

		// Component-Select menu
		option1 := slack.DialogSelectOption{
			Label: "Display name 1",
			Value: "Inner value 1",
		}
		option2 := slack.DialogSelectOption{
			Label: "Display name 2",
			Value: "Inner value 2",
		}
		options := []slack.DialogSelectOption{option1, option2}
		selectInput := slack.NewStaticSelectDialogInput("SelectSample", "Sample label - Select", options)

		// Open a dialog
		elements := []slack.DialogElement{
			textInput,
			textareaInput,
			selectInput,
		}
		dialog := slack.Dialog{
			CallbackID:  "Callback_ID",
			Title:       "Dialog title",
			SubmitLabel: "Submit",
			Elements:    elements,
		}
		api.OpenDialog(message.TriggerID, dialog)

	case slack.InteractionTypeDialogSubmission:
		// Receive a notification of a dialog submission
		log.Printf("Successfully receive a dialog submission.")
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Ok",
	}, nil
}