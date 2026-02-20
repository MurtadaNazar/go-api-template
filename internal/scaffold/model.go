package scaffold

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type State int

const (
	StateMainMenu State = iota
	StateWelcome
	StateProjectName
	StateModuleName
	StateProjectPath
	StateFeatures
	StateEnvVars
	StateConfirm
	StateProcessing
	StateSuccess
	StateError
	StateDevelopmentMenu
	StateBuildTestMenu
	StateCodeQualityMenu
	StateDepsMenu
)

type Feature struct {
	Name        string
	Description string
	Selected    bool
	Default     bool
}

// Fixed container width for consistent layout - account for borders (2) + padding (4)
const CONTAINER_WIDTH = 70

type MenuItem struct {
	Label       string
	Description string
	Action      string // Identifier for the action
}

type Model struct {
	// State
	state       State
	projectName string
	moduleName  string
	projectPath string

	// Environment Variables
	envVars    map[string]string
	envFocus   int
	envEditing bool
	envInput   textinput.Model

	// UI Components
	inputs     []textinput.Model
	focusIndex int
	spinner    spinner.Model
	width      int
	height     int

	// Theme
	styles Styles

	// Features
	features     []Feature
	featureFocus int

	// Menu
	menuItems   []MenuItem
	menuFocus   int
	currentMenu string // Tracks which menu we're in

	// Messages
	err     error
	message string
	warning string

	// Validation
	projectNameValid bool
	moduleNameValid  bool
	projectPathValid bool

	// Feature dependencies
	featureDependencies map[string][]string
}

