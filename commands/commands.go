package commands

import (
	"fmt"
	"net/url"
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
	ui.PrintTitle("Add Entry")

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
