//***********************************************************************************************100
/*
Functions common to all of the app module
*/
package app

import (
	//standard modules
	"fmt"
	
	//local modules
	"github.com/jceaser/stock-watch/common"
	"github.com/jceaser/stock-watch/stocks"
	"github.com/jceaser/stock-watch/services"
	//note, selecting this api for active use
	"github.com/jceaser/stock-watch/services/robowealth"
)

type AppState struct {
	BatchMode bool
	DataFile string
	Initial string
	Version bool
}

// app version number, not data format version number
func Version() string {
	return "0.1.0"
}

func DataLoad(path string) stocks.StockData {
	if !common.FileExists(path) {
		text := fmt.Sprintf(
			"{\"name\": \"generated\", \"version\": \"%s\", \"description\": \"\", \"stocks\": []}",
			stocks.SymanticVersion())
		common.FileWrite(path, text)
	}
	data := common.FileRead(path)
	stocks := common.StructFromJson[stocks.StockData]([]byte(data))
	return stocks
}

func DataStore(path string, data stocks.StockData) {
	encoded := common.JsonFromStructPretty[stocks.StockData](data, true)
	common.FileWrite(path, string(encoded))
}

func FetchData(stock string) (services.Quote, error) {
	//NOTE: change this to use a different API
	quote, err := robowealth.GetQuote(stock)
	return quote, err
}

/* tester function to check if options look like a help request */
func HelpRequested(options ...string) bool {
	count := len(options)
	return (count == 1 && options[0] == "help") || (count < 1) 
}

func AssumeHighLow(liveData services.Quote) (float64, float64) {
	// make an assumed high/low alert
	high := liveData.FiftyTwoWeekHigh
	low := liveData.FiftyTwoWeekLow
	price := liveData.MarketPrice
	
	h := high - ((high - price) / 3)
	l := low + ((price - low) / 3)
	return l, h
}
