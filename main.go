package main

import (
	"fmt"
	"os"

	"github.com/shishiro26/kubera/ui"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}
}

func showHelp() {
	ui.PrintBanner()

	cmds := [][3]string{
		{"init", "Initialize vault", "Master password set once, stored at ~/.kubera/vault.enc"},
		{"add", "Add entry", "Add a new password entry interactively"},
		{"list", "Browse vault", "Interactive browser — navigate, view, add, edit, delete"},
		{"get <site>", "Fetch credentials", "Display username & password for a site"},
		{"edit <site>", "Update entry", "Change username or password for a site"},
		{"delete <site>", "Remove entry", "Delete an entry with confirmation prompt"},
		{"install", "System install", "Add kubera to PATH and optional Windows startup"},
	}

	fmt.Println(ui.HeaderStyle.Render("  COMMANDS"))
	fmt.Println()

	for _, c := range cmds {
		cmd := ui.LabelStyle.Render(fmt.Sprintf("  kubera %-20s", c[0]))
		badge := ui.BadgeGoldStyle.Render(fmt.Sprintf(" %-18s", c[1]))
		desc := ui.SubtleStyle.Render("  " + c[2])
		fmt.Println(cmd + badge + desc)
	}

	fmt.Println()

	vaultPath := ui.ValueStyle.Render("~/.kubera/vault.enc")
	fmt.Println(ui.SubtleStyle.Render("  Vault ") + ui.LabelStyle.Render("→ ") + vaultPath)

	fmt.Println()

	fmt.Println(
		ui.HelpStyle.Render("  Run ") +
			ui.BadgeGoldStyle.Render(" kubera <command> --help ") +
			ui.HelpStyle.Render(" for detailed usage."),
	)

	fmt.Println()
}
