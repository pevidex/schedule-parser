package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-resty/resty/v2"
)

func GenerateGPTPromptPaylod(imageStr string) string {
	contentDataTemplate := `{
		"role": "user",
		"content": [
			{
				"type": "text",
				"text": "Is this image a schedule? If no, return just 'No' with no more text. If it is, return a full json file containing the list of all events in the following format: {\"events\": [{\"name\": \"event_name\", \"time\": \"20:00-21:00\", \"stage\": \"stage_name\"}]} Don't add any more text besides this json file."
			},
			{
				"type": "image_url",
				"image_url": {
					"url": "+imageContent+"
				}
			}
		]
	}`
	imageContent := "data:image/jpeg;base64," + imageStr
	contentRequest := strings.ReplaceAll(contentDataTemplate, "+imageContent+", imageContent)
   	return contentRequest
}

func RequestGPTApi(message RequestMessage) (*resty.Response, error) {
	// Prepare the ChatGPT request
    chatGPTRequest := ChatGPTRequest{
        Model: "gpt-4o-mini",
        Messages: []RequestMessage{message},
    }
    // Make the request to the OpenAI API
    client := resty.New()
    apiKey := os.Getenv("OPENAI_API_KEY")
    resp, err := client.R().
        SetHeader("Content-Type", "application/json").
        SetHeader("Authorization", "Bearer "+apiKey).
        SetBody(chatGPTRequest).
        Post("https://api.openai.com/v1/chat/completions")
	return resp, err
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
    var event MyEvent
    err := json.Unmarshal([]byte(request.Body), &event)
    if err != nil {
        return Response{StatusCode: http.StatusBadRequest, Body: "Invalid request body"}, nil
	}

	gptPayload := GenerateGPTPromptPaylod(event.Base64_image)

	var message RequestMessage
	err = json.Unmarshal([]byte(gptPayload), &message)

	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return Response{StatusCode: http.StatusInternalServerError, Body: "Failed to parse template"}, nil
	}

    resp, err := RequestGPTApi(message)

    if err != nil {
        return Response{StatusCode: http.StatusInternalServerError, Body: "Failed to process request"}, nil
    }
    var chatGPTResponse ChatGPTResponse
    err = json.Unmarshal(resp.Body(), &chatGPTResponse)
    if err != nil {
        return Response{StatusCode: http.StatusInternalServerError, Body: "Failed to process request"}, nil
    }
	
	rawJSON := chatGPTResponse.Choices[0].Message.Content
	start := strings.Index(rawJSON, "```json\n") + len("```json\n")
	rawJSON = rawJSON[start:]

    return Response{
        StatusCode: http.StatusOK,
        Headers:    map[string]string{"Content-Type": "application/json"},
        Body:       rawJSON,
    }, nil
}
