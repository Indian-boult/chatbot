package chat

import (
	"chatbot/pkg/config"
	"context"
	"encoding/base64"
	"fmt"
	openai2 "github.com/sashabaranov/go-openai"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// OpenAI client structure
type OpenAI struct {
	Client *openai2.Client
}

// NewOpenAI initializes a new OpenAI client
func NewOpenAI(cfg *config.Config) *OpenAI {
	client := openai2.NewClient(cfg.OpenAIKey)
	return &OpenAI{
		Client: client,
	}
}

// AskOpenAI sends the chat message to OpenAI API and gets the response.
// This function is used for normal conversation responses from OpenAI.
func (o *OpenAI) AskOpenAI(ctx context.Context, prompt string) (string, error) {
	req := openai2.ChatCompletionRequest{
		Model: openai2.GPT3Dot5Turbo, // Adjust as necessary
		Messages: []openai2.ChatCompletionMessage{
			{
				Role:    openai2.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	resp, err := o.Client.CreateChatCompletion(ctx, req)
	if err != nil {
		log.Printf("AskOpenAI error: %v\n", err)
		return "", err
	}

	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
		return resp.Choices[0].Message.Content, nil
	}
	return "No response from OpenAI.", nil
}

// ImageIntentChecker checks the intent of the message regarding images.
func (o *OpenAI) ImageIntentChecker(ctx context.Context, msg string) (string, error) {
	log.Printf("ImageIntentChecker: Received message: %s\n", msg)
	systemMessage := "User is engaging with a chatbot designed for several tasks including normal conversations, retrieving specific images, and discussing content within images. When the user mention 'retrieve' or 'show' then give output as retrieve_image:<identifier>, if user mention 'discuss' or 'know about' or 'know more' then give output as discuss_image:<identifier>, if user mention that they want to upload the image then give output upload_image, and if nothing matches out of these three then give output normal_conversation. Provide only one of these output. And don't provide any extension for ex: .jpeg, .png, etc"
	resp, err := o.Client.CreateChatCompletion(ctx, openai2.ChatCompletionRequest{
		Model: openai2.GPT3Dot5Turbo,
		Messages: []openai2.ChatCompletionMessage{
			{Role: openai2.ChatMessageRoleSystem, Content: systemMessage},
			{Role: openai2.ChatMessageRoleUser, Content: msg},
		},
	})
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	intent := resp.Choices[0].Message.Content
	log.Printf("ImageIntentChecker: Detected intent: %s\n", intent)
	return intent, nil
}

func (o *OpenAI) DiscussImage(ctx context.Context, imgURL string) (string, error) {
	imgString, err := getImageBase64(imgURL)
	if err != nil {
		return "", fmt.Errorf("error converting image to base64: %v", err)
	}

	resp, err := o.Client.CreateChatCompletion(ctx, openai2.ChatCompletionRequest{
		Model:     openai2.GPT4VisionPreview, // Assuming this is the model for vision tasks
		MaxTokens: 100,
		Messages: []openai2.ChatCompletionMessage{
			{
				Role: openai2.ChatMessageRoleUser,
				MultiContent: []openai2.ChatMessagePart{
					{
						Type: openai2.ChatMessagePartTypeText,
						Text: "What's in this image?",
					},
					{
						Type: openai2.ChatMessagePartTypeImageURL,
						ImageURL: &openai2.ChatMessageImageURL{
							URL:    fmt.Sprintf("data:image/jpg;base64,%s", imgString),
							Detail: openai2.ImageURLDetailLow,
						},
					},
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("error discussing image: %v", err)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}
	return "No response from OpenAI on the image.", nil
}

func getImageBase64(imgPath string) (string, error) {
	// Open the image file
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %v", err)
	}
	defer imgFile.Close()

	// Read the image file's content
	imgData, err := ioutil.ReadAll(imgFile)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %v", err)
	}

	// Encode the file content into base64
	imgBase64 := base64.StdEncoding.EncodeToString(imgData)

	return imgBase64, nil
}

// GeneratePromptWithHistory generates a prompt including the conversation history for context.
func (o *OpenAI) GeneratePromptWithHistory(history []string, userInput string) string {
	var sb strings.Builder
	for _, msg := range history {
		sb.WriteString(msg + "\n")
	}
	sb.WriteString(userInput)
	return sb.String()
}
