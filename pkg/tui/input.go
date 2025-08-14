package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// InputMode represents different input modes
type InputMode int

const (
	InputModeText   InputMode = iota // Full text input with cursor
	InputModeBuffer                  // Simple character accumulation
)

const inputPrompt = "> "

// InputModel handles text input and buffer operations
type InputModel struct {
	width  int
	height int

	textInput   textinput.Model
	inputBuffer string
	mode        InputMode
	isVisible   bool
}

// NewInputModel creates a new input model
func NewInputModel() *InputModel {
	// Initialize text input
	ti := textinput.New()
	ti.Prompt = inputPrompt

	return &InputModel{
		textInput:   ti,
		mode:        InputModeText,
		isVisible:   false, // Default to hidden
		inputBuffer: "",
	}
}

// SetSize updates the model dimensions
func (m *InputModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// Show makes the input visible and focuses it
func (m *InputModel) Show(initialValue string) {
	m.isVisible = true
	if m.mode == InputModeText {
		m.textInput.SetValue(initialValue)
		m.textInput.SetCursor(len(initialValue))
		m.textInput.Focus()
	}
}

// Hide makes the input invisible and clears it
func (m *InputModel) Hide() {
	m.isVisible = false
	if m.mode == InputModeText {
		m.textInput.Blur()
		m.textInput.SetValue("")
	}
	m.inputBuffer = ""
}

// IsVisible returns whether the input is currently visible
func (m *InputModel) IsVisible() bool {
	return m.isVisible
}

// SetMode switches between text input and buffer modes
func (m *InputModel) SetMode(mode InputMode) {
	m.mode = mode
}

// GetValue returns the current input value based on mode
func (m *InputModel) GetValue() string {
	if m.mode == InputModeText {
		return m.textInput.Value()
	}

	return m.inputBuffer
}

// SetBuffer sets the input buffer value (buffer mode only)
func (m *InputModel) SetBuffer(value string) {
	m.inputBuffer = value
}

// GetBuffer returns the current input buffer value
func (m *InputModel) GetBuffer() string {
	return m.inputBuffer
}

// AppendToBuffer appends a character to the input buffer (buffer mode only)
func (m *InputModel) AppendToBuffer(keyStr string) {
	if keyStr == "backspace" {
		if len(m.inputBuffer) > 0 {
			m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
		}
	} else if len(keyStr) == 1 {
		// For single character keys, append to buffer
		m.inputBuffer += keyStr
	}
}

// ClearBuffer clears the input buffer
func (m *InputModel) ClearBuffer() {
	m.inputBuffer = ""
}

// GetTextInput returns the text input model for direct manipulation
func (m *InputModel) GetTextInput() *textinput.Model {
	return &m.textInput
}

// UpdateTextInput updates the text input model
func (m *InputModel) UpdateTextInput(textInput textinput.Model) {
	m.textInput = textInput
}

// InputCompletedMessage indicates that input has been completed
type InputCompletedMessage struct {
	Value string // The final input value
}

// Update handles Bubbletea messages
func (m *InputModel) Update(msg tea.Msg) (*InputModel, tea.Cmd) {
	var cmd tea.Cmd

	// Only process messages when visible and in text mode
	if m.isVisible && m.mode == InputModeText {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				// Input completed
				value := m.textInput.Value()

				return m, func() tea.Msg {
					return InputCompletedMessage{Value: value}
				}
			case "esc":
				// Cancel input
				return m, func() tea.Msg {
					return InputCompletedMessage{Value: ""}
				}
			}
		}

		// Update text input
		m.textInput, cmd = m.textInput.Update(msg)
	}

	return m, cmd
}

// View renders the input view
func (m *InputModel) View() string {
	if !m.isVisible {
		return ""
	}

	switch m.mode {
	case InputModeText:
		return m.textInput.View()
	case InputModeBuffer:
		// In buffer mode, just display the current buffer with a simple prompt
		return inputPrompt + m.inputBuffer
	default:
		return ""
	}
}