func NewModel() *Model {
	// Detect theme from terminal
	theme := DetectTheme()
	styles := BuildStyles(theme)

	inputs := make([]textinput.Model, 3)

	// Project name input
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "my-awesome-project"
	inputs[0].CharLimit = 50
	inputs[0].PromptStyle = styles.Focused
	inputs[0].TextStyle = styles.Focused
	inputs[0].PlaceholderStyle = styles.Blurred
	inputs[0].Cursor.Style = styles.Focused
	inputs[0].Focus()
	inputs[0].Width = CONTAINER_WIDTH - 12

	// Module name input
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "github.com/org/my-project"
	inputs[1].CharLimit = 100
	inputs[1].PromptStyle = styles.Blurred
	inputs[1].TextStyle = styles.Blurred
	inputs[1].PlaceholderStyle = styles.Blurred
	inputs[1].Cursor.Style = styles.Blurred
	inputs[1].Width = CONTAINER_WIDTH - 12

	// Project path input
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "."
	inputs[2].CharLimit = 200
	inputs[2].PromptStyle = styles.Blurred
	inputs[2].TextStyle = styles.Blurred
	inputs[2].PlaceholderStyle = styles.Blurred
	inputs[2].Cursor.Style = styles.Blurred
	inputs[2].Width = CONTAINER_WIDTH - 12

	// Environment variable input
	envInput := textinput.New()
	envInput.Placeholder = "value"
	envInput.CharLimit = 200
	envInput.PromptStyle = styles.Focused
	envInput.TextStyle = styles.Focused
	envInput.PlaceholderStyle = styles.Blurred
	envInput.Cursor.Style = styles.Focused
	envInput.Width = CONTAINER_WIDTH - 20

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.Info

	// Feature dependencies
	featureDependencies := map[string][]string{
		"User Management":      {"Authentication (JWT)"},
		"File Storage":         {"Database"},
		"Authentication (JWT)": {},
		"Database":             {},
		"API Docs":             {},
		"Docker":               {},
		"Podman":               {},
	}

	// Initialize features
	features := []Feature{
		{
			Name:        "Authentication (JWT)",
			Description: "JWT-based auth with token rotation",
			Selected:    true,
			Default:     true,
		},
		{
			Name:        "User Management",
			Description: "User registration, profiles, RBAC",
			Selected:    true,
			Default:     true,
		},
		{
			Name:        "Database",
			Description: "PostgreSQL integration with migrations",
			Selected:    true,
			Default:     true,
		},
		{
			Name:        "File Storage",
			Description: "MinIO S3-compatible file storage",
			Selected:    true,
			Default:     true,
		},
		{
			Name:        "API Docs",
			Description: "Auto-generated Swagger documentation",
			Selected:    true,
			Default:     true,
		},
		{
			Name:        "Docker",
			Description: "Docker & Docker Compose setup",
			Selected:    true,
			Default:     true,
		},
		{
			Name:        "Podman",
			Description: "Podman & Podman Compose setup",
			Selected:    false,
			Default:     false,
		},
	}

	// Initialize main menu items
	mainMenu := []MenuItem{
		{Label: "Create New Project", Description: "Create project from template with feature selection", Action: "create"},
		{Label: "Help", Description: "View keyboard shortcuts and documentation", Action: "help"},
		{Label: "Exit", Description: "Exit the scaffolder", Action: "exit"},
	}

	return &Model{
		state:               StateMainMenu,
		inputs:              inputs,
		spinner:             s,
		width:               80,
		height:              24,
		styles:              styles,
		focusIndex:          0,
		projectPath:         ".",
		features:            features,
		featureFocus:        0,
		menuItems:           mainMenu,
		menuFocus:           0,
		currentMenu:         "main",
		featureDependencies: featureDependencies,
		envVars:             make(map[string]string),
		envFocus:            0,
		envEditing:          false,
		envInput:            envInput,
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if m.width < 60 {
			m.width = 60
		}
		if m.height < 20 {
			m.height = 20
		}

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyUp:
			if m.state == StateMainMenu {
				m.menuFocus--
				if m.menuFocus < 0 {
					m.menuFocus = len(m.menuItems) - 1
				}
			} else if m.state == StateFeatures {
				m.featureFocus--
				if m.featureFocus < 0 {
					m.featureFocus = len(m.features) - 1
				}
			} else if m.state == StateEnvVars && !m.envEditing {
				m.envFocus--
				if m.envFocus < 0 {
					m.envFocus = 7
				}
			}

		case tea.KeyDown:
			if m.state == StateMainMenu {
				m.menuFocus++
				if m.menuFocus >= len(m.menuItems) {
					m.menuFocus = 0
				}
			} else if m.state == StateFeatures {
				m.featureFocus++
				if m.featureFocus >= len(m.features) {
					m.featureFocus = 0
				}
			} else if m.state == StateEnvVars && !m.envEditing {
				m.envFocus++
				if m.envFocus > 7 {
					m.envFocus = 0
				}
			}

		case tea.KeyEscape:
			if m.state == StateEnvVars && m.envEditing {
				m.envEditing = false
				m.envInput.Reset()
				return m, nil
			}

		case tea.KeySpace:
			if m.state == StateFeatures {
				feature := &m.features[m.featureFocus]

				if feature.Selected {
					// Deselecting - check if other features depend on this
					feature.Selected = false
					m.checkDependents()
				} else {
					// Selecting - auto-enable dependencies
					feature.Selected = true
					m.enableDependencies(feature.Name)
					m.checkDependents()
				}
				m.warning = m.getDependencyWarning()
			}

		case tea.KeyTab:
			if m.state == StateProjectName || m.state == StateModuleName || m.state == StateProjectPath {
				m.focusIndex++
				maxInputs := 3
				if m.state == StateProjectName {
					maxInputs = 1
				} else if m.state == StateModuleName {
					maxInputs = 2
				}
				if m.focusIndex >= maxInputs {
					m.focusIndex = 0
				}
				m.updateInputFocus()
				if m.focusIndex < len(m.inputs) {
					return m, m.inputs[m.focusIndex].Focus()
				}
			} else if m.state == StateEnvVars && !m.envEditing {
				m.state = StateConfirm
				return m, nil
			}

		case tea.KeyShiftTab:
			if m.state == StateProjectName || m.state == StateModuleName || m.state == StateProjectPath {
				m.focusIndex--
				if m.focusIndex < 0 {
					maxInputs := 3
					if m.state == StateProjectName {
						maxInputs = 1
					} else if m.state == StateModuleName {
						maxInputs = 2
					}
					m.focusIndex = maxInputs - 1
				}
				m.updateInputFocus()
				if m.focusIndex < len(m.inputs) {
					return m, m.inputs[m.focusIndex].Focus()
				}
			}

		case tea.KeyEnter:
			switch m.state {
			case StateMainMenu:
				// Handle main menu selection
				if m.menuFocus < len(m.menuItems) {
					action := m.menuItems[m.menuFocus].Action
					switch action {
					case "create":
						m.state = StateWelcome
						return m, nil
					case "help":
						m.state = StateWelcome // Show help screen
						m.message = m.getHelpText()
						m.warning = "Press CTRL+C to return to menu"
						return m, nil
					case "exit":
						return m, tea.Quit
					}
				}

			case StateWelcome:
				m.state = StateProjectName
				m.focusIndex = 0
				m.updateInputFocus()
				return m, m.inputs[0].Focus()

			case StateProjectName:
				m.projectName = strings.TrimSpace(m.inputs[0].Value())
				if m.projectName == "" {
					m.err = fmt.Errorf("project name is required")
					m.state = StateError
					return m, nil
				}
				if !isValidProjectName(m.projectName) {
					m.err = fmt.Errorf("invalid format: use lowercase, numbers, hyphens, underscores")
					m.state = StateError
					return m, nil
				}
				m.projectNameValid = true
				m.state = StateModuleName
				m.focusIndex = 1
				m.inputs[1].Reset()
				m.updateInputFocus()
				return m, m.inputs[1].Focus()

			case StateModuleName:
				m.moduleName = strings.TrimSpace(m.inputs[1].Value())
				if m.moduleName == "" {
					m.moduleName = fmt.Sprintf("github.com/example/%s", m.projectName)
				}
				if !isValidModuleName(m.moduleName) {
					m.err = fmt.Errorf("invalid module format: use 'domain.com/org/project'")
					m.state = StateError
					return m, nil
				}
				m.moduleNameValid = true
				m.state = StateProjectPath
				m.focusIndex = 2
				m.inputs[2].Reset()
				m.inputs[2].SetValue(".")
				m.updateInputFocus()
				return m, m.inputs[2].Focus()

			case StateProjectPath:
				m.projectPath = strings.TrimSpace(m.inputs[2].Value())
				if m.projectPath == "" {
					m.projectPath = "."
				}
				if !isValidPath(m.projectPath) {
					m.err = fmt.Errorf("invalid path: use relative path like '.' or './projects'")
					m.state = StateError
					return m, nil
				}
				m.projectPathValid = true
				m.state = StateFeatures
				m.featureFocus = 0
				return m, nil

			case StateFeatures:
				m.state = StateEnvVars
				m.envFocus = 0
				return m, nil

			case StateEnvVars:
				if m.envEditing {
					envKeys := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "JWT_SECRET", "MINIO_ACCESS_KEY", "MINIO_SECRET_KEY"}
					if m.envFocus < len(envKeys) {
						m.envVars[envKeys[m.envFocus]] = m.envInput.Value()
					}
					m.envEditing = false
					m.envInput.Reset()
				} else {
					envKeys := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "JWT_SECRET", "MINIO_ACCESS_KEY", "MINIO_SECRET_KEY"}
					defaults := map[string]string{
						"DB_HOST":          "localhost",
						"DB_PORT":          "5432",
						"DB_USER":          "postgres",
						"DB_PASSWORD":      "postgres",
						"DB_NAME":          m.projectName,
						"JWT_SECRET":       "your-secret-key-change-in-production",
						"MINIO_ACCESS_KEY": "minioadmin",
						"MINIO_SECRET_KEY": "minioadmin",
					}
					if m.envFocus < len(envKeys) {
						currentValue := defaults[envKeys[m.envFocus]]
						if v, ok := m.envVars[envKeys[m.envFocus]]; ok {
							currentValue = v
						}
						m.envInput.SetValue(currentValue)
						m.envEditing = true
						return m, m.envInput.Focus()
					}
				}
				return m, nil

			case StateConfirm:
				m.state = StateProcessing
				return m, tea.Batch(m.spinner.Tick, m.processScaffold())

			case StateError:
				m.state = StateWelcome
				m.inputs[0].Reset()
				m.inputs[1].Reset()
				m.inputs[2].Reset()
				m.projectName = ""
				m.moduleName = ""
				m.projectPath = "."
				m.projectNameValid = false
				m.moduleNameValid = false
				m.projectPathValid = false
				m.err = nil
				m.focusIndex = 0
				return m, nil

			case StateSuccess:
				return m, tea.Quit
			}

		default:
			// Handle 'q' key for exit on final screens
			if msg.Type == tea.KeyRunes && len(msg.Runes) == 1 && msg.Runes[0] == 'q' {
				if m.state == StateSuccess || m.state == StateError {
					return m, tea.Quit
				}
			}

			// Handle other key inputs in input states
			if m.state == StateProjectName || m.state == StateModuleName || m.state == StateProjectPath {
				return m, m.updateInputs(msg)
			}

			// Handle env input editing
			if m.state == StateEnvVars && m.envEditing {
				var cmd tea.Cmd
				m.envInput, cmd = m.envInput.Update(msg)
				return m, cmd
			}
		}

	case spinner.TickMsg:
		if m.state == StateProcessing {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case ProcessCompleteMsg:
		if msg.Err != nil {
			m.err = msg.Err
			m.state = StateError
			return m, nil
		}
		m.message = msg.Message
		m.state = StateSuccess
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) updateInputFocus() {
	for i := 0; i < len(m.inputs); i++ {
		if i == m.focusIndex {
			m.inputs[i].PromptStyle = m.styles.Focused
			m.inputs[i].TextStyle = m.styles.Focused
			m.inputs[i].Cursor.Style = m.styles.Focused
			m.inputs[i].PlaceholderStyle = m.styles.Blurred
		} else {
			m.inputs[i].PromptStyle = m.styles.Blurred
			m.inputs[i].TextStyle = m.styles.Blurred
			m.inputs[i].Cursor.Style = m.styles.Blurred
			m.inputs[i].PlaceholderStyle = m.styles.Blurred
			m.inputs[i].Blur()
		}
	}
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	if m.focusIndex < len(m.inputs) {
		m.inputs[m.focusIndex], _ = m.inputs[m.focusIndex].Update(msg)
	}
	return nil
}

