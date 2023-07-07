/*
Exercise #1 of Gopher Exercises

Input : Read in a quiz provided via CSV file

Main packages used :
OS package -> file IO
FLAGS package -> customizations.

Channels -> coroutines and concurrency

Avoid writing own CSV parser -> time taken there.
	Handle commas in the CSV too! : may have questions with commas
	Parse out final comma only.

https://pkg.go.dev/encoding/csv

*/

package main

// Get command line flags working too
import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mattn/go-tty"
)

var csvFileName string
var numberOfSeconds int

func main() {
	numQuestions := 0
	numQuestionsRight := 0
	// csvDelimeter := ","

	// Step 1 : Open the file and read in the file
	// Single string flag with default value.
	// The heck this is a string pointer : not a string.
	// Gaaah seting up command line utilities aint easy

	// Flags are always reference set? Saw this in Google too
	// Good engineers set default values for their flags, and good descripts too.
	flag.StringVar(&csvFileName, "csv", "problems.csv",
		"a csv file in the format of 'question,answer' [default \"problems.csv\"]")

	flag.IntVar(&numberOfSeconds, "limit", 30,
		"the time limit for the quiz in seconds ( default 30 )")

	// Parse flags from cmdline
	// After all flags defined : before all flags used in program -> call
	flag.Parse()
	f, err := os.Open(csvFileName)
	if err != nil {
		log.Fatal(err)
		// log fatal errors if we get an error
	}

	// defer close to end of currently executing coroutine.
	defer f.Close()

	// Pass in file to csv reader
	csvReader := csv.NewReader(f)
	// Buffered IO operations : accumulate data into buffer
	// Reader number of system calls.

	// User asked to press ENTER / other key before start of timer
	// Read in key press event.
	// TTY utilities.
	// tty, ttyErr := tty.Open()
	// if ttyErr != nil {
	// 	log.Fatal(ttyErr)
	// }
	// // defer a uniquely go keyword
	// defer tty.Close()
	fmt.Printf("User : please press the ENTER key ( or any other key ) to commence the timer\n")
	hitKey := false
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	// defer tty.Close()

	for {
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Key press => " + string(r))
			hitKey = true
			break
		}
	}
	// Forget to close this : rest of program STDIO will not work as expected :-( )
	tty.Close()

	// Go create our channel here
	// Channel for the sleeping coroutine.
	// Make bidirectional channel : send and receive data.
	ic := make(chan int)

	// Set up a signalling channel
	// Emulate a timeout in Go : indirectly - not directly
	fmt.Println(hitKey)
	if hitKey {
		go func() {
			fmt.Printf("Started timer")
			time.Sleep(time.Duration(numberOfSeconds) * time.Second)
			fmt.Printf("Timer has passed.")
			ic <- 0
			close(ic)
		}()
	}

	// User IO till crash : infinite go routine usage?
	// Break on a label to exit out of for-select-case concurrency-channel code WTF???
labelChannel:
	for {
		select {
		case <-ic:
			// numQuestionsWrong := 0
			// How to set a default CSV ( before flag customization )
			// Customize the file name via a FLAG desired
			fmt.Printf("Total number of questions correct = %d.", numQuestionsRight)
			fmt.Printf("Total questions = %d.", numQuestions)
			break labelChannel
		default:
			// Compilers telling declared BUT unusued vars ( powerful feature )
			rec, err := csvReader.Read()
			// IO EOF error huh?
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			// Prefer buffered IO to non-buffered IO for perf reasons.
			question := rec[0]
			answer, _ := strconv.Atoi(rec[1])
			fmt.Printf("What is the solution to %+v\n", question)
			userReader := bufio.NewReader(os.Stdin)
			// Delimeter is being read here? Handle that
			solutionStr, solutionErr := userReader.ReadString('\n')
			if solutionErr != nil {
				log.Fatal(solutionErr)
			}
			// Single valeu context issues?
			solution, solutionErr := strconv.Atoi(solutionStr[:len(solutionStr)-2])
			if solutionErr != nil {
				fmt.Printf("Error reading answer\n")
			}
			if solution == answer {
				numQuestionsRight++
				numQuestions++
			}
		}
	}
}

// Seperate timing function
// Measure and display time
// Pause current goroutine
// GAAAAH why cast around with a duration : direct const multiplier allowed.
// func startTimer(ch chan int) {
// 	// d time.Duration arg gaaah
// 	time.Sleep(time.Duration(numberOfSeconds) * time.Second)
// 	ch <- 0
// 	// good practice to have channels closed
// 	// after all values are sent
// 	close(ch)
// }
