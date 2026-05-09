# Kubera

- An local first encrypted password manager. Passwords that never leave your machine, and are only decrypted in memory when you need them.

- Used Argon2ID for key derivation, and AES-256-GCM for encryption.
- And used TOTP for just overengineering the authentication process, and to make it more secure. and just learn it.

- What did I learn from this:
  - How to use Argon2ID and AES-256-GCM in Go.

## Installation

**Windows**

```powershell
irm https://raw.githubusercontent.com/shishiro26/Kubera/main/install.ps1 | iex
```

**Linux / macOS**

```sh
curl -fsSL https://raw.githubusercontent.com/shishiro26/Kubera/main/install.sh | sh
```

Once installed, run `kubera help` to see available commands.

## Build from source

```sh
git clone https://github.com/shishiro26/Kubera.git
cd Kubera
go build -o kubera .
./kubera help
```