func (m *Model) View() string {
	switch m.state {
	case StateMainMenu:
		return m.viewMainMenu()
	case StateWelcome:
		if m.message != "" {
			// Help screen
			return m.viewHelp()
		}
		return m.viewWelcome()
	case StateProjectName:
		return m.viewProjectName()
	case StateModuleName:
		return m.viewModuleName()
	case StateProjectPath:
		return m.viewProjectPath()
	case StateFeatures:
		return m.viewFeatures()
	case StateEnvVars:
		return m.viewEnvVars()
	case StateConfirm:
		return m.viewConfirm()
	case StateProcessing:
		return m.viewProcessing()
	case StateSuccess:
		return m.viewSuccess()
	case StateError:
		return m.viewError()
	default:
		return ""
	}
}

func (m *Model) viewMainMenu() string {
	logo := m.styles.Info.Render(`██████╗  ██████╗     ██████╗ ██╗     ██╗   ██╗███████╗
██╔════╝ ██╔═══██╗    ██╔══██╗██║     ██║   ██║██╔════╝
██║  ███╗██║   ██║    ██████╔╝██║     ██║   ██║███████║
██║   ██║██║   ██║    ██╔═══╝ ██║     ██║   ██║╚════██║
╚██████╔╝╚██████╔╝    ██║     ███████╗╚██████╔╝███████║
 ╚═════╝  ╚═════╝     ╚═╝     ╚══════╝ ╚═════╝ ╚══════╝`)

	title := m.styles.Title.Render("Go Platform Template")
	subtitle := m.styles.Subtitle.Render("Production-Ready Go API Framework")

	// Build menu items
	var menuLines []string
	menuLines = append(menuLines, m.styles.Info.Render("Choose an option:"), "")

	for i, item := range m.menuItems {
		var line string
		if i == m.menuFocus {
			cursor := m.styles.Focused.Render("▸")
			label := m.styles.Focused.Render(item.Label)
			desc := m.styles.Blurred.Render("    " + item.Description)
			line = fmt.Sprintf("  %s %s\n%s", cursor, label, desc)
		} else {
			label := m.styles.Blurred.Render(item.Label)
			desc := m.styles.Blurred.Render("    " + item.Description)
			line = fmt.Sprintf("    %s\n%s", label, desc)
		}
		menuLines = append(menuLines, line)
	}

	menu := lipgloss.JoinVertical(lipgloss.Left, menuLines...)
	instructions := m.styles.Subtitle.Render("↑/↓ Navigate • ENTER Select • CTRL+C Quit")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		"",
		title,
		subtitle,
		"",
		menu,
		"",
		instructions,
	)

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m *Model) viewHelp() string {
	logo := m.styles.Info.Render(`██████╗  ██████╗     ██████╗ ██╗     ██╗   ██╗███████╗
██╔════╝ ██╔═══██╗    ██╔══██╗██║     ██║   ██║██╔════╝
██║  ███╗██║   ██║    ██████╔╝██║     ██║   ██║███████║
██║   ██║██║   ██║    ██╔═══╝ ██║     ██║   ██║╚════██║
╚██████╔╝╚██████╔╝    ██║     ███████╗╚██████╔╝███████║
 ╚═════╝  ╚═════╝     ╚═╝     ╚══════╝ ╚═════╝ ╚══════╝`)

	title := m.styles.Title.Render("Help & Keyboard Shortcuts")

	// Format help text for nice display
	helpBox := m.styles.Blurred.Render(m.message)

	footer := m.styles.Subtitle.Render("(Press CTRL+C to return to menu)")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		logo,
		"",
		title,
		"",
		helpBox,
		"",
		footer,
	)

	return m.padContent(content)
}

