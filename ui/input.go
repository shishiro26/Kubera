package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

var stdinReader = bufio.NewReader(os.Stdin)

func ReadPassword(prompt string) (string, error) {
	fmt.Print(LabelStyle.Render(prompt))
	password, err := term.ReadPassword(uintptr(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(password)), nil
}

func ReadLine(prompt string) string {
	fmt.Print(LabelStyle.Render(prompt))
	text, _ := stdinReader.ReadString('\n')
	return strings.TrimSpace(text)
}

func ReadLineOptional(prompt string) string {
	fmt.Print(SubtleStyle.Render(prompt))
	text, _ := stdinReader.ReadString('\n')
	return strings.TrimSpace(text)
}

func Confirm(prompt string) bool {
	fmt.Print(WarningStyle.Render(prompt + " [y/N]: "))
	text, _ := stdinReader.ReadString('\n')
	return strings.ToLower(strings.TrimSpace(text)) == "y"
}

func PrintSuccess(msg string) {
	fmt.Println(SuccessStyle.Render("  ❯❯  " + msg))
}

func PrintError(msg string) {
	fmt.Fprintln(os.Stderr, ErrorStyle.Render("  !!  "+msg))
}

func PrintWarning(msg string) {
	fmt.Println(WarningStyle.Render("  ~~  " + msg))
}

func PrintTitle(title string) {
	fmt.Println()
	fmt.Println(TitleStyle.Render("  " + title + "  "))
	fmt.Println()
}

func PrintBanner() {
	art := []string{
		"██╗  ██╗██╗   ██╗██████╗ ███████╗██████╗  █████╗ ",
		"██║ ██╔╝██║   ██║██╔══██╗██╔════╝██╔══██╗██╔══██╗",
		"█████╔╝ ██║   ██║██████╔╝█████╗  ██████╔╝███████║",
		"██╔═██╗ ██║   ██║██╔══██╗██╔══╝  ██╔══██╗██╔══██║",
		"██║  ██╗╚██████╔╝██████╔╝███████╗██║  ██║██║  ██║",
		"╚═╝  ╚═╝ ╚═════╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝",
	}

	gradientColors := []lipgloss.Color{
		lipgloss.Color("#F5C842"),
		lipgloss.Color("#E8B832"),
		lipgloss.Color("#D4A017"),
		lipgloss.Color("#C49010"),
		lipgloss.Color("#B07B0A"),
		lipgloss.Color("#9C6A05"),
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

	lockIcon := lipgloss.NewStyle().
		Foreground(ColorAccent).
		Render("⬡")

	tagline := lipgloss.NewStyle().Foreground(ColorSubtle).Render("─────────── ") +
		lockIcon +
		lipgloss.NewStyle().Foreground(ColorOnSurface).Italic(true).Render("  Secure Password Vault  ") +
		lockIcon +
		lipgloss.NewStyle().Foreground(ColorSubtle).Render(" ───────────")

	versionDot := lipgloss.NewStyle().Foreground(ColorAccent).Render("◆")

	meta := lipgloss.JoinHorizontal(
		lipgloss.Center,
		lipgloss.NewStyle().Foreground(ColorMuted).Render("v1.0.0  "),
		versionDot,
		lipgloss.NewStyle().Foreground(ColorMuted).Render("  100% Local  "),
		versionDot,
		lipgloss.NewStyle().Foreground(ColorMuted).Render("  AES-256 Encrypted"),
	)

	inner := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(artBlock + "\n\n" + tagline + "\n" + meta)

	termWidth, _, err := term.GetSize(uintptr(syscall.Stdout))
	if err != nil || termWidth <= 0 {
		termWidth = 80
	}

	boxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 3).
		MarginTop(1).
		MarginBottom(1)

	boxWidth := lipgloss.Width(boxStyle.Render(inner))
	leftPad := (termWidth - boxWidth) / 2
	if leftPad < 0 {
		leftPad = 0
	}
	box := boxStyle.MarginLeft(leftPad).Render(inner)

	accentText := "◈  All data stored locally. Nothing leaves your machine.  ◈"
	accentPad := (termWidth - lipgloss.Width(accentText)) / 2
	if accentPad < 0 {
		accentPad = 0
	}
	accent := lipgloss.NewStyle().
		Foreground(ColorAccent).
		MarginLeft(accentPad).
		Render(accentText)

	fmt.Println(box)
	fmt.Println(accent)
	fmt.Println()
}
