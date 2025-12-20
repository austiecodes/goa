package set

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/austiecodes/goa/internal/consts"
	"github.com/austiecodes/goa/internal/utils"
)

func createMainMenu() list.Model {
	items := []list.Item{
		MenuItem{title: "provider", desc: "Configure provider settings (API key, base URL)"},
		MenuItem{title: "chat-model", desc: "Set default model for chat"},
		MenuItem{title: "title-model", desc: "Set model for generating conversation titles"},
		MenuItem{title: "think-model", desc: "Set model for thinking"},
		MenuItem{title: "exit", desc: "Exit settings"},
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 60, 14)
	l.Title = "Goa Settings"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)
	return l
}

func createProviderList() list.Model {
	items := []list.Item{
		MenuItem{title: consts.ProviderOpenAI, desc: "OpenAI API (GPT models)"},
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 60, 10)
	l.Title = "Select Provider"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)
	return l
}

func createProviderConfigInputs(config *utils.Config) []textinput.Model {
	inputs := make([]textinput.Model, 2)

	// API Key input
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "sk-..."
	inputs[0].EchoMode = textinput.EchoPassword
	inputs[0].EchoCharacter = '*'
	inputs[0].CharLimit = 256
	inputs[0].Width = 50
	if config.Providers.OpenAI.APIKey != "" {
		inputs[0].SetValue(config.Providers.OpenAI.APIKey)
	}

	// Base URL input
	inputs[1] = textinput.New()
	inputs[1].Placeholder = consts.DefaultBaseURL
	inputs[1].CharLimit = 256
	inputs[1].Width = 50
	if config.Providers.OpenAI.BaseURL != "" {
		inputs[1].SetValue(config.Providers.OpenAI.BaseURL)
	}

	return inputs
}

func createModelList(models []string, mt ModelType) list.Model {
	items := make([]list.Item, len(models))
	for i, modelID := range models {
		items[i] = MenuItem{title: modelID, desc: ""}
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, 60, 20)

	switch mt {
	case ModelTypeChat:
		l.Title = "Select Chat Model"
	case ModelTypeTitle:
		l.Title = "Select Title Model"
	case ModelTypeThink:
		l.Title = "Select Think Model"
	}

	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)
	return l
}

func saveConfig(config *utils.Config) tea.Cmd {
	return func() tea.Msg {
		err := utils.SaveConfig(config)
		return ConfigSavedMsg{Err: err}
	}
}
