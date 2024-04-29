package ui

import "github.com/charmbracelet/lipgloss"

var (
	BASE_STYLE = lipgloss.NewStyle().PaddingLeft(1).
			Foreground(lipgloss.Color("#00FF00")) // Green
	ERROR_STYLE = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#FF0000")) // Red
)
