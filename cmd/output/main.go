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

func Error(a any) {
	fmt.Println(ERROR.Render(fmt.Sprintf("%s", a)))
}

func Success(a any) {
	fmt.Println(GOOD.Render(fmt.Sprintf("%s", a)))
}

func Normal(a any) {
	fmt.Printf("%s\n", a)
}

func Heading(a any) {
	fmt.Println(lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("%s", a)))
}

func BannerWrap(s string) string {
	return "\n" + lipgloss.NewStyle().Padding(1).Foreground(lipgloss.Color("#00FF00")).Border(lipgloss.DoubleBorder()).Render(BANNER) + "\n" + s
}
