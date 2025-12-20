package set

import (
	"fmt"
	"strings"

	"github.com/austiecodes/goa/internal/consts"
	"github.com/austiecodes/goa/internal/types"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) updateMainMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.List.SelectedItem().(MenuItem)
			switch selected.Title() {
			case "provider":
				m.List = createProviderList()
				m.Screen = ScreenProviderSelect
			case "chat-model":
				m.ModelType = ModelTypeChat
				m.List = createProviderList()
				m.Screen = ScreenModelProviderSelect
			case "title-model":
				m.ModelType = ModelTypeTitle
				m.List = createProviderList()
				m.Screen = ScreenModelProviderSelect
			case "think-model":
				m.ModelType = ModelTypeThink
				m.List = createProviderList()
				m.Screen = ScreenModelProviderSelect
			case "exit":
				m.Quitting = true
				return *m, tea.Quit
			}
			return *m, nil
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return *m, cmd
}

func (m *Model) updateProviderSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.List.SelectedItem().(MenuItem)
			if selected.Title() == consts.ProviderOpenAI {
				m.TextInputs = createProviderConfigInputs(m.Config)
				m.FocusedInput = 0
				m.Screen = ScreenProviderConfig
				return *m, m.TextInputs[0].Focus()
			}
			return *m, nil
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return *m, cmd
}

func (m *Model) updateProviderConfig(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "down":
			m.TextInputs[m.FocusedInput].Blur()
			m.FocusedInput = (m.FocusedInput + 1) % len(m.TextInputs)
			return *m, m.TextInputs[m.FocusedInput].Focus()

		case "shift+tab", "up":
			m.TextInputs[m.FocusedInput].Blur()
			m.FocusedInput = (m.FocusedInput - 1 + len(m.TextInputs)) % len(m.TextInputs)
			return *m, m.TextInputs[m.FocusedInput].Focus()

		case "enter":
			// Save config
			apiKey := m.TextInputs[0].Value()
			baseURL := m.TextInputs[1].Value()

			if apiKey == "" {
				m.Err = fmt.Errorf("API key is required")
				return *m, nil
			}

			m.Config.Providers.OpenAI.APIKey = apiKey
			m.Config.Providers.OpenAI.BaseURL = baseURL

			return *m, saveConfig(m.Config)
		}
	}

	// Update focused text input
	var cmd tea.Cmd
	m.TextInputs[m.FocusedInput], cmd = m.TextInputs[m.FocusedInput].Update(msg)
	return *m, cmd
}

func (m *Model) updateModelProviderSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.List.SelectedItem().(MenuItem)
			providerID := selected.Title()
			return *m, loadModelsForProvider(providerID, m.Config)
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return *m, cmd
}

func (m *Model) updateModelSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			selected := m.List.SelectedItem().(MenuItem)
			modelID := selected.Title()

			newModel := &types.Model{
				Provider: consts.ProviderOpenAI,
				ModelID:  modelID,
			}

			switch m.ModelType {
			case ModelTypeChat:
				m.Config.Model.ChatModel = newModel
			case ModelTypeTitle:
				m.Config.Model.TitleModel = newModel
			case ModelTypeThink:
				m.Config.Model.ThinkModel = newModel
			}

			return *m, saveConfig(m.Config)
		}
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return *m, cmd
}

func (m *Model) renderView() string {
	if m.Quitting {
		return "Goodbye!\n"
	}

	var s strings.Builder

	switch m.Screen {
	case ScreenMainMenu:
		s.WriteString(m.List.View())

	case ScreenProviderSelect:
		s.WriteString(TitleStyle.Render("Select Provider"))
		s.WriteString("\n\n")
		s.WriteString(m.List.View())

	case ScreenProviderConfig:
		s.WriteString(TitleStyle.Render("Configure OpenAI Provider"))
		s.WriteString("\n\n")
		for i, input := range m.TextInputs {
			label := ""
			switch i {
			case 0:
				label = "API Key (required)"
			case 1:
				label = "Base URL (optional, default: OpenAI API)"
			}
			s.WriteString(InputLabelStyle.Render(label))
			s.WriteString("\n")
			s.WriteString(input.View())
			s.WriteString("\n\n")
		}
		s.WriteString(HelpStyle.Render("Press Enter to save, Esc to cancel, Tab/Shift+Tab to navigate"))

	case ScreenModelProviderSelect:
		modelName := ""
		switch m.ModelType {
		case ModelTypeChat:
			modelName = "Chat Model"
		case ModelTypeTitle:
			modelName = "Title Model"
		case ModelTypeThink:
			modelName = "Think Model"
		}
		s.WriteString(TitleStyle.Render(fmt.Sprintf("Select Provider for %s", modelName)))
		s.WriteString("\n\n")
		s.WriteString(m.List.View())

	case ScreenModelSelect:
		modelName := ""
		switch m.ModelType {
		case ModelTypeChat:
			modelName = "Chat Model"
		case ModelTypeTitle:
			modelName = "Title Model"
		case ModelTypeThink:
			modelName = "Think Model"
		}
		s.WriteString(TitleStyle.Render(fmt.Sprintf("Select %s", modelName)))
		s.WriteString("\n\n")
		s.WriteString(m.List.View())
	}

	if m.Err != nil {
		s.WriteString("\n\n")
		s.WriteString(ErrorStyle.Render(fmt.Sprintf("Error: %v", m.Err)))
	}

	return s.String()
}
