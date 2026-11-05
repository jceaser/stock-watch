//***********************************************************************************************100
/*
A command line too for managing lists of stocks, fetching current stock details and generating
reports. To use the tool, call the CRUD interfaces to Create stocks of interest. By default the tool
will add alerts near the 52 week high and low values. Then call "fetch" to grab current price data.
From that point you can call "report" to see which stocks are outside of the alert bounds.

The code works by exposing a set of "commands" to the interactive interface. These commands do
simple tasks against the JSON data model. Command can be found in app/actions.go

Stocks contain
* Identifying information like Ticker Name Symbol, Name, Description
* Alert Thresholds
* Historical Price List
* Inventory of stocks
* Current price values

Code Layout
* app - application code not related to being a command line tool
* common - functions used by me in many projects, should not depend on any other libraries
* services - encapsulation of third party stock data APIs
* stocks - data model of stocks

Run with:
	go run . -file test.json -batch-mode -initial report 2> /dev/null
	go build .
*/

package main

import (
	//standard
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	//personal
	"github.com/jceaser/stock-watch/app"

	// third party
	"github.com/peterh/liner"
)

//****************************************************************************80
//MARK - Interactive function

/* configure the interactive liner */
func setup_liner(line *liner.State) string {
	history_fn := filepath.Join(os.TempDir(), ".stock-watch_history") 	//used by liner
	names := []string{"print", "dump", "quit", "exit"}            		//used by liner

	line.SetCtrlCAborts(true)

	line.SetTabCompletionStyle(liner.TabPrints)
	line.SetCompleter(func(line string) (c []string) {
		for _, n := range names {
			fmt.Print(n)
			if strings.HasPrefix(n, strings.ToLower(line)) {
				c = append(c, n)
				fmt.Print(n)
			}
		}
		return
	})
	if f, err := os.Open(history_fn); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
	return history_fn
}

/* Enter into interactive mode */
func InteractiveMode(appState app.AppState) {
	//readline setup
	line := liner.NewLiner()
	defer line.Close()
	history := setup_liner(line)

	fmt.Printf("------------------------------------------------------------\n")
	ShowVersion()

	for {
		if name, err := line.Prompt(">"); err == nil {
			input := strings.Trim(name, " ")
			if len(input) < 1 {
				continue
			}
			if strings.HasPrefix("quit", input) || strings.HasPrefix("exit", input) {
				break
			}
			line.AppendHistory(name)    //save it
			ProcessLine(appState, input) //use it
		} else if err == liner.ErrPromptAborted {
			fmt.Print("Aborted")
		} else {
			fmt.Print("Error reading line: ", err)
		}
		//save the history
		if f, err := os.Create(history); err != nil {
			fmt.Print("Error creating history file: ", err)
		} else {
			line.WriteHistory(f)
			f.Close()
		}
	}
}

/* Act on an entire line of a command, which may contain multiple commands */
func ProcessLine(appState app.AppState, input string) {
	commands := strings.Split(input, ";")
	for _, cmd := range commands {
		cmd = strings.Trim(cmd, " ")
		ProcessCommand(appState, cmd)
	}
}

/* act on one command from a line of commands */
func ProcessCommand(appState app.AppState, line string) {
	response := app.Cmd.Run(appState, line)
	if len(response) > 0 {
		fmt.Printf("%s\n", response)
	}
}

//****************************************************************************80
//MARK - App functions

/*
Print out the application version information, this may not be the same as the data format version
*/
func ShowVersion() {
	fmt.Printf("Stock Watch %s by thomas.cherry@gmail.com\n", app.Version())
}

/* Handle command line paramaters */
func initParams() app.AppState {
	state := app.AppState{}

	flag.BoolVar(&state.BatchMode, "batch", false, "batch mode vs interactive")
	flag.StringVar(&state.DataFile, "file", "stocks.json", "data file")
	flag.StringVar(&state.Initial, "initial", "", "Initial commands")

	//state.LogLevel = flag.Int("log-level", 1, "set log level")
	flag.BoolVar(&state.Version, "version", false, "Show version info")

	flag.Parse()
	return state
}

//****************************************************************************80

func main() {
	appData := initParams()

	if appData.Version {
		ShowVersion()
		os.Exit(0)
	}

	//consume Initial setup instructions
	if len(appData.Initial) > 0 {
		ProcessLine(appData, appData.Initial)
		appData.Initial = ""
	}

	//process stream if it exists, or enter interactive mode
    stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		formula := ""
		for scanner.Scan() {
			token := scanner.Text()
			formula = fmt.Sprintf("%s %s", formula, token)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading input:", err)
		}
		formula = strings.Trim(formula, " ")
		if 0 < len(formula) {
			ProcessLine(appData, formula)
			//appData.BatchMode = true // you can not continue after reading from stdin
		}
	} else {
		if !appData.BatchMode {
			InteractiveMode(appData)
		}
	}

}

