package scaffold

import (
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Primary   lipgloss.TerminalColor
	Secondary lipgloss.TerminalColor
	Success   lipgloss.TerminalColor
	Error     lipgloss.TerminalColor
	Warning   lipgloss.TerminalColor
	Text      lipgloss.TerminalColor
	Muted     lipgloss.TerminalColor
	Border    lipgloss.TerminalColor
}

func DetectTheme() Theme {
	isDark := isDarkBackground()

	if isDark {
		return Theme{
			Primary:   lipgloss.Color("13"), // Bright Magenta
			Secondary: lipgloss.Color("14"), // Bright Cyan
			Success:   lipgloss.Color("10"), // Bright Green
			Error:     lipgloss.Color("9"),  // Bright Red
			Warning:   lipgloss.Color("11"), // Bright Yellow
			Text:      lipgloss.Color("15"), // Bright White
			Muted:     lipgloss.Color("8"),  // Bright Black (Gray)
			Border:    lipgloss.Color("8"),  // Bright Black (Gray)
		}
	}

	return Theme{
		Primary:   lipgloss.Color("5"), // Magenta
		Secondary: lipgloss.Color("6"), // Cyan
		Success:   lipgloss.Color("2"), // Green
		Error:     lipgloss.Color("1"), // Red
		Warning:   lipgloss.Color("3"), // Yellow
		Text:      lipgloss.Color("0"), // Black
		Muted:     lipgloss.Color("8"), // Bright Black (Gray)
		Border:    lipgloss.Color("8"), // Bright Black (Gray)
	}
}

func isDarkBackground() bool {
	bgColor := os.Getenv("COLORFGBG")

	if bgColor != "" {
		parts := strings.Split(bgColor, ";")
		if len(parts) > 0 {
			if num, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
				return num <= 7
			}
		}
	}

	bgColorEnv := os.Getenv("BACKGROUND")
	if bgColorEnv == "dark" {
		return true
	}
	if bgColorEnv == "light" {
		return false
	}

	if os.Getenv("ITERM_PROFILE") != "" {
		return true
	}

	termProgram := os.Getenv("TERM_PROGRAM")
	if termProgram == "Apple_Terminal" || termProgram == "iTerm.app" {
		return true
	}

	colorterm := os.Getenv("COLORTERM")
	if colorterm == "truecolor" || colorterm == "24bit" {
		return true
	}

	return true
}

func BuildStyles(theme Theme) Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true).
			MarginTop(1).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Italic(true).
			MarginBottom(1),

		Label: lipgloss.NewStyle().
			Foreground(theme.Text).
			Bold(true),

		Description: lipgloss.NewStyle().
			Foreground(theme.Muted).
			MarginTop(1).
			MarginBottom(1),

		Focused: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),

		Blurred: lipgloss.NewStyle().
			Foreground(theme.Muted),

		InputBase: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border),

		InputFocused: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary),

		ButtonFocused: lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(theme.Primary).
			Padding(0, 2).
			Bold(true).
			MarginRight(1),

		ButtonBlurred: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border).
			Padding(0, 2).
			MarginRight(1),

		Success: lipgloss.NewStyle().
			Foreground(theme.Success).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(theme.Error).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(theme.Warning).
			Bold(true),

		Info: lipgloss.NewStyle().
			Foreground(theme.Secondary),

		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border).
			Padding(1, 2).
			Width(CONTAINER_WIDTH),

		ContainerPrimary: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Primary).
			Padding(1, 2).
			Width(CONTAINER_WIDTH),

		Help: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Italic(true).
			MarginTop(1),

		ProgressDone: lipgloss.NewStyle().
			Foreground(theme.Success),

		ProgressTodo: lipgloss.NewStyle().
			Foreground(theme.Muted),

		ProgressActive: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),

		Divider: lipgloss.NewStyle().
			Foreground(theme.Border),
	}
}

type Styles struct {
	Title            lipgloss.Style
	Subtitle         lipgloss.Style
	Label            lipgloss.Style
	Description      lipgloss.Style
	Focused          lipgloss.Style
	Blurred          lipgloss.Style
	InputBase        lipgloss.Style
	InputFocused     lipgloss.Style
	ButtonFocused    lipgloss.Style
	ButtonBlurred    lipgloss.Style
	Success          lipgloss.Style
	Error            lipgloss.Style
	Warning          lipgloss.Style
	Info             lipgloss.Style
	Container        lipgloss.Style
	ContainerPrimary lipgloss.Style
	Help             lipgloss.Style
	ProgressDone     lipgloss.Style
	ProgressTodo     lipgloss.Style
	ProgressActive   lipgloss.Style
	Divider          lipgloss.Style
}
