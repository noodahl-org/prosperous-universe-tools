package analysis

import (
	"fmt"
	"math"
	"prosperous-universe-tools/clients/fio/models"
	tx "prosperous-universe-tools/clients/memdb"
	"prosperous-universe-tools/utils"
	"strings"

	"github.com/hashicorp/go-memdb"
	"github.com/samber/lo"
)

type ExchangeSummary struct {
	TickerSummaries map[string][]TickerSummary
}

type TickerSummary struct {
	ExchangeTicker string        `json:"exchange_ticker"`
	Ticker         string        `json:"ticker"`
	Recipe         models.Recipe `json:"-"`
	Ask            Summary       `json:"ask_summary"`
	Bid            Summary       `json:"bid_summary"`
	Demand         int           `json:"demand"`
	Spread         float64       `json:"spread"`
	Cost           float64       `json:"cost"`
	Markup         string        `json:"markup"`
}

type Summary struct {
	Lowest  float64 `json:"lowest"`
	Highest float64 `json:"highest"`
	Average float64 `json:"avg"`
	Volume  int     `json:"volume"`
}

func MarketSummary(db *memdb.MemDB, out *[]TickerSummary) func() error {
	return func() error {
		data := []models.MarketData{}
		utils.Handle(tx.Select("market", "id", db, &data))

		summaries := []TickerSummary{}
		materials := lo.GroupBy(data, func(i models.MarketData) string {
			//return i.MaterialTicker
			return fmt.Sprintf("%s_%s", i.MaterialTicker, i.ExchangeCode)
		})

		for ticker, orders := range materials {
			s := ParseTickerSummary(ticker, orders)
			summaries = append(summaries, s)
		}

		//after we've parsed the initial summary data
		//we can go in and evaluate averages, recipe costs, ect
		for _, s := range summaries {
			s.Cost = lo.SumBy(s.Recipe.Inputs, func(input models.InputOutput) float64 {
				if is, _, ok := lo.FindIndexOf(summaries, func(item TickerSummary) bool {
					return item.Ticker == input.Ticker
				}); ok {
					return float64(input.Amount) * is.Ask.Lowest
				}
				return 0
			})
			s.Markup = fmt.Sprintf("%v%%", math.Round(s.Cost/s.Ask.Lowest*100))

			if len(s.Recipe.Outputs) == 1 {
				s.Cost /= float64(s.Recipe.Outputs[0].Amount)
			}
			s.Spread = math.Round(s.Bid.Highest - s.Ask.Lowest)
			*out = append(*out, s)
		}
		return nil
	}
}

func ParseTickerSummary(ticker string, orders []models.MarketData) TickerSummary {
	result := TickerSummary{
		ExchangeTicker: ticker,
		Ticker:         strings.Split(ticker, "_")[0],
		Ask:            Summary{},
		Bid:            Summary{},
	}
	recs := lo.Filter(models.StaticRecipeList(), func(i models.Recipe, index int) bool {
		for _, output := range i.Outputs {
			if output.Ticker == result.Ticker {
				return true
			}
		}
		return false
	})
	if len(recs) >= 1 {
		result.Recipe = recs[0]
	}

	asks := lo.Filter(orders, func(i models.MarketData, index int) bool {
		return i.Type == "ask"
	})

	bids := lo.Filter(orders, func(i models.MarketData, index int) bool {
		return i.Type == "bid"
	})

	result.Ask.StatisticalSummary(asks)
	result.Bid.StatisticalSummary(bids)

	result.Demand = result.Bid.Volume - result.Ask.Volume
	return result
}

func (s *Summary) StatisticalSummary(inputs []models.MarketData) {
	s.Volume = lo.SumBy(inputs, func(i models.MarketData) int {
		if i.ItemCount <= 0 {
			return 0
		}
		if s.Lowest == 0 || i.ItemCost < s.Lowest {
			s.Lowest = i.ItemCost
		}
		if s.Highest == 0 || i.ItemCost > s.Highest {
			s.Highest = i.ItemCost
		}
		s.Average += i.ItemCost * float64(i.ItemCount)
		return i.ItemCount
	})
	//handle divide by zero
	if s.Volume != 0 {
		s.Average /= float64(s.Volume)
	} else {
		s.Average = 0
	}

}
