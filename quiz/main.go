package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

type problem struct {
	question string
	answer   string
}

func (p problem) String() string {
	return fmt.Sprintf("%s = ", p.question)
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "CSV file with 'question,answer' format;")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Failed to open file: %s", *csvFilename)
		os.Exit(1)
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Failed to read file: %s", *csvFilename)
		os.Exit(1)
	}

	correct := 0
	for index, problem := range parseLines(lines) {
		msg := fmt.Sprintf("Problem #%d: %s\n", index+1, problem)
		fmt.Println(msg)

		var answer string
		fmt.Scanf("%s\n", &answer)

		if answer == problem.answer {
			correct++
		}
	}

	fmt.Printf("Scored %d out of %d\n", correct, len(lines))
}

func parseLines(lines [][]string) (result []problem) {
	result = make([]problem, len(lines))

	for index, line := range lines {
		result[index] = problem{
			question: line[0],
			answer:   line[1],
		}
	}

	return
}
