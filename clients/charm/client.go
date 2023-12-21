package charm

import (
	"bytes"
	"log"
	"prosperous-universe-tools/clients/fio"
	mdb "prosperous-universe-tools/clients/memdb"

	tea "github.com/charmbracelet/bubbletea"
)

type CharmUI struct {
	currentView string
	vb          *bytes.Buffer
	views       map[string]View
	fc          *fio.FIOClient
	db          *mdb.MemDBClient
}

type CharmUIOpts struct{}

var CharmUIOptions CharmUIOpts

type CharmUIOpt func(c *CharmUI)

func NewCharmUI(opts ...CharmUIOpt) *CharmUI {
	ui := &CharmUI{
		currentView: "dashboard",
		vb:          &bytes.Buffer{},
	}
	for _, opt := range opts {
		opt(ui)
	}
	ui.views = map[string]View{
		"dashboard": NewDashboardView(ui.db),
	}
	return ui
}

func (CharmUIOpts) MemDB(db *mdb.MemDBClient) CharmUIOpt {
	return func(c *CharmUI) {
		c.db = db
	}
}

func (c *CharmUI) Start() {
	if _, err := tea.NewProgram(c).Run(); err != nil {
		log.Panic(err)
	}
}

func (c CharmUI) Init() tea.Cmd {
	for _, val := range c.views {
		val.Init()
	}
	return tea.ClearScreen
}

func (c CharmUI) View() string {
	c.vb.Truncate(0)
	c.vb.WriteString("prosperous universe tools")
	c.vb.WriteString(c.views[c.currentView].Render())
	return c.vb.String()
}

func (c CharmUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return c, tea.Quit
		default:
			view, cmd := c.views[c.currentView].HandleCmd(msg)
			c.currentView = view
			return c, cmd
		}
	}
	return c, nil
}
