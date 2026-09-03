package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type clearErrorMsg struct{}

type model struct {
	weightInput         textinput.Model
	quitting            bool
	redWineVinegar      float64
	worcestershireSauce float64
	salt                float64
	pepperCorn          float64
	corianderSeed       float64
	inputErr            string
	errStyle            lipgloss.Style
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m *model) calculate() tea.Cmd {
	weight, err := strconv.ParseFloat(m.weightInput.Value(), 64)
	if err != nil {
		m.inputErr = "Input must be a number"
		return clearErrorCmd()
	}

	m.redWineVinegar = 0.02643 * weight
	m.worcestershireSauce = 0.01322 * weight
	m.salt = 0.02247 * weight
	m.pepperCorn = 0.00749 * weight
	m.corianderSeed = 0.015 * weight

	return nil
}

func clearErrorCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func initialModel() model {
	wi := textinput.New()
	wi.Placeholder = "1000"
	wi.SetVirtualCursor(false)
	wi.Focus()
	wi.CharLimit = 156
	wi.SetWidth(20)

	return model{
		weightInput:         wi,
		redWineVinegar:      26.43,
		worcestershireSauce: 13.22,
		salt:                22.47,
		pepperCorn:          7.49,
		corianderSeed:       15,
		errStyle:            lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case clearErrorMsg:
		m.inputErr = ""
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl + c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			calcCmd := m.calculate()
			if calcCmd != nil {
				cmds = append(cmds, calcCmd)
			}
		}
	}

	m.weightInput, cmd = m.weightInput.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var c *tea.Cursor
	if !m.weightInput.VirtualCursor() {
		c = m.weightInput.Cursor()
		c.Y += lipgloss.Height(m.headerView())
	}

	errView := ""
	if m.inputErr != "" {
		errView = m.errStyle.Render(m.inputErr)
	}

	str := lipgloss.JoinVertical(
		lipgloss.Top,
		m.headerView(),
		m.weightInput.View(),
		errView,
		m.ingredientsView(),
		m.footerView(),
	)

	if m.quitting {
		str += "\n"
	}

	v := tea.NewView(str)
	v.Cursor = c
	return v
}

func (m model) headerView() string { return "Weight in grams" }
func (m model) ingredientsView() string {
	return fmt.Sprintf("Red Wine Vinegar: %.2fml\nWorcestershire Sauce: %.2fml\nSalt: %.2fg\nPepper Corn: %.2fg\nCoriander Seed: %.2fg", m.redWineVinegar, m.worcestershireSauce, m.salt, m.pepperCorn, m.corianderSeed)
}
func (m model) footerView() string { return "\n(esc to quit)" }
