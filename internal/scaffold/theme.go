package scaffold

import (
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Theme-aware colors based on terminal
type Theme struct {
	// Primary colors from terminal
	Primary   string
	Secondary string
	Success   string
	Error     string
	Warning   string
	Text      string
	Muted     string
	Border    string
}

// DetectTheme detects the terminal color scheme
func DetectTheme() Theme {
	// Detect if using dark or light background
	isDark := isDarkBackground()

	// Return appropriate theme
	if isDark {
		return darkTheme
	}
	return lightTheme
}

func isDarkBackground() bool {
	// Check various terminal indicators
	bgColor := os.Getenv("COLORFGBG")

	// Default to dark for most modern terminals
	// Light background typically has high numeric value
	if bgColor != "" {
		parts := strings.Split(bgColor, ";")
		if len(parts) > 0 {
			if num, err := strconv.Atoi(parts[len(parts)-1]); err == nil {
				// Values > 7 typically indicate light background
				return num <= 7
			}
		}
	}

	// Check ITERM_PROFILE or other indicators
	if os.Getenv("ITERM_PROFILE") != "" {
		return true // iTerm defaults to dark
	}

	// Default to dark
	return true
}

// Dark theme - terminal native colors
var darkTheme = Theme{
	Primary:   "5", // Magenta
	Secondary: "6", // Cyan
	Success:   "2", // Green
	Error:     "1", // Red
	Warning:   "3", // Yellow
	Text:      "7", // White
	Muted:     "8", // Bright Black (Gray)
	Border:    "8", // Bright Black
}

// Light theme - terminal native colors for light backgrounds
var lightTheme = Theme{
	Primary:   "4", // Blue
	Secondary: "6", // Cyan
	Success:   "2", // Green
	Error:     "1", // Red
	Warning:   "3", // Yellow
	Text:      "0", // Black
	Muted:     "7", // White
	Border:    "8", // Bright Black
}

// BuildStyles creates responsive styles from theme
func BuildStyles(theme Theme) Styles {
	return Styles{
		// Title styles
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Primary)).
			Bold(true).
			MarginTop(1).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)).
			Italic(true).
			MarginBottom(1),

		// Label styles
		Label: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Text)).
			Bold(true),

		Description: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)).
			MarginTop(1).
			MarginBottom(1),

		// Input styles
		Focused: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Primary)).
			Bold(true),

		Blurred: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)),

		InputBase: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.Border)),

		InputFocused: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.Primary)),

		// Button styles
		ButtonFocused: lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color(theme.Primary)).
			Padding(0, 2).
			Bold(true).
			MarginRight(1),

		ButtonBlurred: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.Border)).
			Padding(0, 2).
			MarginRight(1),

		// Status styles
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Success)).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Error)).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Warning)).
			Bold(true),

		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Primary)),

		// Container styles
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.Border)).
			Padding(1, 2),

		ContainerPrimary: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(theme.Primary)).
			Padding(1, 2),

		// Help text
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)).
			Italic(true).
			MarginTop(1),

		// Progress
		ProgressDone: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Success)),

		ProgressTodo: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Muted)),

		ProgressActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Primary)).
			Bold(true),

		// Divider
		Divider: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Border)),
	}
}

// Styles holds all style definitions
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
