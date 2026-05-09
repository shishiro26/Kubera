package ui

import (
	"fmt"
	"strings"

	"github.com/boombuler/barcode/qr"
	"github.com/charmbracelet/lipgloss"
)

func PrintQR(content string) error {
	code, err := qr.Encode(content, qr.M, qr.Auto)
	if err != nil {
		return fmt.Errorf("failed to generate QR code: %w", err)
	}

	bounds := code.Bounds()
	w := bounds.Max.X
	h := bounds.Max.Y
	pad := 2

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#FFFFFF"))

	isDark := func(x, y int) bool {
		if x < 0 || x >= w || y < 0 || y >= h {
			return false
		}
		r, _, _, _ := code.At(x, y).RGBA()
		return r < 0x8000
	}

	emptyRow := style.Render(strings.Repeat(" ", w+pad*2))

	for i := 0; i < pad; i++ {
		fmt.Println(emptyRow)
	}

	for y := 0; y < h; y += 2 {
		var sb strings.Builder
		sb.WriteString(strings.Repeat(" ", pad))
		for x := 0; x < w; x++ {
			top := isDark(x, y)
			bot := isDark(x, y+1)
			switch {
			case top && bot:
				sb.WriteRune('█')
			case top:
				sb.WriteRune('▀')
			case bot:
				sb.WriteRune('▄')
			default:
				sb.WriteRune(' ')
			}
		}
		sb.WriteString(strings.Repeat(" ", pad))
		fmt.Println(style.Render(sb.String()))
	}

	for i := 0; i < pad; i++ {
		fmt.Println(emptyRow)
	}

	return nil
}
