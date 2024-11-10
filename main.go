package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
)

func GetAnswer(client *openai.Client, question string, query string) string {
	if query == "" {
		fmt.Println("no query data")
	} else {
		fmt.Println(query)
	}

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "use the following data to answer the question and format it nicley, this is the data:" + query,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: question,
			},
		},
	}

	resp, err := client.CreateChatCompletion(context.TODO(), req)
	if err != nil {
		fmt.Println("CreateChatCompletion failed: ", err)
	}

	return resp.Choices[0].Message.Content
}

func main() {
	reset := flag.Bool("reset", false, "set true to reset")
	flag.Parse()

	file, err := os.Open("fallclasses.csv")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	db, err := NewChromaDB(*reset)
	if err != nil {
		fmt.Printf("couldn't create new chromaDB: %v", err)
	}

	db.insertRecords(file, *reset)

	client := openai.NewClient(os.Getenv("api"))

	scanner := bufio.NewScanner(os.Stdin)
	time.Sleep(5 * time.Second)
	fmt.Printf("\nCatalog search> ")

	for scanner.Scan() {
		question := scanner.Text()
		result, err := db.Query(question)
		if err != nil {
			return
		}
		answer := GetAnswer(client, question, result)
		fmt.Println(answer)
		fmt.Print("\nCatalog search> ")

	}

}