func (m *Model) viewWelcome() string {
	logo := m.styles.Info.Render(`██████╗  ██████╗     ██████╗ ██╗     ██╗   ██╗███████╗
██╔════╝ ██╔═══██╗    ██╔══██╗██║     ██║   ██║██╔════╝
██║  ███╗██║   ██║    ██████╔╝██║     ██║   ██║███████║
██║   ██║██║   ██║    ██╔═══╝ ██║     ██║   ██║╚════██║
╚██████╔╝╚██████╔╝    ██║     ███████╗╚██████╔╝███████║
 ╚═════╝  ╚═════╝     ╚═╝     ╚══════╝ ╚═════╝ ╚══════╝`)

	title := m.styles.Title.Render("Go Platform Template Scaffolder")
	subtitle := m.styles.Subtitle.Render("Create a new production-ready Go project")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		"",
		title,
		subtitle,
		"",
		m.styles.Description.Render("This wizard will help you set up a new Go project from the platform template."),
		m.styles.Description.Render("You'll be guided through a few simple steps."),
		"",
		m.styles.Help.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"),
		m.styles.Help.Render("Press ENTER to begin • CTRL+C to exit"),
		m.styles.Help.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"),
	)

	return m.padContent(content)
}

func (m *Model) viewProjectName() string {
	header := m.renderHeader("Project Name", 1, 6)

	input := m.renderInputField(0)

	hint := ""
	value := m.inputs[0].Value()
	if value != "" {
		if isValidProjectName(value) {
			hint = m.styles.Success.Render("✓ Valid project name")
		} else {
			hint = m.styles.Error.Render("✗ Invalid: use lowercase, numbers, hyphens, underscores")
		}
	}

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.Label.Render("Project Name:"),
		input,
		"",
		hint,
		"",
		m.styles.Help.Render("Examples: my-project, awesome_app, api2go"),
	)

	footer := m.renderFooter()
	helpKeys := m.renderKeyboardHelp("Enter", "Next", "TAB", "Cycle")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		m.renderContainer(form),
		"",
		helpKeys,
		"",
		footer,
	)

	return m.padContent(content)
}

