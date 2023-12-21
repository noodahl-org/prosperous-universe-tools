package charm

import (
	"bytes"
	"fmt"
	"prosperous-universe-tools/clients/analysis"
	"prosperous-universe-tools/clients/fio/models"
	mdb "prosperous-universe-tools/clients/memdb"
	"prosperous-universe-tools/utils"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
)

type DashboardView struct {
	vb          *bytes.Buffer
	focus       int
	tableTitle  string
	graph       string
	table       table.Model
	orders      table.Model
	input       textinput.Model
	suggestions []string
	db          *mdb.MemDBClient
}

func NewDashboardView(db *mdb.MemDBClient) View {
	dbv := &DashboardView{
		db: db,
		vb: &bytes.Buffer{},
		input: textinput.Model{
			Placeholder:     "ticker",
			PromptStyle:     FocusedStyle,
			TextStyle:       FocusedStyle,
			ShowSuggestions: true,
			Cursor:          textinput.New().Cursor,
		},
		suggestions: []string{},
		table:       table.Model{},
		focus:       0,
	}
	dbv.input.Focus()
	return dbv
}

func (v *DashboardView) Init() {
	ts := []analysis.TickerSummary{}
	utils.Handle(mdb.Select("ticker", "id", nil, v.db.DB(), &ts))
	tmp := lo.FilterMap(ts, func(x analysis.TickerSummary, _ int) (string, bool) {
		return strings.ToLower(x.Ticker), true
	})
	v.suggestions = lo.Union(tmp)
	sort.Strings(v.suggestions)
	v.input.SetSuggestions(v.suggestions)
}

func (v *DashboardView) HandleCmd(msg tea.Msg) (string, tea.Cmd) {
	var cmd tea.Cmd
	cmds := []tea.Cmd{}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			val := v.input.Value()
			if len(val) > 0 {
				v.input.SetValue(val[:len(val)-1])
			}
		case "tab":
			if v.input.Focused() {
				v.input.Blur()
				v.table.Focus()
			} else {
				v.table.Blur()
				v.input.Focus()
			}
		case "down":
			v.table.MoveDown(1)
			exchangeTicker := []string(v.table.SelectedRow())[0]
			tmp := []models.MarketData{}
			utils.Handle(mdb.Select("market", "exchange_ticker", exchangeTicker, v.db.DB(), &tmp))
			v.graph = MarketSummaryGraph(tmp)
			v.table.Update(msg)
		case "up":
			v.table.MoveUp(1)
			exchangeTicker := []string(v.table.SelectedRow())[0]
			tmp := []models.MarketData{}
			utils.Handle(mdb.Select("market", "exchange_ticker", exchangeTicker, v.db.DB(), &tmp))
			v.graph = MarketSummaryGraph(tmp)
			v.table.Update(msg)
		case "enter":
			mat := models.Material{}
			tk := []analysis.TickerSummary{}
			md := []models.MarketData{}
			exchangeTicker := "AI1." + strings.ToUpper(v.input.Value()) //stupid default, im tired
			utils.Handle(mdb.Select("ticker", "ticker", strings.ToUpper(v.input.Value()), v.db.DB(), &tk))
			utils.Handle(mdb.SelectOne("material", "ticker", strings.ToUpper(v.input.Value()), v.db.DB(), &mat))
			utils.Handle(mdb.Select("market", "exchange_ticker", exchangeTicker, v.db.DB(), &md))

			v.tableTitle = mat.Name
			v.table = TickerSummaryTable(tk)
			v.graph = MarketSummaryGraph(md)
			v.table, cmd = v.table.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	v.input, cmd = v.input.Update(msg)
	cmds = append(cmds, cmd)
	return "dashboard", tea.Batch(cmds...)
}

func (v DashboardView) Render() string {
	return fmt.Sprintf(`
%s
%s
%s
%s
`,
		v.input.View(),
		v.tableTitle,
		v.table.View(),
		v.graph)
}
