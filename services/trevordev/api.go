package trevordev

//go get github.com/TrevorDev/go-finance

import (
	"fmt"
	"github.com/jceaser/stock-watch/services"

	"github.com/TrevorDev/go-finance"
)

func GetQuote() (services.Quote, error) {

	out, err := finance.GetStockInfo([]string{"APPL"},
		[]string{finance.Last_Trade_Price_Only, finance.Price_Per_Earning_Ratio})
	if(err!=nil){
		//This may happen often if yahoo's server is under high volume
		fmt.Println(err)
		return services.Quote{}, err
	}else{
		//outputs map of stock symbols and their specified attributes given above
		fmt.Println(out)

		//outputs microsofts current price per earing ratio as a string
		fmt.Println(out["APPL"][finance.Price_Per_Earning_Ratio])
	}

	return services.Quote{}, nil
}