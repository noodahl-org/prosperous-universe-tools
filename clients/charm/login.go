package charm

import (
	"fmt"
	"prosperous-universe-tools/clients/fio/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginView struct {
	focus  int
	inputs []textinput.Model
	button *string
	auth   models.AuthRequest
}

func NewLoginView() View {
	cursor := textinput.New().Cursor
	inputs := []textinput.Model{
		{
			Placeholder: "fio.net username",
			PromptStyle: FocusedStyle,
			TextStyle:   FocusedStyle,
			Cursor:      cursor,
		},
		{
			Placeholder:   "fio.net password",
			EchoMode:      textinput.EchoPassword,
			EchoCharacter: '•',
			Cursor:        cursor,
		},
	}

	return &LoginView{
		focus:  0,
		inputs: inputs,
		button: &BlurredButton,
	}
}

func (l LoginView) Init() {

}

func (l *LoginView) HandleCmd(msg tea.Msg) (string, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return "", tea.Quit
		// Set focus to next input
		case "backspace":
			val := l.inputs[l.focus].Value()
			if len(val) > 1 {
				l.inputs[l.focus].SetValue(val[:len(val)-1])
			}

		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && l.focus == len(l.inputs) {
				return "login", tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				l.focus--
			} else {
				l.focus++
			}

			switch {
			case l.focus > len(l.inputs):
				l.focus = 0
			case l.focus == len(l.inputs):
				l.button = &FocusedButton
			case l.focus < 0:
				l.focus = len(l.inputs)
			}

			cmds := make([]tea.Cmd, len(l.inputs))
			for i := 0; i <= len(l.inputs)-1; i++ {
				if i == l.focus {
					// Set focused state
					cmds[i] = l.inputs[i].Focus()
					l.inputs[i].PromptStyle = FocusedStyle
					l.inputs[i].TextStyle = FocusedStyle
					continue
				}
				// Remove focused state
				l.inputs[i].Blur()
				l.inputs[i].PromptStyle = NoStyle
				l.inputs[i].TextStyle = NoStyle
			}

			return "login", tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := l.updateInputs(msg)

	return "login", cmd
}

func (l *LoginView) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(l.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range l.inputs {
		l.inputs[i], cmds[i] = l.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (l LoginView) Render() string {
	return fmt.Sprintf(`
%s
%s
%s
`, l.inputs[0].View(), l.inputs[1].View(), *l.button)
}
