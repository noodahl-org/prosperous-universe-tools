package charm

import (
	"cmp"
	"fmt"
	"math"
	"prosperous-universe-tools/clients/analysis"
	"prosperous-universe-tools/clients/fio/models"
	"slices"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/guptarohit/asciigraph"
	"github.com/samber/lo"
)

var (
	FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle         = FocusedStyle.Copy()
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = BlurredStyle.Copy()
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	FocusedButton = FocusedStyle.Copy().Render("[ Submit ]")
	BlurredButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Submit"))
)

func TickerSummaryTable(data []analysis.TickerSummary) table.Model {
	columns := []table.Column{
		{
			Title: "EX",
			Width: 8,
		},
		{
			Title: "BID",
			Width: 6,
		},
		{
			Title: "AVG",
			Width: 6,
		},
		{
			Title: "ASK",
			Width: 6,
		},
		{
			Title: "SPREAD",
			Width: 8,
		},
		{
			Title: "BUY/SELL",
			Width: 16,
		},
		{
			Title: "DEMAND",
			Width: 10,
		},
		{
			Title: "COST",
			Width: 6,
		},
		{
			Title: "MARKUP",
			Width: 6,
		},
	}
	rows := []table.Row{}
	for _, summary := range data {
		rows = append(rows, table.Row{
			summary.ExchangeTicker,
			fmt.Sprintf("%v", math.Round(summary.Bid.Highest)),
			fmt.Sprintf("%v", math.Round(summary.Bid.Average)),
			fmt.Sprintf("%v", math.Round(summary.Ask.Lowest)),
			fmt.Sprintf("%v", math.Round(summary.Spread)),
			fmt.Sprintf("%d-%d", summary.Bid.Volume, summary.Ask.Volume),
			fmt.Sprintf("%v", summary.Demand),
			fmt.Sprintf("%v", math.Round(summary.Cost)),
			fmt.Sprintf("%v", summary.Markup),
		})
	}
	return table.New(
		table.WithColumns(columns),
		table.WithFocused(false),
		table.WithRows(rows),
		table.WithHeight(8),
	)
}

func MarketSummaryGraph(data []models.MarketData) string {
	var graph string
	orders := lo.GroupBy(data, func(i models.MarketData) string {
		return i.Type
	})

	asks := orders["ask"]
	slices.SortFunc(asks, func(prev, next models.MarketData) int {
		return cmp.Compare(prev.ItemCost, next.ItemCost)
	})

	bids := orders["bid"]
	slices.SortFunc(bids, func(prev, next models.MarketData) int {
		return cmp.Compare(prev.ItemCost, next.ItemCost)
	})
	slices.Reverse(bids)
	totalVol := lo.Sum(lo.FlatMap(bids, func(i models.MarketData, index int) []float64 {
		return []float64{float64(i.ItemCount)}
	}))
	if totalVol == 0 {
		return ""
	}
	bidVol := lo.FlatMap(bids, func(order models.MarketData, index int) []float64 {
		result := []float64{}
		percent := (float64(order.ItemCount) / totalVol) * 100
		percent = math.RoundToEven(percent)
		for i := 0; i <= int(percent)/10; i++ {
			result = append(result, float64(order.ItemCount))
		}
		return result
	})

	if len(bidVol) > 0 {
		graph = fmt.Sprintf("bid-vol\n\n%s\n%10v %50v\n%36s", asciigraph.Plot(bidVol, []asciigraph.Option{
			asciigraph.Width(52),
			asciigraph.Height(16),
			asciigraph.Offset(0),
			asciigraph.AxisColor(asciigraph.AliceBlue),
			asciigraph.SeriesColors(
				asciigraph.HotPink),
		}...), math.Round(bids[0].ItemCost), math.Round(bids[len(bids)-1].ItemCost), "price")

	}
	return graph
}
