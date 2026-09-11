package main

import (
	"fmt"
	"maps"
	"os"
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

type optionsModel struct {
	cursor   int
	choices  []string
	selected map[enums.Ingredients]string
}

type weightInputModel struct {
	input    textinput.Model
	inputErr string
	errStyle lipgloss.Style
}

type ingredientAmountsModel struct {
	amounts map[enums.Ingredients]float64
	ratios  map[enums.Ingredients]float64
}

type model struct {
	options     optionsModel
	weightInput weightInputModel
	ingredients ingredientAmountsModel
	focus       int
	quitting    bool
}

const (
	focusInput = iota
	focusMenu
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m *model) calculateAll() tea.Cmd {
	weight, err := m.getWeight()
	if err != nil {
		m.weightInput.inputErr = "Input must be a number"
		return clearErrorCmd()
	}

	// TODO: Update this so that it only calculates for selected ingredients
	// m.ingredients.redWineVinegar = 0.02643 * weight
	// m.ingredients.worcestershireSauce = 0.01322 * weight
	// m.ingredients.salt = 0.02247 * weight
	// m.ingredients.pepperCorn = 0.00749 * weight
	// m.ingredients.corianderSeed = 0.015 * weight

	for ingredient := range m.options.selected {
		m.ingredients.amounts[ingredient] = m.ingredients.ratios[ingredient] * weight
	}

	return nil
}

func (m *model) calculate(ingredient enums.Ingredients) tea.Cmd {

	weight, err := m.getWeight()
	if err != nil {
		m.weightInput.inputErr = "Input must be a number"
		return clearErrorCmd()
	}

	m.ingredients.amounts[ingredient] = m.ingredients.ratios[ingredient] * weight

	return nil
}

func (m *model) getWeight() (float64, error) {
	input := m.weightInput.input.Value()
	if input == "" {
		return 1000, nil
	}

	return strconv.ParseFloat(m.weightInput.input.Value(), 64)
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

	choices := ingredientsMapToSortedSlice(enums.IngredientName)

	selected := make(map[enums.Ingredients]string)
	for _, value := range defaultIngredients {
		selected[enums.Ingredients(value)] = value.String()
	}

	return model{
		options: optionsModel{
			choices:  choices,
			selected: selected,
			cursor:   len(enums.IngredientName) - 1,
		},
		weightInput: weightInputModel{
			input:    wi,
			errStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
		},
		ingredients: ingredientAmountsModel{
			amounts: map[enums.Ingredients]float64{
				enums.RedWineVinegar:      26.43,
				enums.WorcestershireSauce: 13.22,
				enums.Salt:                22.47,
				enums.PepperCorn:          7.49,
				enums.CorianderSeed:       15,
			},
			ratios: map[enums.Ingredients]float64{
				enums.RedWineVinegar:      0.02643,
				enums.WorcestershireSauce: 0.01322,
				enums.Salt:                0.02247,
				enums.PepperCorn:          0.00749,
				enums.CorianderSeed:       0.015,
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
		case "ctrl + c", "esc":
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
				converted := enums.Ingredients(m.options.cursor)
				_, ok := m.options.selected[converted]
				if ok {
					delete(m.options.selected, converted)
				} else {
					m.options.selected[converted] = converted.String()
					calcCmd := m.calculate(converted)
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

	s := "Options:\n"
	for i, choice := range m.options.choices {
		cursor := " "
		if m.options.cursor == i {
			cursor = ">"
		}
		checked := " "
		if _, ok := m.options.selected[enums.Ingredients(i)]; ok {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	borderColor := "7"
	if m.focus == focusMenu {
		borderColor = "5"
	}
	return sectionStyle.
		BorderForeground(lipgloss.Color(borderColor)).
		Render(s)
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
	sorted := ingredientsMapToSortedSlice(m.options.selected)

	for index, value := range sorted {
		amount := m.ingredients.amounts[enums.Ingredients(index)]
		fmt.Fprintf(&s, "%s: %.2fml\n", value, amount)
	}

	return sectionStyle.Render(s.String())
}
func (m model) footerView() string { return "\n(esc to quit)" }

func ingredientsMapToSortedSlice(choicesMap map[enums.Ingredients]string) []string {

	keys := slices.Collect(maps.Keys(choicesMap))
	quicksort.Quicksort(keys, 0, len(keys)-1)
	choices := make([]string, 0, len(keys))
	for _, key := range keys {
		choices = append(choices, enums.IngredientName[key])
	}

	return choices
}
