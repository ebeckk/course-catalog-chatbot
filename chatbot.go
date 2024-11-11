package main

import (
	"context"
	"fmt"

	gopenai "github.com/sashabaranov/go-openai"
)

func (db *Database) getInstructorFromQuestion(question string) (string, error) {

	fmt.Printf("question recived: %s\n", question)

	req := gopenai.ChatCompletionRequest{
		Model: gopenai.GPT4oMini,
		Messages: []gopenai.ChatCompletionMessage{
			{
				Role:    gopenai.ChatMessageRoleSystem,
				Content: "Extract only the first and last name from the following text. If no name is present, respond with 'none'. Return only the name or 'none' as the answer." + question,
			},
			{
				Role:    gopenai.ChatMessageRoleUser,
				Content: question,
			},
		},
	}

	resp, err := db.openaiClient.CreateChatCompletion(context.TODO(), req)
	if err != nil {
		fmt.Println("CreateChatCompletion failed: ", err)
	}

	result := resp.Choices[0].Message.Content

	fmt.Printf("instructor name from api: %s\n", result)

	return result, nil
}
