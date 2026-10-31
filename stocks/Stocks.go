//***********************************************************************************************100
/*
Data Model for tracking stocks. This struct is designed to be loaded and stored in json.
*/

package stocks

import (
	"fmt"
	"strings"
	//"github.com/rymdport/portal/account"
)

/*
Example data this struct will handle
{
	"name": primary,
	"version": "0.1.0",
	"description": "the main list of stocks",
	"stocks": [
		{
			"id": "APPL",
			"name": "Apple",
			"description": "Apple Inc.",
			"note": "buy as much as you can",
			"last": "2026-10-25",
			"time-stamp": 1792785601,
			"alert": {
				"low": 100.0
				"high": 512.0
			},
			"prices": [
				{
					"date-time": "2026-09-01",
					"value": 256.0
				}
			],
			"inventory": [
				{
					"count": 6
					"value": 128.0
				}
			],
			"week-52": {
				"low": 100.0
				"high": 512.0
			}
		}
	]
}
*/

/* wrapper for stock information */
type StockData struct {
	Name   string      `json:"name"`
	Version string	   `json:"version"`
	Description string `json:"description"`
	Stocks []StockItem `json:"stocks"`
}

/* A single stock */
type StockItem struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Note        string           `json:"note"`
	Last        string           `json:"last"`
	TimeStamp	int64			 `json:"time-stamp"`
	Price		float64			 `json:"price"`
	Alert       RangeData        `json:"alert"`
	Prices      []PricePoint     `json:"prices"`
	Inventory   []InventoryItem  `json:"inventory"`
	Week52      RangeData        `json:"week-52"`
}

/* A number range from low to high */
type RangeData struct {
	Low  float64 `json:"low"`
	High float64 `json:"high"`
}

/* A single price point in time */
type PricePoint struct {
	DateTime string  `json:"date-time"`
	Value    float64 `json:"value"`
}

/* An inventory of stock purchased at the same price */
type InventoryItem struct {
	Count float64     `json:"count"` // may need to be a float when not tracking stocks
	Value float64 `json:"value"`
}

//****************************************************************************80
//MARK - Methods

func (s StockItem) String() string {
    return fmt.Sprintf("%s (%s) %s\n" +
    	"\tNote: %s\n" +
    	"\tAlerts: $%0.2f < $%0.2f\n" +
    	"\tInventory: %0.3f x $%0.2f = $%0.2f\n" +
    	"\tActuals:\n" +
    	"\t\tPrices: $%0.2f < $%0.2f < $%0.2f\n" +
    	"\t\tLast updated: %s\n" +
    	"\t\tProfit: $%0.2f\n",
    	s.ID,
    	s.Name,
    	s.Description,
    	s.Note,
    	s.Alert.Low, s.Alert.High,
    	s.InventoryCount(), s.AverageItemCost(), s.Cost(),
    	s.Week52.Low, s.Price, s.Week52.High,
    	s.Last,
    	s.Profit())
}

//AI should implement this function to populate data in StockItem.Prices
func (obj *StockItem) AddPrice(dateTime string, value float64) {
	data := PricePoint{dateTime, value}
	obj.Prices = append(obj.Prices, data)
	if len(obj.Prices) > 10 {
		obj.Prices = obj.Prices[1:]
	}
}

//AI should implement this function to populate data in StockItem.Inventory
func (obj *StockItem) AddInventory(count float64, value float64) {
	item := InventoryItem{count, value}
	obj.Inventory = append(obj.Inventory, item)
}

//****************************************************************************80
//MARK - Helpers

/* look at all the stocks to find a match */
func (data StockData) Exists(stockId string) (StockItem, bool) {
	for _, stock := range data.Stocks {
		if stock.ID == strings.ToUpper(stockId) {
			return stock, true
		}
	}
	return StockItem{}, false
}

/* what did it cost to buy the stock */
func (obj StockItem) Cost() float64 {
	if obj.Inventory == nil || len(obj.Inventory) < 1 {
		return float64(len(obj.Inventory))
	}
	accumulator := 0.0
	for _, item := range obj.Inventory {
		accumulator += item.Value * float64(item.Count)
	}
	return accumulator
}

/* what is the stock worth against the current price */
func (obj StockItem) Value() float64 {
	if obj.Inventory == nil || len(obj.Inventory) < 1 {
		return 0.0
	}
	if obj.Price <= 0.0 {
		return 0.0
	}
	accumulator := 0.0
	for _, item := range obj.Inventory {
		accumulator += (obj.Price * float64(item.Count))
	}
	return accumulator
}

/* what was made or lost on the stock given the current price */
func (obj StockItem) Profit() float64 {
	if obj.Price <= 0.0 {
		return 0.0
	}
	return obj.Value() - obj.Cost()
}

/* Number of stocks held at any price */
func (obj StockItem) InventoryCount() float64 {
	if obj.Inventory == nil || len(obj.Inventory) < 1 {
		return 0
	}
	accumulator := 0.0
	for _, item := range obj.Inventory {
		accumulator += item.Count
	}
	return accumulator
}

/* the average cost of the stock */
func (obj StockItem) AverageItemCost() float64 {
	if obj.Inventory == nil || len(obj.Inventory) < 1 {
		return 0.0
	}
	return obj.Cost() / float64(obj.InventoryCount())
}
