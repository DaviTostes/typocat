package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var frames = []string{
	` /\_/\
( ^.^ )
 / >🍪 
/|_|_|\`,

	` /\_/\
( ^o^ )
 / >☕
/|_|_|\`,
}

var incorrectFrame = `  /\_/\  
 ( T.T ) 
  / >🥀  
 /|_|_|\ `

var colors = []string{
	"#FF0000", // red
	"#FF7F00", // orange
	"#FFFF00", // yellow
	"#00FF00", // green
	"#00BFFF", // sky blue
	"#0000FF", // blue
	"#8B00FF", // violet

	"#FF1493", // deep pink
	"#FF69B4", // hot pink
	"#FF4500", // orange red
	"#FFD700", // gold
	"#ADFF2F", // green yellow
	"#00FA9A", // medium spring green
	"#20B2AA", // light sea green
	"#1E90FF", // dodger blue
	"#4169E1", // royal blue
	"#6A5ACD", // slate blue
	"#9370DB", // medium purple
	"#BA55D3", // orchid
	"#FF00FF", // magenta
}

var texts = []string{
	`Once upon a midnight dreary, While I pondered, weak and weary, Over many a quaint and curious Volume of forgotten lore- While I nodded, nearly napping, Suddenly there came a tapping, As of some one gently rapping, Rapping at my chamber door. "'T is some visitor," I muttered, "Tapping at my chamber door Only this and nothing more."`,
}

var combosTexts = map[int]func(c int) string{
	0:   func(c int) string { return "Yay!" },
	10:  func(c int) string { return fmt.Sprintf("Nice. %dx combo", c) },
	20:  func(c int) string { return fmt.Sprintf("Amazing! %dx combo", c) },
	50:  func(c int) string { return fmt.Sprintf("Incredible!! %dx combo", c) },
	100: func(c int) string { return fmt.Sprintf("SPLENDID!!!. %dX COMBO", c) },
}

func randomColor() string {
	r := rand.Intn(len(colors) - 1)
	return colors[r]
}

type model struct {
	viewport    viewport.Model
	keys        []string
	textKeys    []string
	correct     bool
	correctText string
	frame       int
	combo       int
}

type comboMsg struct {
	fn func(int) string
}

func (m *model) comboCmd(fn func(int) string) tea.Cmd {
	return func() tea.Msg {
		return comboMsg{
			fn: fn,
		}
	}
}

func InitModel() model {
	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	textKeys := strings.Split(texts[0], "")

	return model{
		viewport:    vp,
		keys:        []string{},
		textKeys:    textKeys,
		correct:     false,
		correctText: "",
		frame:       0,
		combo:       0,
	}
}

func (m *model) renderViewportContent() {
	frame := incorrectFrame
	if m.correct {
		if m.combo < 10 {
			m.correctText = "Yay!"
		}
		fn := combosTexts[m.combo]
		if fn != nil {
			m.correctText = fn(m.combo)
		}

		m.frame ^= 1
		frame = frames[m.frame]
	}

	cat := lipgloss.NewStyle().
		Foreground(lipgloss.Color(randomColor())).
		Align(lipgloss.Center).
		Render(
			"TYPOCAT\n\n" +
				frame,
		)

	text := ""
	for i := range m.textKeys {
		if len(m.keys)-1 < i {
			text += lipgloss.NewStyle().Foreground(lipgloss.Color("#5D5D5D")).Render(m.textKeys[i])
			continue
		}

		if strings.Compare(m.keys[i], m.textKeys[i]) == 0 {
			text += lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Render(m.textKeys[i])
			continue
		}

		text += lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render(m.textKeys[i])
	}

	text = lipgloss.NewStyle().
		Width(70).
		Align(lipgloss.Center).
		Render(text)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		cat,
		"",
		text,
		"",
		m.correctText,
	)

	content = lipgloss.Place(
		m.viewport.Width(),
		m.viewport.Height(),
		lipgloss.Center,
		lipgloss.Center,
		content,
	)

	m.viewport.SetContent(content)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height)
		m.viewport.GotoBottom()

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "backspace":
			if len(m.keys) > 0 {
				m.keys = m.keys[:len(m.keys)-1]
			}
			m.correct = false

		default:
			m.keys = append(m.keys, msg.String())
			if len(m.keys) > 0 {
				m.correct = m.keys[len(m.keys)-1] == m.textKeys[len(m.keys)-1] ||
					(msg.String() == "space" && m.textKeys[len(m.keys)-1] == " ")

				if m.correct {
					m.combo++
					return m, m.comboCmd(combosTexts[m.combo])
				} else {
					m.combo = 0
				}
			}
		}

	case comboMsg:
		if msg.fn != nil {
			m.correctText = msg.fn(m.combo)
		}
	}

	m.renderViewportContent()

	return m, nil
}

func (m model) View() tea.View {
	vpView := m.viewport.View()
	v := tea.NewView(vpView)
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(InitModel())
	if _, err := p.Run(); err != nil {
		log.Fatalln(err.Error())
		os.Exit(1)
	}
}
