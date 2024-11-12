package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

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

	fmt.Printf("\nCatalog search> ")

	for scanner.Scan() {
		question := scanner.Text()
		result, err := db.Query(question)
		if err != nil {
			return
		}
		var allQuery string
		for _, row := range result[0] {
			allQuery = allQuery + row
		}

		bot := &chatBot{
			question: question,
			data:     "use the following data to answer the question and format it, this is the data:" + allQuery,
		}

		answer, err := bot.callAPI(client)
		if err != nil {
			fmt.Printf("api call didn't work: %v", err)
		}

		fmt.Println(answer)
		fmt.Print("\nCatalog search> ")

	}

}
