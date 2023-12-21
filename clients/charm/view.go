package charm

import tea "github.com/charmbracelet/bubbletea"

type View interface {
	Init()
	HandleCmd(tea.Msg) (string, tea.Cmd)
	Render() string
}
