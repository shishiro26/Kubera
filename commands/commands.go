package commands

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/shishiro26/kubera/models"
	"github.com/shishiro26/kubera/storage"
	"github.com/shishiro26/kubera/ui"
)

func Init() error {
	ui.PrintBanner()
	ui.PrintTitle("Initialize Vault")

	if storage.Exists() {
		ui.PrintWarning("A vault already exists. The master password cannot be changed.")
		fmt.Println(ui.SubtleStyle.Render("  To start fresh, delete: " + storage.VaultPath))
		return nil
	}

	fmt.Println(ui.SubtleStyle.Render("  Choose a strong master password. It cannot be changed later."))
	fmt.Println()

	password, err := ui.ReadPassword("Master Password: ")
	if err != nil {
		return err
	}

	if len(password) < 8 {
		ui.PrintError("Master password must be at least 8 characters")
		return nil
	}

	confirm, err := ui.ReadPassword("Confirm Password: ")
	if err != nil {
		return err
	}

	if password != confirm {
		ui.PrintError("Passwords do not match")
		return nil
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Kubera",
		AccountName: "Kubera Vault",
	})
	if err != nil {
		return fmt.Errorf("failed to generate TOTP: %w", err)
	}

	secret := key.Secret()

	if err := storage.SaveTOTPSecret(secret); err != nil {
		return fmt.Errorf("failed to save TOTP secret: %w", err)
	}

	vaultPassword := password + secret
	if err := storage.Save(vaultPassword, []models.Entry{}); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(ui.LabelStyle.Render("  Scan with your authenticator app:"))
	fmt.Println()
	if err := ui.PrintQR(key.URL()); err != nil {
		fmt.Println(ui.SubtleStyle.Render("  (QR render failed, use the URL below)"))
	}
	fmt.Println()
	fmt.Println(ui.LabelStyle.Render("  Manual entry secret:"))
	fmt.Println(ui.ValueStyle.Render("  " + secret))
	fmt.Println()
	ui.PrintSuccess("Vault initialized. You will need your TOTP code to unlock the vault.")
	fmt.Println()
	return nil
}

func Unlock() ([]models.Entry, string, error) {
	if !storage.Exists() {
		return nil, "", fmt.Errorf("no vault found — run 'kubera init' first")
	}

	secret, err := storage.LoadTOTPSecret()
	if err != nil {
		return nil, "", err
	}

	password, err := ui.ReadPassword("Master Password: ")
	if err != nil {
		return nil, "", err
	}

	code, err := ui.ReadPassword("TOTP Code: ")
	if err != nil {
		return nil, "", err
	}

	if !totp.Validate(code, secret) {
		return nil, "", fmt.Errorf("invalid TOTP code")
	}

	vaultPassword := password + secret
	entries, err := storage.Load(vaultPassword)
	if err != nil {
		return nil, "", fmt.Errorf("invalid master password")
	}

	return entries, vaultPassword, nil
}

func Add() error {
	ui.PrintTitle("Add Password")

	entries, vaultPassword, err := Unlock()
	if err != nil {
		return err
	}

	fmt.Println()
	site := ui.ReadLine("Site / Service:      ")
	if site == "" {
		return fmt.Errorf("site cannot be empty")
	}

	username := ui.ReadLine("Username / Email:    ")

	for _, e := range entries {
		if strings.EqualFold(e.Site, site) && strings.EqualFold(e.Username, username) {
			ui.PrintError(fmt.Sprintf("An entry for '%s' with username '%s' already exists.", site, username))
			return nil
		}
	}

	password, err := ui.ReadPassword("Password:            ")
	if err != nil {
		return err
	}

	notes := ui.ReadLineOptional("Notes (optional):    ")

	var totpSecret string
	if ui.Confirm("Add a TOTP secret for this entry?") {
		totpSecret = ui.ReadLine("TOTP Secret (base32): ")
		if totpSecret != "" {
			label := url.PathEscape(site + ":" + username)
			otpURL := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=%s",
				label, totpSecret, url.QueryEscape(site))

			fmt.Println()
			fmt.Println(ui.LabelStyle.Render("  TOTP QR Code for " + site + ":"))
			fmt.Println()
			if err := ui.PrintQR(otpURL); err != nil {
				fmt.Println(ui.SubtleStyle.Render("  (QR render failed)"))
			}
			fmt.Println()
		}
	}

	entry := models.Entry{
		Site:      site,
		Username:  username,
		Password:  password,
		TOTP:      totpSecret,
		Notes:     notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	entries = append(entries, entry)

	if err := storage.Save(vaultPassword, entries); err != nil {
		return err
	}

	ui.PrintSuccess("Entry added: " + site)
	fmt.Println()
	return nil
}