func (m *Model) viewModuleName() string {
	header := m.renderHeader("Go Module", 2, 6)

	input := m.renderInputField(1)

	hint := ""
	value := m.inputs[1].Value()
	defaultModule := fmt.Sprintf("github.com/example/%s", m.projectName)

	if value != "" {
		if isValidModuleName(value) {
			hint = m.styles.Success.Render("✓ Valid module name")
		} else {
			hint = m.styles.Error.Render("✗ Invalid: use 'domain.com/org/project' format")
		}
	} else {
		hint = m.styles.Info.Render("→ Will default to: " + defaultModule)
	}

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.Label.Render("Go Module Name (optional):"),
		input,
		"",
		hint,
		"",
		m.styles.Help.Render("Examples: github.com/acme/myapp, gitlab.com/team/project"),
	)

	footer := m.renderFooter()
	helpKeys := m.renderKeyboardHelp("Enter", "Next", "TAB", "Cycle")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		m.renderContainer(form),
		"",
		helpKeys,
		"",
		footer,
	)

	return m.padContent(content)
}

func (m *Model) viewProjectPath() string {
	header := m.renderHeader("Project Location", 3, 6)

	input := m.renderInputField(2)

	hint := ""
	value := m.inputs[2].Value()

	if value != "" {
		if isValidPath(value) {
			hint = m.styles.Success.Render("✓ Valid path")
		} else {
			hint = m.styles.Error.Render("✗ Invalid path")
		}
	} else {
		hint = m.styles.Info.Render("→ Default: current directory (.)")
	}

	fullPath := value + "/" + m.projectName
	if value == "." {
		fullPath = "./" + m.projectName
	}

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.Label.Render("Project Path:"),
		input,
		"",
		hint,
		"",
		m.styles.Description.Render("Where to create the project:"),
		m.styles.Info.Render(fullPath),
	)

	footer := m.renderFooter()
	helpKeys := m.renderKeyboardHelp("Enter", "Next", "TAB", "Cycle")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		m.renderContainer(form),
		"",
		helpKeys,
		"",
		footer,
	)

	return m.padContent(content)
}

func (m *Model) viewFeatures() string {
	header := m.renderHeader("Select Features", 4, 6)

	featuresList := ""
	for i, feat := range m.features {
		var checkbox string
		if feat.Selected {
			checkbox = m.styles.Success.Render("[✓]")
		} else {
			checkbox = m.styles.Blurred.Render("[ ]")
		}

		var featureText string
		if i == m.featureFocus {
			cursor := m.styles.Focused.Render("▸")
			name := m.styles.Focused.Render(feat.Name)
			featureText = fmt.Sprintf("  %s %s %s", cursor, checkbox, name)
		} else {
			featureText = fmt.Sprintf("    %s %s", checkbox, feat.Name)
		}

		featuresList += featureText
		if i < len(m.features)-1 {
			featuresList += "\n"
		}
	}

	// Show description of focused feature
	if m.featureFocus >= 0 && m.featureFocus < len(m.features) {
		desc := m.features[m.featureFocus].Description
		featuresList += "\n\n" + m.styles.Blurred.Render("    "+desc)
	}

	selectedCount := 0
	for _, feat := range m.features {
		if feat.Selected {
			selectedCount++
		}
	}

	summary := fmt.Sprintf(
		"%s / %d features selected",
		m.styles.Info.Render(fmt.Sprintf("%d", selectedCount)),
		len(m.features),
	)

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.Label.Render("Choose features to include:"),
		"",
		featuresList,
		"",
		summary,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		m.renderContainer(form),
		"",
	)

	// Add warning/info message if present
	if m.warning != "" {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			m.styles.Description.Render(m.warning),
			"",
		)
	}

	footer := m.renderFooter()
	helpKeys := m.styles.Help.Render("SPACE = Toggle  •  UP/DOWN = Navigate  •  ENTER = Next")

	content = lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		helpKeys,
		"",
		footer,
	)

	return m.padContent(content)
}

