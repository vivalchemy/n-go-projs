package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"time"
)

type problem struct {
	question string
	answer   string
}

var questionsCount int

func problemsPuller(fileName string) ([]problem, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	problems, err := parseProblems(records[1:])
	if err != nil {
		return nil, err
	}

	return problems, nil

}

// can use shuffled lines to get unique questions
func parseProblems(lines [][]string) ([]problem, error) {
	var problems []problem
	if len(lines) < questionsCount {
		return nil, fmt.Errorf("not enough questions in the file")
	}
	for _ = range questionsCount {
		randomIndex := rand.IntN(len(lines))
		if len(lines[randomIndex]) != 2 {
			return nil, fmt.Errorf("invalid line: %s", lines[randomIndex])
		}
		problems = append(problems, problem{lines[randomIndex][0], lines[randomIndex][1]})
	}
	return problems, nil
}

func main() {
	// input the name of the file
	fileName := flag.String("file", "quiz.csv", "the path of the csv file (default 'quiz.csv')")

	// set the duration of the timer
	timer := flag.Duration("t", 30*time.Second, "the duration of the timer in seconds (default 30s)")

	// number of questions to ask
	questionsCount = *flag.Int("q", 10, "the number of questions to ask (default 10)")
	flag.Parse()

	// pull the problems from the file (calling our problems puller func)
	problems, err := problemsPuller(*fileName)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}

	// counter to count the numebr of correct answers
	var correctAnswers int8

	// using the duration of the timer, we want to initialize the timer
	tObj := time.NewTimer(*timer)
	ansChan := make(chan string)

	// loop through the problems, print the questions, we'll accept the answers
problemLoop:

	for i, problem := range problems {
		var answer string
		fmt.Printf("%v. %v = ", i+1, problem.question)
		go func() {
			fmt.Scanf("%s", &answer)
			ansChan <- answer
		}()
		select {
		case <-tObj.C:
			fmt.Println("Time's up!")
			break problemLoop
		case answer := <-ansChan:
			if answer == problem.answer {
				correctAnswers++
			}
			if i == len(problems)-1 {
				close(ansChan)
			}
		}
	}
	// well will calculate and print out the result along with the time taken for the test
	fmt.Println("You got", correctAnswers, "out of", len(problems), "correct.")

	fmt.Println("Press enter to exit")
	<-ansChan
}
