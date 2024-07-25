package main

type ChatGPTRequest struct {
    Model    string   `json:"model"`
    Messages []RequestMessage `json:"messages"`
}

type RequestMessage struct {
    Role    string `json:"role"`
    Content []RequestContentItem `json:"content"`
}

type RequestContentItem struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL RequestImageURL  `json:"image_url,omitempty"`
}

type RequestImageURL struct {
	URL string `json:"url"`
}

type ResponseContentItem struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
}

type ResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTResponse struct {
    Choices []Choice `json:"choices"`
}

type Choice struct {
    Message ResponseMessage `json:"message"`
}

type MyEvent struct {
    Base64_image string `json:"base64_image"`
}

type ParsedEvent struct { // TODO improve struct names
	Name  string `json:"name"`
	Time  string `json:"time"`
	Stage string `json:"stage"`
}

type Response struct {
    StatusCode int               `json:"statusCode"`
    Headers    map[string]string `json:"headers"`
    Body       string            `json:"body"`
}
