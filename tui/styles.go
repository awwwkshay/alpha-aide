package main

import "github.com/charmbracelet/lipgloss"

var (
	cyan   = lipgloss.Color("#00AFFF")
	green  = lipgloss.Color("#00D7AF")
	yellow = lipgloss.Color("#FFD700")
	red    = lipgloss.Color("#FF5F5F")
	dim    = lipgloss.Color("#555555")
	bright = lipgloss.Color("#EEEEEE")
	muted  = lipgloss.Color("#888888")

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cyan)

	logoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cyan)

	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(muted)

	activeItemStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	inactiveItemStyle = lipgloss.NewStyle().
				Foreground(bright)

	cursorStyle = lipgloss.NewStyle().
			Foreground(cyan).
			Bold(true)

	userLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(green)

	agentLabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(cyan)

	toolStyle = lipgloss.NewStyle().
			Foreground(yellow)

	toolResultStyle = lipgloss.NewStyle().
			Foreground(dim)

	errorStyle = lipgloss.NewStyle().
			Foreground(red)

	dimStyle = lipgloss.NewStyle().
			Foreground(dim)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(dim)

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(cyan)

	inputPrefixStyle = lipgloss.NewStyle().
				Foreground(green).
				Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(dim)

	tokenStyle = lipgloss.NewStyle().
			Foreground(muted)
)
