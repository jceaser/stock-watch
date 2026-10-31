//***********************************************************************************************100
/*
Interactive commands that act on the json data model. An function that starts with Do is assumed to
be a command which can be called by the interface.

Commands have the following features:
* trigger word
* Function
* optional parameter - one command can do multiple things
* Help Message

Take care in naming thetrigger word in that the selector will find the first action that matches
what the user types in, partial words are allowed

*/

package app

import (
	//standard
	"fmt"
	"strconv"
	"strings"

	//personal
	"github.com/jceaser/stock-watch/common"
	"github.com/jceaser/stock-watch/stocks"
)

//****************************************************************************80

func init() {
	Cmd = Commands{List: map[string]Command{}}
	Cmd.A("help", DoHelp, "", "Ask for help")
	Cmd.A("echo", DoEcho, "", "Repeat yourself, Repeat yourself")

	// CRUD
	Cmd.A("create", DoCreate, "", "Create a stock record")
	Cmd.A("read", DoRead, "", "Read all stock records")
	Cmd.A("update", DoUpdate, "", "Update a stock details")
	Cmd.A("delete", DoDelete, "", "remove stock from record")

	Cmd.A("alert", DoAlertSet, "", "Set High and Low Alerts for a stock")

	Cmd.A("fetch", DoFetch, "", "Fetch latest prices for all stacks")
	Cmd.A("assume-alerts", DoLowHighSet, "", "Make an assumption on a good alert low/high value")

	Cmd.A("report-full", DoReport, "", "Run report on stocks")
	Cmd.A("report-csv", DoReportFormat, "csv", "Run report on stocks")
	Cmd.A("report-markdown", DoReportFormat, "markdown", "Run report on stocks")
	Cmd.A("report-shell", DoReportFormat, "shell", "Run report on stocks")

	//Cmd.A("add-price", DoNoImpl, "", "add a new price to the list")

	// Stock innovatory CRuD
	Cmd.A("stock-create", DoCreateStock, "", "add stocks share count and price")
	Cmd.A("stock-read", DoReadStock, "", "show the stock count and price list")
	Cmd.A("stock-update", DoUpdateStock, "", "update stock count and price")
	Cmd.A("stock-delete", DoDeleteStock, "", "delete stocks shares")
}

//****************************************************************************80
//MARK - Basic CRUD

/* create a new stock and save it to the database */
func DoCreate(state any, command Command, flag string, options ...string) string {
	count := len(options)
	if HelpRequested(options...) || 3 < count {
		return fmt.Sprintf("%s id description? note?", command.Word)
	}

	stockId := strings.ToUpper(options[0])

	// 1. check if stock already exists
	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		if found, exists := data.Exists(stockId) ; exists {
			return fmt.Sprintf("Stock %s (%s) already exists.", found.Name, found.ID)
		}
 	}

	// 2. check if id exists
	liveData, err := FetchData(stockId)
	if err != nil {
		return fmt.Sprintf("Could not find stock: %v", err)
	}

	// 3. create stock
	stock := stocks.StockItem{}
	stock.ID = stockId
	stock.Name = liveData.ShortName
	if count > 1 {
		stock.Description = options[1]
	}
	if count > 2 {
		stock.Note = options[2]
	}

	// 4. make an assumed high/low alert
	stock.Alert.Low, stock.Alert.High = AssumeHighLow(liveData)

	// 5. save stock data
	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		data.Stocks = append(data.Stocks, stock)
		DataStore(appState.DataFile, data)
	}
	return ""
}

/* read the data file and print it out */
func DoRead(state any, command Command, flag string, options ...string) string {
	const ALL = ""
	count := len(options)
	//allow now parameters, the help parameter, but not more then 1
	if count!=0 && HelpRequested(options...) || 1 < len(options) {
		return fmt.Sprintf("%s id?", command.Word)
	}

	id := ALL
	if len(options) > 0 {
		id = options[0]
	}

	var buffer strings.Builder
	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		for i, stock := range data.Stocks {
			if id == ALL{
				buffer.WriteString(fmt.Sprintf("%d: %s\n", i, stock))
			} else if stock.ID == id {
				buffer.WriteString(fmt.Sprintf("%s\n", stock))
				break
			}
		}
	}
	return fmt.Sprintf("%s", buffer.String())
}

/* read the data file and print it out */
func DoUpdate(state any, command Command, flag string, options ...string) string {
	count := len(options)
	if HelpRequested(options...) || count < 2 || 3 < count {
		return fmt.Sprintf("%s id description note(optional)", command.Word)
	}
	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		for i, stock := range data.Stocks {
			if stock.ID == options[0] {
				if options[1] != "" {
					data.Stocks[i].Description = options[1]
				}
				if count > 2 && options[2] != "" {
					data.Stocks[i].Note = options[2]
				}
			}
		}
		DataStore(appState.DataFile, data)
	}
	return ""
}

