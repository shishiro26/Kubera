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
	var err error

	cmd := os.Args[1]
	switch cmd {
	case "init":
		fmt.Print("Initializing vault... ")
	case "add":
		fmt.Print("Adding new entry... ")
	case "list":
		fmt.Print("OPening vault browser... ")
	case "get":
		requireArg(cmd, 3)
		fmt.Print("Fetching credentials for " + os.Args[2] + "... ")
	case "edit", "update":
		requireArg(cmd, 3)
		fmt.Print("Updating entry for " + os.Args[2] + "... ")
	case "delete", "remove", "rm":
		requireArg(cmd, 3)
		fmt.Print("Deleting entry for " + os.Args[2] + "... ")
	case "install":
		fmt.Print("Installing kubera... ")
	case "help", "--help", "-h", "man":
		showHelp()
	default:
		ui.PrintError("Unknown command: " + cmd)
		showHelp()
		os.Exit(1)
	}

	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
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

func requireArg(cmd string, minArgs int) {
	if len(os.Args) < minArgs {
		ui.PrintError(fmt.Sprintf("Usage: kubera %s <site>", cmd))
		os.Exit(1)
	}
}
