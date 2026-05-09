package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/shishiro26/kubera/crypto"
	"github.com/shishiro26/kubera/models"
)

const vaultVersion = 1

var (
	VaultPath       string
	TOTPSecretPath  string
)

func init() {
	usr, _ := user.Current()
	VaultPath = filepath.Join(usr.HomeDir, ".kubera", "vault.enc")
	TOTPSecretPath = filepath.Join(usr.HomeDir, ".kubera", "totp.secret")
}

func SaveTOTPSecret(secret string) error {
	dir := GetVaultDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return os.WriteFile(TOTPSecretPath, []byte(secret), 0600)
}

func LoadTOTPSecret() (string, error) {
	data, err := os.ReadFile(TOTPSecretPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("TOTP secret not found, run 'kubera init' first")
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func GetVaultDir() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".kubera")
}

func Exists() bool {
	_, err := os.Stat(VaultPath)
	return err == nil
}

func Load(password string) ([]models.Entry, error) {
	data, err := os.ReadFile(VaultPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("vault not found, run 'kubera init' to create one")
		}
		return nil, err
	}

	var vault models.Vault
	if err := json.Unmarshal(data, &vault); err != nil {
		return nil, fmt.Errorf("corrupted vault file: ")
	}

	plainText, err := crypto.Decrypt(vault.Salt, vault.CipherText, password)
	if err != nil {
		return nil, err
	}

	var entries []models.Entry
	if err := json.Unmarshal(plainText, &entries); err != nil {
		return nil, fmt.Errorf("corrupted vault data")
	}

	return entries, nil
}

func Save(password string, entries []models.Entry) error {
	data, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	saltHex, cipherText, err := crypto.Encrypt(data, password)
	if err != nil {
		return err
	}

	vault := models.Vault{
		Version:    vaultVersion,
		Salt:       saltHex,
		CipherText: cipherText,
	}

	raw, err := json.Marshal(vault)
	if err != nil {
		return err
	}

	dir := GetVaultDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	return os.WriteFile(VaultPath, raw, 0600)
}