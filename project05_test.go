package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestProject5(t *testing.T) {
	tests := []struct {
		name     string
		question string
		expected string
	}{
		{
			name:     "TestPhil",
			question: "What courses is Phil Peterson teaching in Fall 2024?",
			expected: `Philip Peterson is teaching the following courses in Fall 2024:

1. CS 272 (Software Development)
   - In-Person
   - Schedule: Tuesday and Thursday, 8:00 AM - 9:45 AM
   - Location: LS G12

2. SCCS 272L 01 (Software Development Lab)
   - In-Person
   - Schedule: Wednesday, 1:00 PM - 2:30 PM
   - Location: MH 122

3. SCCS 272L 02 (Software Development Lab)
   - In-Person
   - Schedule: Wednesday, 2:55 PM - 4:25 PM
   - Location: MH 122

4. CS 272 (Software Development)
   - In-Person
   - Schedule: Tuesday and Thursday, 2:40 PM - 4:25 PM
   - Location: LS G12
			`,
		},
		{
			name:     "TestPHIL",
			question: "Which philosophy courses are offered this semester?",
			expected: `
The philosophy courses offered this semester are:

1. **PHIL 110 03**: Great Philosophical Questions (In-Person)
   - Schedule: MWF 11:45 AM - 12:50 PM
   - Instructor: Jea Oh (joh18@usfca.edu)

2. **PHIL 110 01**: Great Philosophical Questions (In-Person)
   - Schedule: MWF 9:15 AM - 10:20 AM
   - Instructor: Deena Lin (dmlin@usfca.edu)

3. **PHIL 110 02**: Great Philosophical Questions (In-Person)
   - Schedule: MWF 1:00 PM - 2:05 PM
   - Instructor: Jea Oh (joh18@usfca.edu)

4. **PHIL 110 04**: Great Philosophical Questions (Hybrid)
   - Schedule: RE F 3:30 PM - 4:35 PM (Online)
   - Schedule: IP MW 3:30 PM - 4:35 PM (In-Person)
   - Instructor: Richie Kim (rkim7@usfca.edu)

5. **PHIL 204 01**: Philosophy of Science (In-Person)
   - Schedule: MW 4:45 PM - 6:25 PM
   - Instructor: Krupa Patel (kpatel28@usfca.edu)

6. **PHIL 204 02**: Philosophy of Science (In-Person)
   - Schedule: MW 6:30 PM - 8:15 PM
   - Instructor: Krupa Patel (kpatel28@usfca.edu)

7. **PHIL 310 01**: Ancient & Medieval Philosophy (In-Person)
   - Schedule: TR 9:55 AM - 11:40 AM
   - Instructor: Thomas Cavanaugh (cavanaught@usfca.edu)

8. **PHIL 220 01**: Asian Philosophy (In-Person)
   - Schedule: MW 4:45 PM - 6:25 PM
   - Instructor: Joshua Stoll (jstoll@usfca.edu)

9. **PHIL 202 01**: Philosophy of Religion (In-Person)
   - Schedule: MWF 1:00 PM - 2:05 PM
   - Instructor: Deena Lin (dmlin@usfca.edu)

10. **PHIL 205 01**: Philosophy of Biology (In-Person)
    - Schedule: MWF 10:30 AM - 11:35 AM
    - Instructor: Stephen Friesen (smfriesen@usfca.edu)

11. **PHIL 240 07**: Ethics (In-Person)
    - Schedule: MWF 1:00 PM - 2:05 PM
    - Instructor: Vida Pavesich (vpavesich@usfca.edu)

12. **PHIL 240 03**: Ethics (In-Person)
    - Schedule: MWF 9:15 AM - 10:20 AM
    - Instructor: Vida Pavesich (vpavesich@usfca.edu)

13. **PHIL 240 05**: Ethics (In-Person)
    - Schedule: MW 6:30 PM - 8:15 PM
    - Instructor: Greig Mulberry (grmulberry@usfca.edu)

14. **PHIL 240 04**: Ethics (In-Person)
    - Schedule: MW 4:45 PM - 6:25 PM
    - Instructor: Greig Mulberry (grmulberry@usfca.edu)

15. **PHIL 244 02**: Environmental Ethics (In-Person)
    - Schedule: MWF 9:15 AM - 10:20 AM
    - Instructor: Stephen Friesen (smfriesen@usfca.edu)

16. **PHIL 480 01**: Topics in Contemporary Philosophy (In-Person)
    - Schedule: MWF 10:30 AM - 11:35 AM

17. **PHIL 319 01**: Logic (In-Person)
    - Schedule: MWF 11:45 AM - 12:50 PM
    - Instructor: Nick Leonard (nleonard@usfca.edu) `,
		},
		{
			name:     "TestBio",
			question: "Where does BIOL 422 and BIOL 423 meet?",
			expected: `
BIOL 422 (Bioinformatics) meets in KA 311 on Mondays and Wednesdays from 09:00 to 10:15. BIOL 423 (Bioinformatics Lab) meets in LS 205 on Mondays from 13:00 to 15:50. Therefore, BIOL 422 and BIOL 423 do not meet at the same time or location.  
			`,
		},
		{
			name:     "TestGuitar",
			question: "Can I learn guitar this semester?",
			expected: `
Yes, you can learn guitar this semester. Here are the available courses for guitar lessons:
        
1. **Course Title:** Guitar and Bass Lessons (Section 01)
- **Course Code:** MUS 121
- **Format:** In-Person
- **Dates:** 8/20/24 - 11/28/24
- **Time:** TBA
- **Instructor:** Christopher Ruscoe
- **Email:** cgruscoe@usfca.edu
- **Credits:** 5

2. **Course Title:** Guitar and Bass Lessons (Section 02)
- **Course Code:** MUS 121
- **Format:** In-Person 
- **Dates:** 8/20/24 - 11/28/24
- **Time:** TBA
- **Instructor:** Christopher Ruscoe
- **Email:** cgruscoe@usfca.edu
- **Credits:** 2
			`,
		},
		{
			name:     "TestMultiple",
			question: "I would like to take a Rhetoric course from Phil Choong. What can I take?",
			expected: `
			You can take the following Rhetoric courses from Philip Choong:

1. **RHET 103: Public Speaking**
- **Class Number:** 40166
- **Schedule:** MWF 11:45 AM - 12:50 PM
- **Location:** LM 346A
- **Dates:** 8/20/24 - 12/4/24

2. **LARHET 103: Public Speaking**
- **Class Number:** 40146
- **Schedule:** MWF 10:30 AM - 11:35 AM
- **Location:** LM 346A
- **Dates:** 8/20/24 - 12/4/24

3. **LARHET 328: Speaking Center Internship**
- **Class Number:** 42533
- **Schedule:** T 4:35 PM - 6:25 PM
- **Location:** LM 345
- **Dates:** 8/20/24 - 12/3/24

4. **LARHET 195: FYS: Podcasts: Eloquentia & Aud**
- **Class Number:** 40215
- **Schedule:** MWF 2:15 PM - 3:20 PM
- **Location:** LM 352
- **Dates:** 8/20/24 - 12/4/24
			`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			file, err := os.Open("fallclasses.csv")
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			defer file.Close()

			db, err := NewChromaDB(false)
			if err != nil {
				fmt.Printf("couldn't create new chromaDB: %v", err)
			}
			db.insertRecords(file, false)

			client := openai.NewClient(os.Getenv("api"))

			result, err := db.Query(test.question)
			if err != nil {
				return
			}

			var allQuery string
			for _, row := range result[0] {
				allQuery = allQuery + row
			}

			bot := &chatBot{
				question: test.question,
				data:     "use the following data to answer the question, keep the courses at their fullnames, no shortnames" + allQuery,
			}

			answer, err := bot.callAPI(client)
			if err != nil {
				fmt.Printf("api call didn't work: %v", err)
			}

			bot = &chatBot{
				data:     "Compare these two course listings and return ONLY 'yes' if the information is atleast 50 percent similar. The only thing that matters are the courses, instructors and meeting times.(ignoring formatting), or 'no' if they differ:\n\nFirst listing:\n" + answer + "\n\nSecond listing:\n" + test.expected,
				question: "Are these course listings equivalent in content? Answer only 'yes' or 'no'.",
			}

			actual, err := bot.callAPI(client)
			if err != nil {
				fmt.Printf("api call didn't work: %v", err)
			}

			if actual != "yes" {
				t.Errorf("Results are not similar\nGot:\n%s\n\nExpected:\n%s", answer, test.expected)

				fmt.Printf("Comparison response: %s\n", actual)
			}

		})
	}
}
