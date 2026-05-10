package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pquerna/otp/totp"
	"github.com/shishiro26/kubera/models"
)

type ListMode int

const (
	modeList ListMode = iota
	modeView
	modeConfirmDelete
)

type ListModel struct {
	entries   []models.Entry
	cursor    int
	selected  int
	mode      ListMode
	action    string
	status    string
	termWidth int
}

func NewListModel(entries []models.Entry, startCursor int, status string) ListModel {
	cursor := startCursor
	if cursor >= len(entries) {
		cursor = 0
	}
	return ListModel{
		entries: entries,
		cursor:  cursor,
		status:  status,
	}
}

func (m ListModel) Init() tea.Cmd { return nil }

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch m.mode {
		case modeList:
			return m.handleList(msg)
		case modeView:
			return m.handleView(msg)
		case modeConfirmDelete:
			return m.handleConfirmDelete(msg)
		}
	}
	return m, nil
}

func (m ListModel) handleList(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.status = ""
	switch k.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.entries)-1 {
			m.cursor++
		}
	case "enter", " ":
		if len(m.entries) > 0 {
			m.mode = modeView
		}
	case "d":
		if len(m.entries) > 0 {
			m.selected = m.cursor
			m.mode = modeConfirmDelete
		}
	case "e":
		if len(m.entries) > 0 {
			m.action = "edit"
			m.selected = m.cursor
			return m, tea.Quit
		}
	case "a":
		m.action = "add"
		return m, tea.Quit
	}
	return m, nil
}

func (m ListModel) handleView(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "esc", "backspace", "q":
		m.mode = modeList
	case "d":
		m.selected = m.cursor
		m.mode = modeConfirmDelete
	case "e":
		m.action = "edit"
		m.selected = m.cursor
		return m, tea.Quit
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m ListModel) handleConfirmDelete(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "y", "Y":
		m.action = "delete"
		return m, tea.Quit
	case "n", "N", "esc":
		m.mode = modeList
	case "ctrl+c":
		return m, tea.Quit
	}
	return m, nil
}

func (m ListModel) viewTitle() string {
	art := []string{
		"█▄▀ █ █ █▄▄ █▀▀ █▀█ ▄▀█",
		"█░█ █▄█ █▄█ ██▄ █▀▄ █▀█",
	}

	gradientColors := []lipgloss.Color{
		lipgloss.Color("#F5C842"),
		lipgloss.Color("#D4A017"),
	}

	var artLines []string
	for i, line := range art {
		styled := lipgloss.NewStyle().
			Foreground(gradientColors[i%len(gradientColors)]).
			Bold(true).
			Render(line)
		artLines = append(artLines, styled)
	}
	artBlock := strings.Join(artLines, "\n")

	totpCount := 0
	for _, e := range m.entries {
		if e.TOTP != "" {
			totpCount++
		}
	}
	counter := lipgloss.NewStyle().Foreground(ColorAccent).Render("◆") +
		lipgloss.NewStyle().Foreground(ColorMuted).Render(fmt.Sprintf("  %d entries", len(m.entries)))
	if totpCount > 0 {
		counter += lipgloss.NewStyle().Foreground(ColorSubtle).Render("  ·  ") +
			lipgloss.NewStyle().Foreground(ColorSuccess).Render(fmt.Sprintf("%d with 2FA", totpCount))
	}

	inner := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(artBlock + "\n" + counter)

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(ColorPrimary).
		Padding(0, 3).
		MarginBottom(1)

	if m.termWidth > 0 {
		titleWidth := lipgloss.Width(style.Render(inner))
		leftPad := (m.termWidth - titleWidth) / 2
		if leftPad < 0 {
			leftPad = 0
		}
		return style.MarginLeft(leftPad).Render(inner)
	}
	return style.Render(inner)
}

func (m ListModel) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(m.viewTitle() + "\n\n")

	switch m.mode {
	case modeList:
		b.WriteString(m.viewList())
		if m.status != "" {
			statusBox := lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(ColorSuccess).
				Padding(0, 2).
				MarginLeft(2).
				MarginTop(1).
				Render(SuccessStyle.Render("  ❯❯  " + m.status))
			b.WriteString("\n" + statusBox)
		}
		b.WriteString("\n\n" + m.viewHelp())

	case modeView:
		b.WriteString(m.viewList())
		b.WriteString("\n" + m.viewDetail())
		b.WriteString("\n\n" + HelpStyle.Render("  esc  back   e  edit   d  delete"))

	case modeConfirmDelete:
		b.WriteString(m.viewList())
		if m.selected < len(m.entries) {
			e := m.entries[m.selected]
			warning := lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(ColorError).
				Padding(0, 2).
				MarginLeft(2).
				MarginTop(1).
				Render(
					WarningStyle.Render("  !!  ") +
						ValueStyle.Render("Delete ") +
						LabelStyle.Render(e.Username) +
						SubtleStyle.Render(" @ ") +
						LabelStyle.Render(e.Site) +
						ValueStyle.Render("?  ") +
						LabelStyle.Render("y") +
						SubtleStyle.Render(" confirm   ") +
						LabelStyle.Render("n") +
						SubtleStyle.Render(" cancel"),
				)
			b.WriteString("\n" + warning)
		}
	}

	return b.String()
}

