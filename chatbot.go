package main

import (
	"context"
	"fmt"
	"strings"

	gopenai "github.com/sashabaranov/go-openai"
)

func (db *Database) getInstructorFromQuestion(question string) (string, error) {
	prompt := gopenai.ChatCompletionMessage{
		Role: gopenai.ChatMessageRoleUser,
		Content: fmt.Sprintf(
			`Extract the instructor name from this question. If there's no instructor name, return "NONE". Only return the name, nothing else.\nQuestion: %s`,
			question,
		),
	}

	req := gopenai.ChatCompletionRequest{
		Model:    gopenai.GPT4oMini,
		Messages: []gopenai.ChatCompletionMessage{prompt},
	}

	resp, err := db.openaiClient.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("failed to extract instructor name: %v", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return strings.TrimSpace(resp.Choices[0].Message.Content), nil
}