func (m *Model) viewEnvVars() string {
	header := m.renderHeader("Environment Variables", 5, 6)

	envFields := []struct {
		key   string
		label string
		desc  string
	}{
		{"DB_HOST", "Database Host", "PostgreSQL host"},
		{"DB_PORT", "Database Port", "PostgreSQL port"},
		{"DB_USER", "Database User", "PostgreSQL username"},
		{"DB_PASSWORD", "Database Password", "PostgreSQL password"},
		{"DB_NAME", "Database Name", "Database name"},
		{"JWT_SECRET", "JWT Secret", "Secret key for JWT"},
		{"MINIO_ACCESS_KEY", "MinIO Access Key", "MinIO access key"},
		{"MINIO_SECRET_KEY", "MinIO Secret Key", "MinIO secret key"},
	}

	defaults := map[string]string{
		"DB_HOST":          "localhost",
		"DB_PORT":          "5432",
		"DB_USER":          "postgres",
		"DB_PASSWORD":      "postgres",
		"DB_NAME":          m.projectName,
		"JWT_SECRET":       "your-secret-key-change-in-production",
		"MINIO_ACCESS_KEY": "minioadmin",
		"MINIO_SECRET_KEY": "minioadmin",
	}

	var lines []string
	for i, field := range envFields {
		value, exists := m.envVars[field.key]
		if !exists {
			value = defaults[field.key]
		}

		cursor := "  "
		style := m.styles.Blurred
		if i == m.envFocus {
			cursor = m.styles.Focused.Render("▸ ")
			style = m.styles.Focused
		}

		var line string
		if m.envEditing && i == m.envFocus {
			line = fmt.Sprintf("%s%s: %s",
				cursor,
				style.Render(field.label),
				m.envInput.View(),
			)
		} else {
			line = fmt.Sprintf("%s%s: %s",
				cursor,
				style.Render(field.label),
				m.styles.Description.Render(value),
			)
		}
		lines = append(lines, line)
	}

	instruction := m.styles.Info.Render("Press ENTER to edit selected value")
	if m.envEditing {
		instruction = m.styles.Info.Render("Press ENTER to save, ESC to cancel")
	}

	skipText := m.styles.Help.Render("TAB = Skip to next step")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		m.styles.Label.Render("Configure environment variables:"),
		"",
		lipgloss.JoinVertical(lipgloss.Left, lines...),
		"",
		instruction,
		skipText,
		"",
		m.styles.Help.Render("UP/DOWN = Navigate"),
	)

	return m.padContent(content)
}

func (m *Model) viewConfirm() string {
	header := m.renderHeader("Review & Confirm", 6, 6)

	fullPath := m.projectPath + "/" + m.projectName
	if m.projectPath == "." {
		fullPath = "./" + m.projectName
	}

	selectedFeatures := ""
	for _, feat := range m.features {
		if feat.Selected {
			selectedFeatures += "✓ " + feat.Name + "\n"
		}
	}
	if selectedFeatures != "" {
		selectedFeatures = strings.TrimSuffix(selectedFeatures, "\n")
	}

	details := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderKeyValue("Project Name", m.projectName),
		m.renderKeyValue("Go Module", m.moduleName),
		m.renderKeyValue("Project Path", fullPath),
		"",
		m.styles.Label.Render("Selected Features:"),
		selectedFeatures,
	)

	confirmBox := m.styles.ContainerPrimary.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.Focused.Render("✓ Everything looks good!"),
			"",
			details,
		),
	)

	buttonWidth := 20

	createBtn := lipgloss.NewStyle().
		Foreground(lipgloss.Color("0")).
		Background(m.styles.Focused.GetForeground()).
		Padding(0, 1).
		Bold(true).
		Width(buttonWidth).
		Align(lipgloss.Center).
		Render("Create Project")

	cancelBtn := lipgloss.NewStyle().
		Foreground(m.styles.Blurred.GetForeground()).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.Blurred.GetForeground()).
		Padding(0, 1).
		Width(buttonWidth).
		Align(lipgloss.Center).
		Render("CTRL+C Cancel")

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		createBtn,
		"   ",
		cancelBtn,
	)

	buttons = lipgloss.NewStyle().
		Width(CONTAINER_WIDTH).
		Align(lipgloss.Center).
		Render(buttons)

	footer := m.renderFooter()
	helpKeys := m.styles.Help.Render("Press ENTER to create project or CTRL+C to cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		confirmBox,
		"",
		buttons,
		"",
		helpKeys,
		"",
		footer,
	)

	return m.padContent(content)
}