/* create a new stock and save it to the database */
func DoDelete(state any, command Command, flag string, options ...string) string {
	if HelpRequested(options...) || len(options) != 1 {
		return fmt.Sprintf("%s id", command.Word)
	}
	ret := ""
	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		for si, stock := range data.Stocks {
			ret = stock.String()
			if stock.ID == options[0] {
				data.Stocks = append(data.Stocks[:si], data.Stocks[si+1:]...)
				break
			}
		}

		DataStore(appState.DataFile, data)
	}
	return ret
}

// **********************************40
//MARK - CRuD Stock innovatory

/* add innovatory, but not negative innovatory */
func DoCreateStock(state any, command Command, flag string, options ...string) string {
	count := len(options)
	if HelpRequested(options...) || count < 2 || count > 3 {
		return fmt.Sprintf("%s id count value", command.Word)
	}

	id := options[0]

	quantity, err := strconv.ParseFloat(options[1], 64)
	if err != nil {
		 fmt.Println("Error:", err)
	}
	if quantity <= 0.0 {
		return "Must have a positive number of stocks"
	}

	value, err :=  strconv.ParseFloat(options[2], 64)
	if err != nil {
 		fmt.Println("Error:", err)
 	}
	if value < 0.0 {
		return "Must have a positive value"
	}

	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		for y, row := range data.Stocks {
			if row.ID == id {
				data.Stocks[y].AddInventory(quantity, value)
				break
			}
		}
		DataStore(appState.DataFile, data)
	}

	return ""
}

/* List out the current stock innovatory */
func DoReadStock(state any, command Command, flag string, options ...string) string {
	if HelpRequested(options...) || len(options) != 1 {
		return fmt.Sprintf("%s id", command.Word)
	}
	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		for _, stock := range data.Stocks {
			if stock.ID == options[0] {
				var buffer strings.Builder
				for y, row := range stock.Inventory {
					buffer.WriteString(fmt.Sprintf("%d: %v\n", y, row))
				}
				return buffer.String()
			}
		}
	}
	return ""
}

/* Update a specific stock innovatory, returning the old values */
func DoUpdateStock(state any, command Command, flag string, options ...string) string {
	count := len(options)
	if HelpRequested(options...) || count < 3 || count > 4 {
		return fmt.Sprintf("%s id row count value(?)", command.Word)
	}

	// 0: get stock id
	stockId := strings.ToUpper(strings.TrimSpace(options[0]))
	if len(stockId) < 1 {
		return "Must have a stock id: " + options[0]
	}

	// 1: get index to update
	index, err :=  strconv.Atoi(options[1])
	if err != nil {
		 fmt.Println("Error:", err)
	}
	if index < 0 {
		return "Must have a positive row number: " + options[1]
	}

	// 2: get quantity, mutual funds allow partial ownership, thus a float
	quantity, err := strconv.ParseFloat(options[2], 64)
	if err != nil {
		 fmt.Println("Error:", err)
	}
	if quantity < 0.0 || quantity != -1.0 {
		return "Must have a positive quantity number or be -1.0 : " + options[2]
	}

	// 3: get the stock value
	value := 0.0
	if 3 < count {
		value, err = strconv.ParseFloat(options[3], 64)
		if err != nil {
			fmt.Println("Error:", err)
		}
		if value < 0.0 || value != -1.0 {
			return "Must have a positive value number or be -1.0: " + options[3]
		}
	}

	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		for si, stock := range data.Stocks {
			if stock.ID == stockId {
				ret := ""
				if len(stock.Inventory) > index {
					if quantity > 0.0 {
						ret += fmt.Sprintf("Count: %.3f", data.Stocks[si].Inventory[index].Count)
						data.Stocks[si].Inventory[index].Count = quantity
					}
					if quantity > 0.0 && value > 0.0 {
						ret += " "
					}
					if value > 0.0 {
						ret += fmt.Sprintf("Value: %.3f", data.Stocks[si].Inventory[index].Value)
						data.Stocks[si].Inventory[index].Value = value
					}
					DataStore(appState.DataFile, data)
					return ret
				}
			}
		}
	}
	return ""
}

func DoDeleteStock(state any, command Command, flag string, options ...string) string {
	if HelpRequested(options...) || len(options) != 2 {
		return fmt.Sprintf("%s id row", command.Word)
	}

	index, err :=  strconv.Atoi(options[1])
	if err != nil {
		 fmt.Println("Error:", err)
	}
	if index < 1 {
		return "Must have a positive row number"
	}

	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		for si, stock := range data.Stocks {
			if stock.ID == options[0] {
				if len(stock.Inventory) > index {
					dropped := stock.Inventory[index]
					data.Stocks[si].Inventory = append(data.Stocks[si].Inventory[:index],
						data.Stocks[si].Inventory[index+1:]...)
					DataStore(appState.DataFile, data)
					return fmt.Sprintf("%v", dropped)
				}
			}
		}
	}

	return ""
}

