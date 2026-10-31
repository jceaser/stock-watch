package finnhub

/*
https://finnhub.io/dashboard

experimental support, not impressed

*/


import (
	//standard
	"context"
	"fmt"

	_ "embed"

	//local
	"github.com/jceaser/stock-watch/services"

	//public
    finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

var (
	// Put the api key in api-key.text.
	//
	//go:embed api-key.txt
	apiKey string
)

func GetQuote() (services.Quote, error) {

	cfg := finnhub.NewConfiguration()
    cfg.AddDefaultHeader("X-Finnhub-Token", apiKey)
    finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	//quote, _, err := finnhubClient.Quote(context.Background(), "AAPL")
    quote, _, err := finnhubClient.Quote(context.Background()).Symbol("AAPL").Execute()
    fmt.Printf("%+v\n", quote)
    if err != nil {
    	panic(err)
    }

	//price, _, err := finnhubClient.PriceMetrics(context.Background(), "AAPL")
    price, _, err := finnhubClient.PriceMetrics(context.Background()).Symbol("AAPL").Execute()
    fmt.Printf("%+v\n", price)
    if err != nil {
    	panic(err)
    }


	fmt.Printf("%v\n", price)

	return services.CreateQuote("AAPL",
		"",//quote.ShortName,
		"",
		0,
		float64(*quote.C),
		0.0,//q.FiftyTwoWeekLow,
		0.0)//q.FiftyTwoWeekHigh)
}
