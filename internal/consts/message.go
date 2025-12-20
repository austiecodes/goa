package consts

type MessageRole string

const (
	MessageSystem    MessageRole = "system"
	MessageUser      MessageRole = "user"
	MessageAssistant MessageRole = "assistant"
	MessageDeveloper MessageRole = "developer"
)
