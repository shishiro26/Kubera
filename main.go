package main

import (
	"fmt"
	"os"

	"github.com/shishiro26/kubera/commands"
	"github.com/shishiro26/kubera/ui"
)

// Injected at build time by GoReleaser via ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
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
		err = commands.Init()
	case "add":
		err = commands.Add()
	case "list":
		err = commands.List()
	case "get":
		requireArg(cmd, 3)
		err = commands.Get(os.Args[2])
	case "edit", "update":
		requireArg(cmd, 3)
		err = commands.Edit(os.Args[2])
	case "delete", "remove", "rm":
		requireArg(cmd, 3)
		err = commands.Delete(os.Args[2])
	case "install":
		err = commands.Install()
	case "version", "--version", "-v":
		fmt.Printf("kubera %s (commit %s, built %s)\n", version, commit, date)
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
		{"edit <site>", "Update entry", "Change username, password, or TOTP for a site"},
		{"delete <site>", "Remove entry", "Delete an entry with confirmation prompt"},
		{"install", "System install", "Copy kubera to PATH and make it available everywhere"},
		{"version", "Show version", fmt.Sprintf("Current build: %s (%s)", version, commit)},
		{"help", "Show this help", "Usage: kubera <command> [args]"},
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
}

func requireArg(cmd string, minArgs int) {
	if len(os.Args) < minArgs {
		ui.PrintError(fmt.Sprintf("Usage: kubera %s <site>", cmd))
		os.Exit(1)
	}
}
