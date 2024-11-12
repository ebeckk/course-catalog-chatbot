package main

import (
	"context"
	"fmt"

	"github.com/amikos-tech/chroma-go/types"
)

type courseRecord struct {
	document   string
	instructor string
}

func (db *Database) insertCourses(courses []courseRecord) error {
	ctx := context.Background()

	fmt.Printf("Attempting to insert %d courses\n", len(courses))
	if len(courses) == 0 {
		fmt.Println("nothing to insert")
		return nil
	}

	const batchSize = 2000
	for i := 0; i < len(courses); i += batchSize {
		end := i + batchSize
		if end > len(courses) {
			end = len(courses)
		}

		rs, err := types.NewRecordSet(
			types.WithEmbeddingFunction(db.courseCollection.EmbeddingFunction),
			types.WithIDGenerator(types.NewUUIDGenerator()),
		)
		if err != nil {
			return fmt.Errorf("error creating course record set: %v", err)
		}

		for _, course := range courses[i:end] {
			rs.WithRecord(
				types.WithDocument(course.document),
				types.WithMetadata("instructor", course.instructor),
			)
			//fmt.Printf("course: %v, instructor: %v\n", course.document, course.instructor)
		}

		if _, err := rs.BuildAndValidate(ctx); err != nil {
			return fmt.Errorf("error building and validating courses: %v", err)
		}

		if _, err := db.courseCollection.AddRecords(ctx, rs); err != nil {
			return fmt.Errorf("error adding courses: %v", err)
		}
	}

	return nil
}

func (db *Database) insertInstructor(instructorNames []string) error {
	ctx := context.TODO()

	if len(instructorNames) == 0 {
		return nil
	}

	irs, err := types.NewRecordSet(
		types.WithEmbeddingFunction(db.instructorCollection.EmbeddingFunction),
		types.WithIDGenerator(types.NewUUIDGenerator()),
	)

	if err != nil {
		return fmt.Errorf("error creating instructor record set")
	}

	for _, instructor := range instructorNames {
		irs.WithRecord(types.WithDocument(instructor))
		fmt.Printf("added instructor: %s\n", instructor)
	}

	_, err = irs.BuildAndValidate(ctx)
	if err != nil {
		return fmt.Errorf("failed to build and validate instructors: %v", err)
	}
	fmt.Println("finished build and validate")

	fmt.Println("ITS ABOUT TO ADD RECORED")
	_, err = db.instructorCollection.AddRecords(ctx, irs)
	if err != nil {
		return fmt.Errorf("error adding instructors: %v", err)
	}

	fmt.Println("DONE WITH ADDING ")

	return nil
}
