package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

func (p problem) String() string {
	return fmt.Sprintf("%s = ", p.question)
}

var (
	csvFilename = flag.String("csv", "problems.csv", "CSV file with 'question,answer' format;")
	timeLimit   = flag.Int("limit", 30, "limit time to solve the quiz")
)

func init() {
	flag.Parse()
}

func main() {
	sliceOfProblems := loadProblems()
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

problemLoop:
	for index, problem := range sliceOfProblems {
		fmt.Printf("Problem #%d: %s\n", index+1, problem)

		answerChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			showScore(correct, len(sliceOfProblems))
			break problemLoop
		case answer := <-answerChannel:
			if answer == problem.answer {
				correct++
			}
		}
	}

	showScore(correct, len(sliceOfProblems))
}

func loadProblems() []problem {
	file := openFile()
	lines := readFileLinesFrom(file)

	return parseLinesIntoProblems(lines)
}

func openFile() (file *os.File) {
	file, err := os.Open(*csvFilename)
	defer file.Close()
	if err != nil {
		fmt.Printf("Failed to open file: %s", *csvFilename)
		os.Exit(1)
	}
	return
}

func readFileLinesFrom(file *os.File) [][]string {
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Failed to read file: %s", *csvFilename)
		os.Exit(1)
	}

	return lines
}

func parseLinesIntoProblems(lines [][]string) (result []problem) {
	result = make([]problem, len(lines))

	for index, line := range lines {
		result[index] = problem{
			question: line[0],
			answer:   line[1],
		}
	}

	return
}

func showScore(correctAmount, problemAmount int) {
	fmt.Printf("\nYou scored %d out of %d", correctAmount, problemAmount)
}
