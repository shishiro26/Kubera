package commands

import (
	"fmt"

	"github.com/shishiro26/kubera/models"
	"github.com/shishiro26/kubera/storage"
	"github.com/shishiro26/kubera/ui"
)

func Init() error {
	ui.PrintBanner()
	ui.PrintTitle("Initialize vault")

	if storage.Exists() {
		ui.PrintWarning("A vault already exists. The master password cannot be changed.")
		fmt.Println(ui.SubtleStyle.Render(" To Start fresh, delete:" + storage.VaultPath))
		return nil
	}

	fmt.Println(ui.SubtleStyle.Render("Choose a string master password. It cannot be changed later"))
	fmt.Println()
	password, err := ui.ReadPassword(("Master Password:"))
	if err != nil {
		return err
	}

	if len(password) < 8 {
		ui.PrintError("Master Password must be atleast 8 characters")
		return nil
	}

	confirm, err := ui.ReadPassword(" Confirm password ")
	if err != nil {
		return err
	}

	if password != confirm {
		ui.PrintError("Passwords do not match")
		return nil
	}
	fmt.Println()
	if err := storage.Save(password, []models.Entry{}); err != nil {
		fmt.Println()
		return err
	}

	fmt.Println(ui.SuccessStyle.Render("done"))
	fmt.Println()
	ui.PrintSuccess("Vault initialized. Master password is permanent and cannot be changed.")
	fmt.Println()
	return nil
}