func (m *Model) viewProcessing() string {
	header := m.renderHeader("Creating Project", 5, 5)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"",
		m.styles.Info.Render(m.spinner.View()+" Processing..."),
		"",
		m.styles.Description.Render("Setting up project structure..."),
		m.styles.Description.Render("Creating directories and files..."),
		m.styles.Description.Render("Initializing git repository..."),
		"",
	)

	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		m.renderContainer(content),
		"",
		m.styles.Help.Render("(This may take a moment)"),
	)

	return m.padContent(fullContent)
}

func (m *Model) viewSuccess() string {
	header := m.renderHeader("Success!", 5, 5)

	fullPath := m.projectPath + "/" + m.projectName
	if m.projectPath == "." {
		fullPath = "./" + m.projectName
	}

	successContent := m.styles.ContainerPrimary.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.Success.Render("✓ Project created successfully!"),
			"",
			m.renderKeyValue("Location", fullPath),
			m.renderKeyValue("Module", m.moduleName),
		),
	)

	nextSteps := m.renderContainer(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.Focused.Render("📋 Next Steps:"),
			"",
			"1. "+m.styles.Description.Render(fmt.Sprintf("cd %s", fullPath)),
			"2. "+m.styles.Description.Render("cp .env.example .env"),
			"3. "+m.styles.Description.Render("make dev-d"),
			"4. "+m.styles.Description.Render("Visit http://localhost:8080/swagger"),
		),
	)

	// Build selected features list
	selectedFeaturesList := []string{}
	featureDescriptions := map[string]string{
		"Authentication (JWT)": "✓ JWT Authentication & Token Rotation",
		"User Management":      "✓ User Management with RBAC",
		"Database":             "✓ PostgreSQL Database Integration",
		"File Storage":         "✓ MinIO File Storage",
		"API Docs":             "✓ Auto-Generated Swagger Docs",
		"Docker":               "✓ Docker & Docker Compose Setup",
		"Podman":               "✓ Podman & Podman Compose Setup",
	}

	for _, feat := range m.features {
		if feat.Selected {
			if desc, ok := featureDescriptions[feat.Name]; ok {
				selectedFeaturesList = append(selectedFeaturesList, desc)
			}
		}
	}

	// Always include base features
	selectedFeaturesList = append(selectedFeaturesList, "✓ Structured Logging (Zap)")
	selectedFeaturesList = append(selectedFeaturesList, "✓ Error Handling & Response Formatting")

	// Build features content with header
	featureContent := []string{
		m.styles.Focused.Render("🚀 Included Features:"),
		"",
	}
	featureContent = append(featureContent, selectedFeaturesList...)

	features := m.renderContainer(
		lipgloss.JoinVertical(
			lipgloss.Left,
			featureContent...,
		),
	)

	footer := m.renderFooter()
	helpKeys := m.renderKeyboardHelp("Enter", "Exit", "Q", "Quit")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		successContent,
		"",
		nextSteps,
		"",
		features,
		"",
		helpKeys,
		"",
		footer,
	)

	return m.padContent(content)
}

func (m *Model) viewError() string {
	header := m.renderHeader("Error", 5, 5)

	errorBox := m.styles.ContainerPrimary.Render(
		m.styles.Error.Render("✗ " + m.err.Error()),
	)

	footer := m.renderFooter()
	helpKeys := m.renderKeyboardHelp("Enter", "Try Again", "CTRL+C", "Exit")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"",
		errorBox,
		"",
		helpKeys,
		"",
		footer,
	)

	return m.padContent(content)
}

// Helper methods
func (m *Model) renderInputField(idx int) string {
	if idx >= len(m.inputs) {
		return ""
	}

	input := m.inputs[idx].View()

	if m.inputs[idx].Focused() {
		return m.styles.InputFocused.Render(input)
	}
	return m.styles.InputBase.Render(input)
}

func (m *Model) renderHeader(title string, step, totalSteps int) string {
	steps := m.renderStepIndicator(step, totalSteps)
	titleRendered := m.styles.Title.Render(title)
	divider := m.styles.Divider.Render(strings.Repeat("─", CONTAINER_WIDTH))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		steps,
		titleRendered,
		divider,
	)
}

func (m *Model) renderFooter() string {
	return m.styles.Divider.Render(strings.Repeat("─", CONTAINER_WIDTH))
}

func (m *Model) renderContainer(content string) string {
	return m.styles.Container.Render(content)
}

func (m *Model) renderKeyboardHelp(key1, action1, key2, action2 string) string {
	help := fmt.Sprintf("%s = %s", m.styles.Info.Render(key1), action1)
	if key2 != "" && action2 != "" {
		help += "  •  " + fmt.Sprintf("%s = %s", m.styles.Info.Render(key2), action2)
	}
	return m.styles.Help.Render(help)
}

