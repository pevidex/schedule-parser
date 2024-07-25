package main_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	main "github.com/pevidex/schedule_parser"
	"github.com/stretchr/testify/assert"
)

func encodeImageToBase64(filePath string) (string, error) {
	// Read the image file
	imageData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	// Encode the image data to Base64
	encodedImage := base64.StdEncoding.EncodeToString(imageData)
	return encodedImage, nil
}

func TestGPTHandler(t *testing.T) {
	b64Image, err := encodeImageToBase64("./test_image.jpg")
	if err != nil {
		assert.Fail(t, "Failed to load test_image.jpg")
	}

	bodyBytes, err := json.Marshal(struct {Base64_image string}{Base64_image: b64Image})
	if err != nil {
		assert.Fail(t, "Failed to load test_image.jpg")
	}

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: string(bodyBytes)},
			expect:  "Hello Paul",
			err:     nil,
		},
		{
			// Test that the handler responds ErrNameNotProvided
			// when no name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: ""},
			expect:  "",
			err:     nil,
		},
	}

	for _, test := range tests {
		response, err := main.HandleRequest(context.TODO(), test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}

}