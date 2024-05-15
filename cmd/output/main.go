package output

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	GOOD   = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#00FF00")) // Green
	ERROR  = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#FF0000")) // Red
	BANNER = `
████████╗███████╗██████╗ ███╗   ███╗███████╗ ██████╗ ██╗     
╚══██╔══╝██╔════╝██╔══██╗████╗ ████║██╔════╝██╔═══██╗██║     
   ██║   █████╗  ██████╔╝██╔████╔██║███████╗██║   ██║██║     
   ██║   ██╔══╝  ██╔══██╗██║╚██╔╝██║╚════██║██║▄▄ ██║██║     
   ██║   ███████╗██║  ██║██║ ╚═╝ ██║███████║╚██████╔╝███████╗
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝ ╚══▀▀═╝ ╚══════╝
   `
)

func Error(s string) {
	fmt.Println(ERROR.Render(s))
}

func Success(s string) {
	fmt.Println(GOOD.Render(s))
}

func BannerWrap(s string) string {
	return lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("#00FF00")).Border(lipgloss.DoubleBorder()).Render(BANNER) + "\n" + s
}