func List() error {
	ui.PrintBanner()
	entries, vaultPassword, err := Unlock()
	if err != nil {
		return err
	}

	fmt.Println()
	cursor := 0
	status := ""

	for {
		action, idx, err := ui.RunList(entries, cursor, status)
		if err != nil {
			return err
		}
		status = ""

		switch action {
		case "delete":
			if idx < len(entries) {
				site := entries[idx].Site
				entries = append(entries[:idx], entries[idx+1:]...)
				if saveErr := storage.Save(vaultPassword, entries); saveErr != nil {
					ui.PrintError("Save failed: " + saveErr.Error())
					break
				}
				status = fmt.Sprintf("'%s' deleted.", site)
				cursor = idx
				if cursor >= len(entries) && cursor > 0 {
					cursor--
				}
				continue
			}

		case "edit":
			if idx < len(entries) {
				entries, err = performEdit(vaultPassword, entries, idx)
				if err != nil {
					ui.PrintError(err.Error())
				} else {
					status = fmt.Sprintf("'%s' updated.", entries[idx].Site)
				}
				cursor = idx
				continue
			}

		case "add":
			var newSite string
			entries, newSite, err = performAdd(vaultPassword, entries)
			if err != nil {
				ui.PrintError(err.Error())
			} else if newSite != "" {
				status = fmt.Sprintf("'%s' added.", newSite)
				cursor = len(entries) - 1
				if cursor < 0 {
					cursor = 0
				}
			}
			continue
		}

		break
	}

	return nil
}

func performAdd(vaultPassword string, entries []models.Entry) ([]models.Entry, string, error) {
	fmt.Println()
	ui.PrintTitle("Add Entry")

	site := ui.ReadLine("  Site / App name: ")
	if site == "" {
		ui.PrintError("Site name cannot be empty.")
		return entries, "", nil
	}

	username := ui.ReadLine("  Username / Email: ")

	for _, e := range entries {
		if strings.EqualFold(e.Site, site) && strings.EqualFold(e.Username, username) {
			ui.PrintError(fmt.Sprintf("'%s' with username '%s' already exists.", site, username))
			return entries, "", nil
		}
	}
	entryPass, err := ui.ReadPassword("  Password: ")
	if err != nil {
		return entries, "", err
	}
	notes := ui.ReadLineOptional("  Notes (optional): ")

	var totpSecret string
	if ui.Confirm("  Add a TOTP secret for this entry?") {
		totpSecret = ui.ReadLine("  TOTP Secret (base32): ")
		if totpSecret != "" {
			label := url.PathEscape(site + ":" + username)
			otpURL := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=%s",
				label, totpSecret, url.QueryEscape(site))
			fmt.Println()
			if err := ui.PrintQR(otpURL); err != nil {
				fmt.Println(ui.SubtleStyle.Render("  (QR render failed)"))
			}
			fmt.Println()
		}
	}

	now := time.Now()
	entries = append(entries, models.Entry{
		Site:      site,
		Username:  username,
		Password:  entryPass,
		TOTP:      totpSecret,
		Notes:     notes,
		CreatedAt: now,
		UpdatedAt: now,
	})

	if err := storage.Save(vaultPassword, entries); err != nil {
		return entries, "", err
	}

	fmt.Println()
	return entries, site, nil
}

func performEdit(vaultPassword string, entries []models.Entry, idx int) ([]models.Entry, error) {
	e := entries[idx]

	fmt.Println()
	ui.PrintTitle(fmt.Sprintf("Edit — %s", e.Site))
	fmt.Println(ui.SubtleStyle.Render("  Press Enter to keep the current value."))
	fmt.Println()

	newUsername := ui.ReadLine(fmt.Sprintf("  Username [%s]: ", e.Username))
	if newUsername == "" {
		newUsername = e.Username
	}

	if !strings.EqualFold(newUsername, e.Username) {
		for i, other := range entries {
			if i != idx && strings.EqualFold(other.Site, e.Site) && strings.EqualFold(other.Username, newUsername) {
				return entries, fmt.Errorf("an entry for '%s' with username '%s' already exists", e.Site, newUsername)
			}
		}
	}

	newPass, err := ui.ReadPassword("  New password (Enter to keep): ")
	if err != nil {
		return entries, err
	}
	if newPass == "" {
		newPass = e.Password
	}

	newNotes := ui.ReadLineOptional(fmt.Sprintf("  Notes [%s] (Enter to keep): ", e.Notes))
	if newNotes == "" {
		newNotes = e.Notes
	}

	newTOTP := e.TOTP
	if ui.Confirm("  Update TOTP secret? (Enter to skip)") {
		input := ui.ReadLine("  New TOTP Secret (base32, Enter to clear): ")
		newTOTP = input
	}

	entries[idx].Username = newUsername
	entries[idx].Password = newPass
	entries[idx].Notes = newNotes
	entries[idx].TOTP = newTOTP
	entries[idx].UpdatedAt = time.Now()

	if err := storage.Save(vaultPassword, entries); err != nil {
		return entries, err
	}

	fmt.Println()
	return entries, nil
}

