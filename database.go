package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	chroma "github.com/amikos-tech/chroma-go"
	"github.com/amikos-tech/chroma-go/openai"
	"github.com/amikos-tech/chroma-go/types"
	"github.com/joho/godotenv"
	gopenai "github.com/sashabaranov/go-openai"
)

type Database struct {
	client               *chroma.Client
	courseCollection     *chroma.Collection
	instructorCollection *chroma.Collection
	openaiEf             *openai.OpenAIEmbeddingFunction
	openaiClient         *gopenai.Client
}

func NewChromaDB(reset bool) (*Database, error) {
	if err := godotenv.Load("api.env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	apiKey := os.Getenv("api")
	if apiKey == "" {
		return nil, fmt.Errorf("api key not found")
	}

	ctx := context.Background()

	client, err := chroma.NewClient("http://0.0.0.0:8000")
	if err != nil {
		return nil, fmt.Errorf("error creating client: %s", err)
	}

	openaiEf, err := openai.NewOpenAIEmbeddingFunction(apiKey)
	if err != nil {
		fmt.Printf("error initializing OpenAI embedding function: %v", err)
	}

	var courseCollection, instructorCollection *chroma.Collection
	// if reset, delete the courses collection and create a new ones
	if reset {
		_, err = client.DeleteCollection(ctx, "courses")
		if err != nil {
			fmt.Printf("failed to delete courses collection: %v", err)
		}

		_, err = client.DeleteCollection(ctx, "instructors")
		if err != nil {
			fmt.Printf("failed to delete courses collection: %v", err)
		}

		// Create new collections
		_, err = client.CreateCollection(ctx, "courses", nil, true, openaiEf, types.L2)
		if err != nil {
			return nil, fmt.Errorf("failed to create courses collection: %v", err)
		}

		_, err = client.CreateCollection(ctx, "instructors", nil, true, openaiEf, types.L2)
		if err != nil {
			return nil, fmt.Errorf("failed to create instructors collection: %v", err)
		}
	}

	// get collections again
	courseCollection, err = client.GetCollection(ctx, "courses", openaiEf)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses collection: %v", err)
	}

	instructorCollection, err = client.GetCollection(ctx, "instructors", openaiEf)
	if err != nil {
		return nil, fmt.Errorf("failed to get create instructors collection: %v", err)
	}

	return &Database{
		client:               client,
		courseCollection:     courseCollection,
		instructorCollection: instructorCollection,
		openaiEf:             openaiEf,
		openaiClient:         gopenai.NewClient(apiKey),
	}, nil

}

func (db *Database) insertRecords(r io.Reader, reset bool) error {
	if !reset {
		return nil
	}

	ctx := context.Background()
	records, err := ReadCSV(r)
	if err != nil {
		return err
	}

	if len(records) < 1 {
		return fmt.Errorf("no records to insert")
	}

	fmt.Printf("Processing %d records from CSV\n", len(records)-1)

	instructors := make(map[string]bool)
	instructorsNames := []string{}
	var courseDocuments []courseRecord

	for _, row := range records[1:] {

		firstName := row[17]
		lastName := row[18]
		fullName := firstName + " " + lastName

		if !instructors[fullName] {
			instructors[fullName] = true
			instructorsNames = append(instructorsNames, fullName)
		}

		courseDocuments = append(courseDocuments, courseRecord{
			document:   strings.Join(row, " "),
			instructor: fullName,
		})

	}

	fmt.Printf("Found %d unique instructors and %d courses\n", len(instructors), len(courseDocuments))

	if err := db.insertCourses(ctx, courseDocuments); err != nil {
		return fmt.Errorf("failed to insert courses: %v", err)
	}

	if err := db.insertInstructor(instructorsNames); err != nil {
		return fmt.Errorf("failed to insert instructors: %v", err)
	}

	return nil
}

func (db *Database) Query(question string) ([][]string, error) {

	fmt.Printf("question is: %s\n", question)

	// Extract potential instructor name from the question
	instructorName, err := db.getInstructorFromQuestion(question)
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

		metadata = map[string]interface{}{"instructor": name.Documents[0][0]}
	}

	courseQR, err := db.courseCollection.Query(context.TODO(), []string{question}, 5, metadata, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error querying course collection: %v", err)
	}

	if len(courseQR.Documents) == 0 {
		return nil, fmt.Errorf("no matching courses found")
	}

	return courseQR.Documents, nil
}
