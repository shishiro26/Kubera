package ui

import "github.com/charmbracelet/lipgloss"

// Thank you AI for generating this beautiful style palette 
var (
	ColorPrimary   = lipgloss.Color("#D4A017")
	ColorAccent    = lipgloss.Color("#00CED1")
	ColorSuccess   = lipgloss.Color("#3DDC84")
	ColorError     = lipgloss.Color("#FF4C4C")
	ColorWarning   = lipgloss.Color("#FFB347")
	ColorSubtle    = lipgloss.Color("#4A5568")
	ColorMuted     = lipgloss.Color("#718096")
	ColorOnSurface = lipgloss.Color("#E2E8F0")
	ColorSurface   = lipgloss.Color("#1A1E2E")
	ColorHighlight = lipgloss.Color("#2D3250")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(ColorPrimary).
			Padding(0, 2)

	HeaderStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(ColorSubtle)

	LabelStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	ValueStyle = lipgloss.NewStyle().
			Foreground(ColorOnSurface)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			Italic(true)

	BoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(ColorSubtle).
			Padding(1, 2)

	FocusedBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	VaultBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(ColorAccent).
			Padding(1, 3)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Background(ColorHighlight).
			Bold(true).
			PaddingLeft(1)

	NormalStyle = lipgloss.NewStyle().
			Foreground(ColorOnSurface).
			PaddingLeft(1)

	DimmedStyle = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			PaddingLeft(1)

	BadgeSuccessStyle = lipgloss.NewStyle().
				Foreground(ColorSurface).
				Background(ColorSuccess).
				Bold(true).
				Padding(0, 1)

	BadgeErrorStyle = lipgloss.NewStyle().
			Foreground(ColorSurface).
			Background(ColorError).
			Bold(true).
			Padding(0, 1)

	BadgeGoldStyle = lipgloss.NewStyle().
			Foreground(ColorSurface).
			Background(ColorPrimary).
			Bold(true).
			Padding(0, 1)
)