func Get(site string) error {
	ui.PrintTitle("Get Password")

	entries, _, err := Unlock()
	if err != nil {
		return err
	}
	fmt.Println()

	found := false
	for _, e := range entries {
		if strings.EqualFold(e.Site, site) {
			found = true
			content := ui.LabelStyle.Render("Site:     ") + ui.ValueStyle.Render(e.Site) + "\n" +
				ui.LabelStyle.Render("Username: ") + ui.ValueStyle.Render(e.Username) + "\n" +
				ui.LabelStyle.Render("Password: ") + ui.ValueStyle.Render(e.Password)
			if e.Notes != "" {
				content += "\n" + ui.LabelStyle.Render("Notes:    ") + ui.ValueStyle.Render(e.Notes)
			}
			if e.TOTP != "" {
				code, totpErr := totp.GenerateCode(e.TOTP, time.Now())
				if totpErr == nil {
					content += "\n" + ui.LabelStyle.Render("TOTP:     ") + ui.ValueStyle.Render(code)
				}
			}
			fmt.Println(ui.BoxStyle.Render(content))
		}
	}

	if !found {
		ui.PrintError(fmt.Sprintf("No entry found for '%s'.", site))
	}
	return nil
}

func Edit(site string) error {
	ui.PrintTitle(fmt.Sprintf("Edit — %s", site))

	entries, vaultPassword, err := Unlock()
	if err != nil {
		return err
	}
	fmt.Println()

	var matches []int
	for i, e := range entries {
		if strings.EqualFold(e.Site, site) {
			matches = append(matches, i)
		}
	}

	if len(matches) == 0 {
		ui.PrintError(fmt.Sprintf("No entry found for '%s'.", site))
		return nil
	}

	idx := matches[0]
	if len(matches) > 1 {
		fmt.Println(ui.SubtleStyle.Render(fmt.Sprintf("  Multiple accounts for '%s':", site)))
		fmt.Println()
		for i, mi := range matches {
			fmt.Println(ui.LabelStyle.Render(fmt.Sprintf("  %d. ", i+1)) + ui.ValueStyle.Render(entries[mi].Username))
		}
		fmt.Println()
		input := ui.ReadLine("  Select account (number): ")
		n, parseErr := strconv.Atoi(strings.TrimSpace(input))
		if parseErr != nil || n < 1 || n > len(matches) {
			ui.PrintError("Invalid selection.")
			return nil
		}
		idx = matches[n-1]
	}

	if _, err := performEdit(vaultPassword, entries, idx); err != nil {
		ui.PrintError(err.Error())
		return nil
	}
	ui.PrintSuccess(fmt.Sprintf("'%s' updated.", entries[idx].Site))
	return nil
}

func Delete(site string) error {
	ui.PrintTitle(fmt.Sprintf("Delete — %s", site))

	entries, vaultPassword, err := Unlock()
	if err != nil {
		return err
	}
	fmt.Println()

	var matches []int
	for i, e := range entries {
		if strings.EqualFold(e.Site, site) {
			matches = append(matches, i)
		}
	}

	if len(matches) == 0 {
		ui.PrintError(fmt.Sprintf("No entry found for '%s'.", site))
		return nil
	}

	idx := matches[0]
	if len(matches) > 1 {
		fmt.Println(ui.SubtleStyle.Render(fmt.Sprintf("  Multiple accounts for '%s':", site)))
		fmt.Println()
		for i, mi := range matches {
			fmt.Println(ui.LabelStyle.Render(fmt.Sprintf("  %d. ", i+1)) + ui.ValueStyle.Render(entries[mi].Username))
		}
		fmt.Println()
		input := ui.ReadLine("  Select account (number): ")
		n, parseErr := strconv.Atoi(strings.TrimSpace(input))
		if parseErr != nil || n < 1 || n > len(matches) {
			ui.PrintError("Invalid selection.")
			return nil
		}
		idx = matches[n-1]
	}

	e := entries[idx]
	if !ui.Confirm(fmt.Sprintf("  Delete '%s' @ %s?", e.Username, e.Site)) {
		fmt.Println(ui.SubtleStyle.Render("  Aborted."))
		return nil
	}
	entries = append(entries[:idx], entries[idx+1:]...)
	if err := storage.Save(vaultPassword, entries); err != nil {
		return err
	}
	fmt.Println()
	ui.PrintSuccess(fmt.Sprintf("'%s' @ %s deleted.", e.Username, e.Site))
	return nil
}
