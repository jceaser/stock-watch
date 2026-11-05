//***********************************************************************************************100
/*
Functions to build reports
*/

package app

import (
	//standard
	"fmt"
	"strings"

	//personal
	"github.com/jceaser/stock-watch/common"
)

// ***************************************************************************80

//used in the alert text
const(
	marker_above = "Above"
	marker_below = "Below"
)

//terminal colors
const (
	black = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

//terminal specific letters
const (
	RuneVLine    = '│'
	RuneHLine    = '─'
    RunePlus     = '┼'
    RuneTTee     = '┬'
    RuneRTee     = '┤'
    RuneLTee     = '├'
    RuneBTee     = '┴'
    RuneULCorner = '┌'
    RuneURCorner = '┐'
)

// ***************************************************************************80
//MARK - functions

func color(color int, text string) string {
	return fmt.Sprintf("\033[0;%dm%s\033[0m", color, text)
}

func negativeColor(num float64, neg, pos int) int {
	if num < 0.0 {
		return neg
	}
	return pos
}

func aboveColor(text, mark string, markColor, normalColor int) int {
	if strings.Contains(text, mark) {
		return markColor
	}
	return normalColor
}

// ***************************************************************************80
//MARK - Reports

/* Create a report based on the current data, then hand the display of the report to a handler */
func DoReportFormat(state any, command Command, flag string, options ...string) string {
	count := len(options)
	if count != 0 && (flag == "help" || HelpRequested(options...) || 1 < count) {
		return fmt.Sprintf("%s mode", command.Word)
	}

	modes := []string{}
	if len(options) > 0 {
		modes = strings.Split(options[0], ",")
	}

	if appState, okay := state.(AppState) ; okay {
		// shortcuts to make the new functions of NumStr
		n := common.NewNumStrNumber
		s := common.NewNumStrString

		header := []common.NumStr{s("SYM"), s("52w Low"), s("Price"), s("52w High"), s("Cost"),
			s("Value"), s("Avg(c)"), s("Profit"), s("Alert"), s("Note")}

		//create a 2D array of NumStr values, values that can be number or strings
		table := [][]common.NumStr{}
		table = append(table, header)

		//loop through the data and build each row
		data := DataLoad(appState.DataFile)
		for _, stock := range data.Stocks {
			//create special for alert text
			alert := ""
			if stock.Price < stock.Alert.Low {
				alert = fmt.Sprintf("%s %.2f", marker_below, stock.Alert.Low)
			}
			if stock.Price > stock.Alert.High {
				alert = fmt.Sprintf("%s %.2f", marker_above, stock.Alert.High)
			}

			//now that we have the special values, get out of requested
			if 	common.Contains(modes, "alerts") && len(alert) < 1 {
				continue
			}
			if common.Contains(modes, "invitory") && len(stock.Inventory) < 1 {
				continue
			}

			//assemble the row
			row := []common.NumStr{s(stock.ID),
				n(stock.Week52.Low),
				n(stock.Price),
				n(stock.Week52.High),
				n(stock.Cost()),
				n(stock.Value()),
				n(stock.AverageItemCost()),
				n(stock.Profit()),
				s(alert),
				s(stock.Note)}
			table = append(table, row)
		}

		//send the table to the requested formater
		var buffer strings.Builder
		if flag == "csv" {
			buffer.WriteString(reportDeliminated(table, ","))
		} else if flag == "tab" {
			buffer.WriteString(reportDeliminated(table, "\t"))
		} else if flag == "markdown" {
			buffer.WriteString(reportMarkDown(data.Name, table))
		} else if flag == "shell" {
			buffer.WriteString(reportShell(data.Name, table))
		} else {
			for i, r := range table {
				buffer.WriteString(fmt.Sprintf("%d) %v\n", i, r))
			}
		}
		return buffer.String()
	}
	return "" // no data, no report
}

/* Generate the report as a table formated as a CSV export */
func reportDeliminated(table [][]common.NumStr, deliminator string) string {
	var buffer strings.Builder

	for _, row := range table {
		for x, cell := range row {
			if x > 0 {
				buffer.WriteString(deliminator)
			}
			if cell.HasNumber() {
				buffer.WriteString(fmt.Sprintf("%.2f", cell.GetNumber()))
			} else {
				text := cell.GetString()
				if strings.Contains(text, deliminator) {
					buffer.WriteString(fmt.Sprintf("\"%v\"", cell))
				} else {
					buffer.WriteString(fmt.Sprintf("%v", cell))
				}
			}
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

/* Generate the report as a table formated as markdown */
func reportMarkDown(title string, table [][]common.NumStr) string {
	var buffer strings.Builder

	head := "%5s | %7s | %7s | %8s | %7s | %7s | %7s | %7s | %-13s | %s\n"
	format := "%5s | %7.2f | %7.2f | %8.2f | %7.2f | %7.2f | %7.2f | %7.2f | %-13s | %s\n"

	buffer.WriteString(fmt.Sprintf("## %s\n\n", title))
	for y, row := range table {
		if y==0 {
			line := fmt.Sprintf(head,
				row[0].GetString(),
				row[1].GetString(),
				row[2].GetString(),
				row[3].GetString(),
				row[4].GetString(),
				row[5].GetString(),
				row[6].GetString(),
				row[7].GetString(),
				row[8].GetString(),
				row[9].GetString())
			buffer.WriteString(line)
			line = fmt.Sprintf(head,
				strings.Repeat("-", 4),
				strings.Repeat("-", 7),
				strings.Repeat("-", 7),
				strings.Repeat("-", 8),
				strings.Repeat("-", 7),
				strings.Repeat("-", 7),
				strings.Repeat("-", 7),
				strings.Repeat("-", 7),
				strings.Repeat("-", 13),
				strings.Repeat("-", 13))
			buffer.WriteString(line)

		} else {
			line := fmt.Sprintf(format,
				row[0].GetString(),
				row[1].GetNumber(),
				row[2].GetNumber(),
				row[3].GetNumber(),
				row[4].GetNumber(),
				row[5].GetNumber(),
				row[6].GetNumber(),
				row[7].GetNumber(),
				row[8].GetString(),
				row[9].GetString())
			buffer.WriteString(line)
		}
	}
	return buffer.String()
}

/* Generate the report as a table formated and styles for a terminal */
func reportShell(title string, table [][]common.NumStr) string {
	var buffer strings.Builder

	head := " %5s │ %8s │ %8s │ %8s │ %8s │ %8s │ %8s │ %9s │ %-13s │ %s\n"
	sep := " %5s┼%8s┼%8s┼%8s┼%8s┼%8s┼%8s┼%9s┼%-13s┼%s\n"
	format := " %5s │ %8.2f │ %8.2f │ %8.2f │ %8.2f │ %8.2f │ %8.2f │ \033[0;%dm%9.2f\033[0m │ \033[0;%dm%-13s\033[0m │ %s\n"

	buffer.WriteString(fmt.Sprintf("\033[0;4m%s\033[0m\n\n", title))
	for y, row := range table {
		if y==0 {
			line := fmt.Sprintf(head,
				row[0].GetString(),
				row[1].GetString(),
				row[2].GetString(),
				row[3].GetString(),
				row[4].GetString(),
				row[5].GetString(),
				row[6].GetString(),
				row[7].GetString(),
				row[8].GetString(),
				row[9].GetString())
			buffer.WriteString(line)
			line = fmt.Sprintf(sep,
				strings.Repeat("─", 6),
				strings.Repeat("─", 10),
				strings.Repeat("─", 10),
				strings.Repeat("─", 10),
				strings.Repeat("─", 10),
				strings.Repeat("─", 10),
				strings.Repeat("─", 10),
				strings.Repeat("─", 11),
				strings.Repeat("─", 15),
				strings.Repeat("─", 13))
			buffer.WriteString(line)

		} else {
			profit := row[7].GetNumber()
			alert := row[8].GetString()
			line := fmt.Sprintf(format,
				row[0].GetString(),
				row[1].GetNumber(),
				row[2].GetNumber(),
				row[3].GetNumber(),
				row[4].GetNumber(),
				row[5].GetNumber(),
				row[6].GetNumber(),
				negativeColor(profit, red, white), profit,
				aboveColor(alert, marker_above, green, magenta), alert,
				row[9].GetString())
			buffer.WriteString(line)
		}
	}
	return buffer.String()
}
