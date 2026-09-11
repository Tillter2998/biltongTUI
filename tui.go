package tui

import (
	"fmt"
	"maps"
	"slices"
	"strconv"

	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/Tillter2998/biltongTUI/internal/enums"
	"github.com/Tillter2998/biltongTUI/internal/quicksort"
)

type clearErrorMsg struct{}

type ingredientValues struct {
	amount          float64
	ratio           float64
	measurementUnit string
}

type optionsModel struct {
	cursor   int
	choices  []enums.Ingredients
	selected map[enums.Ingredients]struct{}
}

type weightInputModel struct {
	input    textinput.Model
	inputErr string
	errStyle lipgloss.Style
}

type ingredientValuesModel struct {
	values map[enums.Ingredients]*ingredientValues
}

type model struct {
	options     optionsModel
	weightInput weightInputModel
	ingredients ingredientValuesModel
	focus       int
	quitting    bool
}

const (
	focusInput = iota
	focusMenu
)

func NewModel() tea.Model {
	return initialModel()
}

func (m *model) calculateAll() tea.Cmd {
	weight, err := m.getWeight()
	if err != nil {
		m.weightInput.inputErr = "Input must be a number"
		return clearErrorCmd()
	}
	if weight <= 0 {
		m.weightInput.inputErr = "Input must be greater than 0"
		return clearErrorCmd()
	}

	for ingredient := range m.options.selected {
		m.ingredients.values[ingredient].amount = m.ingredients.values[ingredient].ratio * weight
	}

	return nil
}

func (m *model) calculate(ingredient enums.Ingredients) tea.Cmd {

	weight, err := m.getWeight()
	if err != nil {
		m.weightInput.inputErr = "Input must be a number"
		return clearErrorCmd()
	}
	if weight <= 0 {
		m.weightInput.inputErr = "Input must be greater than 0"
		return clearErrorCmd()
	}

	m.ingredients.values[ingredient].amount = m.ingredients.values[ingredient].ratio * weight

	return nil
}

func (m *model) getWeight() (float64, error) {
	input := m.weightInput.input.Value()
	if input == "" {
		return 1000, nil
	}

	return strconv.ParseFloat(input, 64)
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
	wi.SetVirtualCursor(true)
	wi.Blur()
	wi.CharLimit = 156
	wi.SetWidth(20)

	defaultIngredients := []enums.Ingredients{
		enums.RedWineVinegar,
		enums.WorcestershireSauce,
		enums.Salt,
		enums.PepperCorn,
		enums.CorianderSeed,
	}

	choices := enums.AllIngredients()

	selected := make(map[enums.Ingredients]struct{})
	for _, value := range defaultIngredients {
		selected[enums.Ingredients(value)] = struct{}{}
	}

	return model{
		options: optionsModel{
			choices:  choices,
			selected: selected,
			cursor:   len(defaultIngredients),
		},
		weightInput: weightInputModel{
			input:    wi,
			errStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		},

		ingredients: ingredientValuesModel{
			values: map[enums.Ingredients]*ingredientValues{
				enums.RedWineVinegar: {
					amount:          26.43,
					ratio:           0.02643,
					measurementUnit: "ml",
				},
				enums.WorcestershireSauce: {
					amount:          13.22,
					ratio:           0.01322,
					measurementUnit: "ml",
				},
				enums.Salt: {
					amount:          22.47,
					ratio:           0.02247,
					measurementUnit: "g",
				},
				enums.PepperCorn: {
					amount:          7.49,
					ratio:           0.00749,
					measurementUnit: "g",
				},
				enums.CorianderSeed: {
					amount:          15,
					ratio:           0.015,
					measurementUnit: "g",
				},
			},
		},
		focus: focusMenu,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case clearErrorMsg:
		m.weightInput.inputErr = ""
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			if m.focus == focusInput {
				m.focus = focusMenu
				m.weightInput.input.Blur()
			} else {
				m.focus = focusInput
				m.weightInput.input.Focus()
			}
			return m, nil
		}

		if m.focus == focusInput {
			switch msg.String() {
			case "enter":
				calcCmd := m.calculateAll()
				if calcCmd != nil {
					cmds = append(cmds, calcCmd)
				}
			}

			m.weightInput.input, cmd = m.weightInput.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch msg.String() {
			case "k", "up":
				if m.options.cursor > 0 {
					m.options.cursor--
				}
			case "j", "down":
				if m.options.cursor < len(m.options.choices)-1 {
					m.options.cursor++
				}
			case "enter":
				ingredient := m.options.choices[m.options.cursor]
				_, ok := m.options.selected[ingredient]
				if ok {
					delete(m.options.selected, ingredient)
				} else {
					m.options.selected[ingredient] = struct{}{}
					calcCmd := m.calculate(ingredient)
					if calcCmd != nil {
						cmds = append(cmds, calcCmd)
					}

				}
			}
		}
	}

	// if m.focus == focusInput {
	// 	m.weightInput.input, cmd = m.weightInput.input.Update(msg)
	// 	cmds = append(cmds, cmd)
	// }
	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	var c *tea.Cursor
	if !m.weightInput.input.VirtualCursor() {
		c = m.weightInput.input.Cursor()

		// FIX: Only modify c if it's not nil (it will be nil when blurred!)
		if c != nil {
			topHeight := lipgloss.Height(m.optionsView()) + lipgloss.Height(m.inputView())
			c.Y += topHeight
		}
	}

	str := lipgloss.JoinVertical(
		lipgloss.Top,
		m.optionsView(),
		m.inputView(),
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

var sectionStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("7")).
	Padding(0, 1).
	Width(40)

func (m model) optionsView() string {

	var s strings.Builder
	s.WriteString("Options:\n")
	for i, choice := range m.options.choices {
		cursor := " "
		if m.options.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.options.selected[choice]; ok {
			checked = "x"
		}
		fmt.Fprintf(&s, "%s [%s] %s\n", cursor, checked, choice)
	}

	borderColor := "7"
	if m.focus == focusMenu {
		borderColor = "5"
	}
	return sectionStyle.
		BorderForeground(lipgloss.Color(borderColor)).
		Render(s.String())
}
func (m model) inputView() string {
	s := "Weight in Grams:\n"

	s += fmt.Sprintf("%s\n", m.weightInput.input.View())

	errView := ""
	if m.weightInput.inputErr != "" {
		errView = m.weightInput.errStyle.Render(m.weightInput.inputErr)
	}

	s += fmt.Sprintf("%s", errView)

	borderColor := "7"
	if m.focus == focusInput {
		borderColor = "5"
	}
	return sectionStyle.
		BorderForeground(lipgloss.Color(borderColor)).
		Render(s)

}
func (m model) ingredientsView() string {
	var s strings.Builder

	sortedSelected := slices.Collect(maps.Keys(m.options.selected))
	quicksort.Quicksort(sortedSelected, 0, len(m.options.selected)-1)

	for _, ingredient := range sortedSelected {
		amount := m.ingredients.values[ingredient].amount
		measurement := m.ingredients.values[ingredient].measurementUnit

		fmt.Fprintf(&s, "%s: %.2f%s\n", ingredient.String(), amount, measurement)
	}

	return sectionStyle.Render(s.String())
}
func (m model) footerView() string { return "\n(esc to quit)" }