func (m *Model) renderStepIndicator(current, total int) string {
	steps := "Step " + fmt.Sprintf("%d/%d", current, total) + "  "
	for i := 1; i <= total; i++ {
		if i == current {
			steps += m.styles.ProgressActive.Render("●")
		} else if i < current {
			steps += m.styles.ProgressDone.Render("●")
		} else {
			steps += m.styles.ProgressTodo.Render("○")
		}
		if i < total {
			steps += " "
		}
	}
	return steps
}

func (m *Model) renderKeyValue(key, value string) string {
	return m.styles.Label.Render(key+":") + "  " + m.styles.Info.Render(value)
}

func (m *Model) padContent(content string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Top, content)
}

// Validators
func isValidProjectName(name string) bool {
	if len(name) == 0 || len(name) > 50 {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-z0-9_-]+$`, name)
	return matched
}

func isValidModuleName(name string) bool {
	if len(name) == 0 || !strings.Contains(name, "/") {
		return false
	}
	matched, _ := regexp.MatchString(`^[a-z0-9./_-]+$`, name)
	return matched
}

func isValidPath(path string) bool {
	if len(path) == 0 {
		return true
	}
	// Allow relative paths like ., .., ./path, ../path, path/to/dir
	matched, _ := regexp.MatchString(`^\.{1,2}(/[a-zA-Z0-9._-]+)*$|^[a-zA-Z0-9._-]+(/[a-zA-Z0-9._-]+)*$|^\./?$`, path)
	return matched
}

// enableDependencies ensures required dependencies are selected
func (m *Model) enableDependencies(featureName string) {
	deps, ok := m.featureDependencies[featureName]
	if !ok {
		return
	}

	for _, depName := range deps {
		for i := range m.features {
			if m.features[i].Name == depName {
				m.features[i].Selected = true
				// Recursively enable dependencies
				m.enableDependencies(depName)
				break
			}
		}
	}
}

// checkDependents ensures features that depend on deselected feature are disabled
func (m *Model) checkDependents() {
	for i := range m.features {
		if !m.features[i].Selected {
			continue
		}

		// Check if all dependencies are met
		deps, ok := m.featureDependencies[m.features[i].Name]
		if !ok {
			continue
		}

		for _, depName := range deps {
			found := false
			for _, feat := range m.features {
				if feat.Name == depName && feat.Selected {
					found = true
					break
				}
			}

			if !found {
				// Dependency not met, deselect this feature
				m.features[i].Selected = false
				// Recursively check dependents of this feature
				m.checkDependents()
				return
			}
		}
	}
}

// getHelpText returns keyboard shortcuts and instructions
func (m *Model) getHelpText() string {
	return `KEYBOARD SHORTCUTS & INSTRUCTIONS

📋 NAVIGATION
  ↑ / ↓          Navigate menu items or features
  TAB / SHIFT+TAB Switch between input fields
  ENTER          Proceed / Confirm selection
  CTRL+C         Cancel and exit anytime

✓ FEATURE SELECTION (when available)
  SPACE          Toggle feature selection
  Dependencies are auto-managed (enabled/disabled as needed)

📝 TEXT INPUT
  Type normally   Enter text
  Backspace/Del  Delete characters
  Arrows         Move cursor (if supported)

ℹ️  WORKFLOW
  1. Select 'Create New Project' from main menu
  2. Enter project name (lowercase, hyphens/underscores)
  3. Enter Go module path (or press ENTER for default)
  4. Select features you need
  5. Confirm to create project

💡 TIPS
  • Project names: my-project, my_api, api2go
  • Module format: github.com/org/project
  • Features auto-select dependencies
  • Project created in parent directory (../project-name)

Press CTRL+C to return to main menu`
}

// getDependencyWarning returns warning message about feature dependencies
func (m *Model) getDependencyWarning() string {
	var warnings []string

	for _, feat := range m.features {
		if !feat.Selected {
			continue
		}

		deps, ok := m.featureDependencies[feat.Name]
		if !ok || len(deps) == 0 {
			continue
		}

		// Check which dependencies are not enabled
		missing := []string{}
		for _, depName := range deps {
			found := false
			for _, f := range m.features {
				if f.Name == depName && f.Selected {
					found = true
					break
				}
			}
			if !found {
				missing = append(missing, depName)
			}
		}

		if len(missing) > 0 {
			msg := fmt.Sprintf("ℹ %s requires: %s", feat.Name, strings.Join(missing, ", "))
			warnings = append(warnings, msg)
		}
	}

	if len(warnings) > 0 {
		return strings.Join(warnings, " | ")
	}

	// Show what was auto-enabled
	autoEnabled := []string{}
	for _, feat := range m.features {
		if feat.Selected && !feat.Default {
			autoEnabled = append(autoEnabled, feat.Name)
		}
	}

	if len(autoEnabled) > 0 {
		return fmt.Sprintf("ℹ Auto-enabled: %s", strings.Join(autoEnabled, ", "))
	}

	return ""
}
