package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Parse flags
	csvName := flag.String("csv", "problems.csv", "csv file containing problems")
	questionTimeout := flag.Int("question-timeout", 5, "timeout for answering a question")
	quizTimeout := flag.Int("timeout", 30, "timeout for the total quiz")
	flag.Parse()

	// Open CSV
	csvFile := openCsv(*csvName)
	defer csvFile.Close()
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	// Do the quiz and show the result
	correctAnswers := askQuestions(csvLines, quizTimeout, questionTimeout)
	fmt.Printf("You scored %d out of %d.\n", correctAnswers, len(csvLines))
}

func openCsv(fileName string) *os.File {
	csvFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	return csvFile
}

func askQuestions(csvLines [][]string, quizTimeout *int, questionTimeout *int) int {
	inputChannel := make(chan int)
	go readInput(inputChannel)
	quizChannel := time.After(time.Second * time.Duration(*quizTimeout))
	correctAnswers := 0
quiz:
	for i, line := range csvLines {
		fmt.Printf("Problem #%d: %s = ", i+1, line[0])
		answer, _ := strconv.Atoi(line[1])
		select {
		case response := <-inputChannel:
			if answer == response {
				correctAnswers++
			}
		case <-time.After(time.Second * time.Duration(*questionTimeout)):
			fmt.Println()
		case <-quizChannel:
			fmt.Println()
			break quiz
		}
	}
	return correctAnswers
}

func readInput(channel chan<- int) {
	for {
		var in string
		_, err := fmt.Scanf("%s\n", &in)
		if err != nil {
			panic(err)
		}
		in = strings.TrimSpace(in)
		a, _err := strconv.Atoi(in)
		if _err != nil {
			fmt.Println("Please enter your answer as integer!")
		} else {
			channel <- a
		}
	}
}
