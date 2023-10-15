package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

type problem struct{
	question string
	answer string
}

func main() {

	// Recover from panic
	defer func(){
		if panicCheck := recover(); panicCheck != nil {
			stacktrace := debug.Stack()
			fmt.Println("Panic occured: ", panicCheck)
			fmt.Println(string(stacktrace))
			os.Exit(1)
		}
	}()

	// Define flags
	csvfile := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	// Read the CSV file
	file, err := os.Open(*csvfile)
	if err != nil {
		fmt.Println("Failed to open the CSV file: ", *csvfile, "\n", err.Error())
		os.Exit(1)
	}

	csvreader := csv.NewReader(file)
	records, err := csvreader.ReadAll()
	if err != nil {
		fmt.Println("Failed to read the CSV file: ", *csvfile, "\n", err.Error())
		os.Exit(1)
	}

	// Parse the record in the format of 'question,answer'
	problems , err := parseRecords(records)
	if err != nil {
		fmt.Println("Failed to parse the records: ",err.Error())
		os.Exit(1)
	}
	totalQuestions := len(problems)

	answerCh := make(chan string, 1)
	score := make(chan int, 1)
	next := make(chan bool, 1)	
	fmt.Println("You have", *timeLimit, "seconds to answer", totalQuestions, "questions.")
	fmt.Println("Press enter to start the quiz.")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	scanner.Text()

	timer := time.NewTimer(time.Duration(*timeLimit) *time.Second)
	go askQuestions(problems, answerCh, next)
	go answerQuestions(answerCh, score, timer, next)


	result := <- score
	fmt.Println("You scored", result, "out of", totalQuestions)
}

func parseRecords(records [][]string) ([]problem, error){
	
	problems := make([]problem, len(records))

	for i, record := range records {
		if len(record) != 2 {
			return nil, fmt.Errorf("Invalid record: %s", record)
		}
		problems[i].question = record[0]
		problems[i].answer = record[1]
	}
	return problems, nil
}

func askQuestions(problems []problem, answerCh chan string, next chan bool){
	defer close(answerCh)
	for _, problem := range problems {
		answerCh <- problem.answer
		fmt.Printf("Problem: %s = ", problem.question)
		<-next
	}
}

func answerQuestions(answerCh chan string, score chan int, timeLimit *time.Timer, next chan bool){
	var count int
	defer close(score)
	defer close(next)
	go func(){
		for  {
			answer, ok := <- answerCh
			if !ok {
				break
			}
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			userAnswer := scanner.Text()
			if userAnswer == answer {
				count++
			}
			next <- true
		}
		score <- count
	}()

	select {
	case <- timeLimit.C:
		fmt.Println("\nTime's up!")
		score <- count		
	}
}