func (m ListModel) viewList() string {
	if len(m.entries) == 0 {
		empty := lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorSubtle).
			Padding(1, 4).
			MarginLeft(2).
			Render(
				SubtleStyle.Render("  No entries yet.  ") +
					LabelStyle.Render("a") +
					SubtleStyle.Render("  to add one"),
			)
		return empty + "\n"
	}

	var b strings.Builder

	siteHeader := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render(fmt.Sprintf("  %-32s", "SITE"))
	userHeader := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render(fmt.Sprintf("%-24s", "USERNAME"))
	otpHeader := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render(" 2FA")
	b.WriteString(siteHeader + userHeader + otpHeader + "\n")
	b.WriteString(lipgloss.NewStyle().Foreground(ColorSubtle).Render("  "+strings.Repeat("─", 62)) + "\n")

	for i, e := range m.entries {
		site := truncate(e.Site, 30)
		user := truncate(e.Username, 22)

		if i == m.cursor {
			cursor := lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render("❯ ")
			siteCol := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render(fmt.Sprintf("%-32s", site))
			userCol := lipgloss.NewStyle().Foreground(ColorOnSurface).Render(fmt.Sprintf("%-24s", user))
			badge := "    "
			if e.TOTP != "" {
				badge = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render(" ◆  ")
			}
			row := lipgloss.NewStyle().
				Background(ColorHighlight).
				Render(cursor + siteCol + userCol + badge)
			b.WriteString(row + "\n")
		} else {
			siteCol := lipgloss.NewStyle().Foreground(ColorOnSurface).Render(fmt.Sprintf("  %-32s", site))
			userCol := lipgloss.NewStyle().Foreground(ColorMuted).Render(fmt.Sprintf("%-24s", user))
			badge := ""
			if e.TOTP != "" {
				badge = lipgloss.NewStyle().Foreground(ColorSuccess).Render(" ◆")
			}
			b.WriteString(siteCol + userCol + badge + "\n")
		}
	}

	return b.String()
}

func (m ListModel) viewDetail() string {
	if m.cursor >= len(m.entries) {
		return ""
	}
	e := m.entries[m.cursor]

	siteHeader := lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true).Render("  " + e.Site)
	sep := lipgloss.NewStyle().Foreground(ColorSubtle).Render("  " + strings.Repeat("─", 40))

	userRow := LabelStyle.Render("  Username  ") +
		lipgloss.NewStyle().Foreground(ColorAccent).Render("◆  ") +
		ValueStyle.Render(e.Username)

	passRow := LabelStyle.Render("  Password  ") +
		lipgloss.NewStyle().Foreground(ColorAccent).Render("◆  ") +
		lipgloss.NewStyle().Foreground(ColorWarning).Render(e.Password)

	rows := []string{siteHeader, sep, userRow, passRow}

	if e.TOTP != "" {
		code, err := totp.GenerateCode(e.TOTP, time.Now())
		if err == nil {
			totpRow := LabelStyle.Render("  TOTP      ") +
				lipgloss.NewStyle().Foreground(ColorSuccess).Render("◆  ") +
				lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true).Render(code)
			rows = append(rows, totpRow)
		}
	}

	if e.Notes != "" {
		notesRow := LabelStyle.Render("  Notes     ") +
			lipgloss.NewStyle().Foreground(ColorAccent).Render("◆  ") +
			SubtleStyle.Render(e.Notes)
		rows = append(rows, notesRow)
	}

	if !e.UpdatedAt.IsZero() {
		updatedRow := lipgloss.NewStyle().Foreground(ColorSubtle).Render(
			fmt.Sprintf("  Updated    ◆  %s", e.UpdatedAt.Format("2006-01-02  15:04")),
		)
		rows = append(rows, updatedRow)
	}

	content := strings.Join(rows, "\n")

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(ColorAccent).
		MarginLeft(2).
		MarginTop(1).
		Padding(1, 2).
		Render(content)
}

func (m ListModel) viewHelp() string {
	sep := lipgloss.NewStyle().Foreground(ColorSubtle).Render("  │  ")

	bind := func(key, desc string) string {
		return lipgloss.NewStyle().Foreground(ColorAccent).Bold(true).Render(key) +
			lipgloss.NewStyle().Foreground(ColorMuted).Render("  "+desc)
	}

	return HelpStyle.Render("  ") +
		bind("↑↓ / jk", "navigate") +
		sep +
		bind("enter", "view") +
		sep +
		bind("a", "add") +
		sep +
		bind("e", "edit") +
		sep +
		bind("d", "delete") +
		sep +
		bind("q", "quit")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-2] + ".."
}

func RunList(entries []models.Entry, startCursor int, status string) (action string, idx int, err error) {
	m := NewListModel(entries, startCursor, status)
	p := tea.NewProgram(m, tea.WithAltScreen())

	final, runErr := p.Run()
	if runErr != nil {
		return "", 0, runErr
	}

	result := final.(ListModel)
	return result.action, result.selected, nil
}
