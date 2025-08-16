package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const inputPrompt = "> "

// InputModel handles text input
type InputModel struct {
	width  int
	height int

	textInput textinput.Model
	isVisible bool
}

// InputCompletedMessage indicates that input has been completed
type InputCompletedMessage struct {
	Value string // The final input value
}

// NewInputModel creates a new input model
func NewInputModel() *InputModel {
	// Initialize text input
	ti := textinput.New()
	ti.Prompt = inputPrompt

	return &InputModel{
		textInput: ti,
		isVisible: false, // Default to hidden
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
	m.textInput.SetValue(initialValue)
	m.textInput.SetCursor(len(initialValue))
	m.textInput.Focus()
}

// Hide makes the input invisible and clears it
func (m *InputModel) Hide() {
	m.isVisible = false
	m.textInput.Blur()
	m.textInput.SetValue("")
}

// IsVisible returns whether the input is currently visible
func (m *InputModel) IsVisible() bool {
	return m.isVisible
}

// GetValue returns the current input value
func (m *InputModel) GetValue() string {
	return m.textInput.Value()
}

// Update handles Bubbletea messages
func (m *InputModel) Update(msg tea.Msg) tea.Cmd {
	if !m.isVisible {
		return nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	return cmd
}

// View renders the input view
func (m *InputModel) View() string {
	if !m.isVisible {
		return ""
	}

	return m.textInput.View()
}