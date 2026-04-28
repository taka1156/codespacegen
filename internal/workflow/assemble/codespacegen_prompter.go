package assemble

type CodespacegenPrompter interface {
	PromptProjectName(explicitProjectName string) (string, error)
	PromptLanguage(explicitLanguage string) (string, error)
	PromptWorkspaceFolder(explicitWorkspaceFolder string) (string, error)
	PromptServiceName(explicitServiceName string) (string, error)
	PromptPortMapping(explicitPort string) (string, error)
	PromptTimezone(defaultTimezone string) (string, error)
}
