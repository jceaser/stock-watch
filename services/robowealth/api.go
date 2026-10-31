package robowealth

import (
	"time"

	"github.com/jceaser/stock-watch/services"

	"github.com/robowealth-mutual-fund/finance-go/quote"
)

func GetQuote(symbol string) (services.Quote, error) {
	q, err := quote.Get(symbol)
	if err != nil {
		return services.Quote{}, err
	}

	marketTime := int64(q.RegularMarketTime)
    t := time.Unix(marketTime, 0)

	// Success!

	quote, _ := services.CreateQuote(q.Symbol,
		q.ShortName,
      	t.Format(time.RFC3339), //iso format
      	marketTime,
		q.RegularMarketPrice,
		q.FiftyTwoWeekLow,
		q.FiftyTwoWeekHigh)

	return quote, nil
}
