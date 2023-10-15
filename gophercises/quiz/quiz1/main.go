package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
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

	flag.Parse()

	// Read the CSV file
	file, err := os.Open(*csvfile)
	if err != nil {
		fmt.Println("Failed to open the CSV file: ", *csvfile, "\n", err.Error())
		os.Exit(1)
	}

	csvreader := csv.NewReader(file)

	var record []string
	var count int
	// Read the CSV file line by line
	for record, err = csvreader.Read(); err == nil; record, err = csvreader.Read() {

		// Parse the record in the format of 'question,answer'
		problem , err := parseRecord(record)
		if err != nil {
			fmt.Println("Failed to parse the record: ", record, "\n", err.Error())
			os.Exit(1)
		}
		
		// Ask the question and check answer
		fmt.Printf("Problem: %s = ", problem.question)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		answer := scanner.Text()

		if answer == problem.answer {
			count++
		}

	}
	if err.Error() != "EOF"{
		fmt.Println("Failed to read the record: ", err.Error())
		os.Exit(1)
	}

	totalQuestions,_ := csvreader.FieldPos(0)
	
	fmt.Printf("You scored %d out of %d\n", count, totalQuestions)


}

func parseRecord(record []string) (problem, error){
	
	if len(record) != 2 {
		return problem{}, fmt.Errorf("Invalid record: %v", record)
	}

	return problem{
		question: record[0],
		answer: record[1],
	}, nil
}