//***********************************************************************************************100

package services

type Quote struct {
	Symbol				string
	ShortName			string
	DateTime			string
	TimeStamp			int64
	MarketPrice			float64
	FiftyTwoWeekLow		float64
	FiftyTwoWeekHigh	float64
}

func CreateQuote(symbol, name, dateTime string,
		timeStamp int64,
		price, week52Low, week52High float64) (Quote, error) {
	return Quote{Symbol: symbol,
		ShortName: name,
		DateTime: dateTime,
		TimeStamp: timeStamp,
		MarketPrice: price,
		FiftyTwoWeekLow: week52Low,
		FiftyTwoWeekHigh: week52High},
		nil
}