// ***************************************************************************80
//MARK - Other actions

func DoAlertSet(state any, command Command, flag string, options ...string) string {
	if HelpRequested(options...) || len(options) != 3 {
		return fmt.Sprintf("%s id low high", command.Word)
	}

	low, err := strconv.ParseFloat(options[1], 64)
	if err != nil {
		 fmt.Println("Error:", err)
	}
	if low <= 0.0 {
		return "Must have a positive row number"
	}

	high, err :=  strconv.ParseFloat(options[2], 64)
	if err != nil {
		 fmt.Println("Error:", err)
	}
	if high <= 0.0 {
		return "Must have a positive value"
	}

	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		for si, stock := range data.Stocks {
			if stock.ID == options[0] {
				data.Stocks[si].Alert.Low = low
				data.Stocks[si].Alert.High = high
				DataStore(appState.DataFile, data)
				break
			}
		}
	}
	return ""
}

func DoFetch(state any, command Command, flag string, options ...string) string {
	var buffer strings.Builder
	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		for si, row := range data.Stocks {
			stockData, err := FetchData(row.ID)
			if err != nil {
				buffer.WriteString(fmt.Sprintf("%v", err))
				continue
			}

			data.Stocks[si].Name = stockData.ShortName
			data.Stocks[si].Week52.High = stockData.FiftyTwoWeekHigh
			data.Stocks[si].Price = stockData.MarketPrice
			data.Stocks[si].Week52.Low = stockData.FiftyTwoWeekLow

			when := stockData.DateTime
			data.Stocks[si].AddPrice(when, stockData.MarketPrice)
			DataStore(appState.DataFile, data)

			buffer.WriteString(fmt.Sprintf("%s - %f < %f < %f\n",
				data.Stocks[si].Name,
				data.Stocks[si].Week52.Low,
				data.Stocks[si].Price,
				data.Stocks[si].Week52.High))
		}
	}
	return buffer.String()
}

func DoLowHighSet(state any, command Command, flag string, options ...string) string {
	if HelpRequested(options...) || len(options) < 1 {
		return fmt.Sprintf("%s id", command.Word)
	}
	id := options[0]

	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		for si, row := range data.Stocks {
			if row.ID == id {
				stockData, err := FetchData(id)
				if err != nil {
					return fmt.Sprintf("%v", err)
				}
				data.Stocks[si].Alert.Low, data.Stocks[si].Alert.High = AssumeHighLow(stockData)
				DataStore(appState.DataFile, data)
				break
			}
		}
	}
	return ""
}

func DoReport(state any, command Command, flag string, options ...string) string {
	count := len(options)
	if count != 0 && (flag == "help" || HelpRequested(options...) || 1 < count) {
		return fmt.Sprintf("%s mode", command.Word)
	}

	modes := []string{}
	if len(options) > 0 {
		modes = strings.Split(options[0], ",")
	}

	head := "%5s | %7s | %7s | %8s | %7s | %7s | %7s | %7s | %-13s | %s\n"
	format := "%5s | %7.2f | %7.2f | %8.2f | %7.2f | %7.2f | %7.2f | %7.2f | %-13s | %s\n"

	var buffer strings.Builder
	if appState, okay := state.(AppState) ; okay {
		data := DataLoad(appState.DataFile)
		buffer.WriteString(fmt.Sprintf("\nStock List: %s - %s\n", data.Name, data.Description))
		buffer.WriteString("-------------------------------------------------------------\n\n")
		h1 := fmt.Sprintf(head, "SYM",
			"52w Low", "Price", "52w High", "Cost", "Value", "Avg(c)", "Profit", "Alert", "Note")
		h2 := fmt.Sprintf(head, "-----",
			"-------", "-------", "--------", "-------", "-------", "-------", "-------",
			"-------------", "----")
		buffer.WriteString(h1)
		buffer.WriteString(h2)

		for _, stock := range data.Stocks {
			alert := ""
			if stock.Price < stock.Alert.Low {
				alert = fmt.Sprintf("Below %-7.2f", stock.Alert.Low)
			}
			if stock.Price > stock.Alert.High {
				alert = fmt.Sprintf("\033[0;32mAbove %-7.2f\033[0m", stock.Alert.High)
			}
			if 	common.Contains(modes, "alerts") && len(alert) < 1 {
				continue
			}
			if common.Contains(modes, "invitory") && len(stock.Inventory) < 1 {
				continue
			}
			out := fmt.Sprintf(format,
				stock.ID,

				stock.Week52.Low,
				stock.Price,
				stock.Week52.High,

				stock.Cost(),
				stock.Value(),
				stock.AverageItemCost(),
				stock.Profit(),
				alert,
				stock.Note)
			buffer.WriteString(out)
		}
	}

	return buffer.String()
}
