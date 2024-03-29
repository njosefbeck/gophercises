package main

import (
	"flag"
	"os"
	"fmt"
	"encoding/csv"
	"strings"
	"time"
)

func main() {
	// csvFilename is a pointer to a string
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()
	
	// we use an asterick here because we want to use the VALUE
	// of csvFilename and not the pointer -- TO DO: research pointers
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	
	// Reads the CSV, implements an I/O Reader interface
	// TO DO: read up on interfaces and I/O Reader
	r := csv.NewReader(file)

	// Reads all the lines in the CSV
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)

		// Go channel
		answerCh := make(chan string)

		// Goroutine
		go func() {
			var answer string
			// Using a pointer value here so whenever Scanf sets the value
			// we can then access it with our var
			fmt.Scanf("%s\n", &answer)

			// Send the answer into the answer channel (happens when we get a response)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <- answerCh:
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	// Since we know how many lines are in the file,
	// it's recommended here to use len(), so that Go doesn't have to do the
	// extra work of sizing the slice
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
