package main

import (
	"context"
	"fmt"

	gopenai "github.com/sashabaranov/go-openai"
)

type chatBot struct {
	question string
	data     string
}

func (bot *chatBot) callAPI(client *gopenai.Client) (string, error) {
	req := gopenai.ChatCompletionRequest{
		Model: gopenai.GPT4oMini,
		Messages: []gopenai.ChatCompletionMessage{
			{
				Role:    gopenai.ChatMessageRoleSystem,
				Content: bot.data,
			},
			{
				Role:    gopenai.ChatMessageRoleUser,
				Content: bot.question,
			},
		},
	}

	resp, err := client.CreateChatCompletion(context.TODO(), req)
	if err != nil {
		return "", fmt.Errorf("CreateChatCompletion failed: %v", err)
	}

	return resp.Choices[0].Message.Content, nil

}

func (db *Database) Query(question string) ([][]string, error) {

	fmt.Printf("question is: %s\n", question)

	// Extract potential instructor name from the question
	bot := &chatBot{
		question: question,
		data:     "Extract only the first and last name from the following text. If no name is present, respond with 'none'. Return only the name or 'none' as the answer.",
	}

	instructorName, err := bot.callAPI(db.openaiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to extract instructor from question: %v", err)
	}

	fmt.Printf("Instructor name extracted: %s\n", instructorName)

	var metadata map[string]interface{}

	if instructorName != "none" {
		name, err := db.instructorCollection.Query(context.TODO(), []string{instructorName}, 1, nil, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("error querying instructors collection: %v", err)
		}
		fmt.Printf("real name is: %s\n", name.Documents[0][0])

		metadata = map[string]interface{}{"instructor": name.Documents[0][0]}
	}

	courseQR, err := db.courseCollection.Query(context.TODO(), []string{question}, 40, metadata, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error querying course collection: %v", err)
	}

	if len(courseQR.Documents) == 0 {
		return nil, fmt.Errorf("no matching courses found")
	}

	return courseQR.Documents, nil
}
