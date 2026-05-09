package commands

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shishiro26/kubera/ui"
)

func Install() error {
	ui.PrintTitle("Install Kubera")

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot locate current binary: %w", err)
	}

	installDir, installPath := resolveInstallPaths()

	fmt.Println(ui.LabelStyle.Render("  Destination: ") + ui.ValueStyle.Render(installPath))
	fmt.Println()

	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("cannot create install directory: %w", err)
	}

	if err := copyBinary(exe, installPath); err != nil {
		return fmt.Errorf("cannot copy binary: %w", err)
	}
	ui.PrintSuccess("Binary copied")

	if err := persistPath(installDir); err != nil {
		fmt.Println()
		ui.PrintWarning("Could not update PATH automatically: " + err.Error())
		fmt.Println(ui.SubtleStyle.Render("  Add this to your shell profile manually:"))
		if runtime.GOOS == "windows" {
			fmt.Println(ui.ValueStyle.Render(`    setx PATH "%PATH%;` + installDir + `"`))
		} else {
			fmt.Println(ui.ValueStyle.Render(`    export PATH="` + installDir + `:$PATH"`))
		}
	} else {
		ui.PrintSuccess("Added " + installDir + " to PATH")
	}

	fmt.Println()
	ui.PrintSuccess("Done! Open a new terminal and run: kubera help")
	fmt.Println()
	return nil
}

func resolveInstallPaths() (dir, path string) {
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		dir = filepath.Join(home, "bin")
		path = filepath.Join(dir, "kubera.exe")
	} else {
		dir = filepath.Join(home, ".local", "bin")
		path = filepath.Join(dir, "kubera")
	}
	return
}

func copyBinary(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func persistPath(dir string) error {
	if runtime.GOOS == "windows" {
		return persistPathWindows(dir)
	}
	return persistPathUnix(dir)
}

func persistPathWindows(dir string) error {
	script := fmt.Sprintf(
		`$p=[Environment]::GetEnvironmentVariable('Path','User');`+
			`if($p -notlike '*%s*'){[Environment]::SetEnvironmentVariable('Path',$p+';%s','User')}`,
		dir, dir,
	)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func persistPathUnix(dir string) error {
	export := fmt.Sprintf(`export PATH="%s:$PATH"`, dir)
	home, _ := os.UserHomeDir()

	profiles := []string{
		filepath.Join(home, ".zshrc"),
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".bash_profile"),
		filepath.Join(home, ".profile"),
	}

	for _, profile := range profiles {
		if _, err := os.Stat(profile); err != nil {
			continue
		}
		data, _ := os.ReadFile(profile)
		if strings.Contains(string(data), dir) {
			return nil // already present
		}
		f, err := os.OpenFile(profile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			continue
		}
		fmt.Fprintf(f, "\n# kubera\n%s\n", export)
		f.Close()
		return nil
	}

	return fmt.Errorf("no shell profile found (~/.zshrc, ~/.bashrc, ~/.profile)")
}
