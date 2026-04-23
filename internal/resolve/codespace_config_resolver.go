package resolve

import (
	"bufio"
	"io"

	"codespacegen/internal/domain/entity"
	"codespacegen/internal/i18n"
	"fmt"
	"strings"
)

type CodespaceConfigResolver struct {
	reader *bufio.Reader
}

func NewCodespaceConfigResolver(r io.Reader) *CodespaceConfigResolver {
	return &CodespaceConfigResolver{
		reader: bufio.NewReader(r),
	}
}

func (cscr *CodespaceConfigResolver) ResolveLanguage(explicitLanguage string) (string, error) {
	defaultLanguage := strings.TrimSpace(explicitLanguage)
	value, err := promptWithDefault(cscr.reader, i18n.T("prompt_language"), defaultLanguage)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_language"), err)
	}
	return strings.ToLower(strings.TrimSpace(value)), nil
}

func (cscr *CodespaceConfigResolver) ResolveWorkspaceFolder(explicitWorkspaceFolder string) (string, error) {
	defaultWorkspaceFolder := strings.TrimSpace(explicitWorkspaceFolder)
	if defaultWorkspaceFolder == "" {
		defaultWorkspaceFolder = "/workspace"
	}
	value, err := promptWithDefault(cscr.reader, i18n.T("prompt_workspace_folder", map[string]interface{}{"Default": defaultWorkspaceFolder}), defaultWorkspaceFolder)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_workspace_folder"), err)
	}
	return strings.TrimSpace(value), nil
}

func (cscr *CodespaceConfigResolver) ResolveTimezone(explicitTimezone string, configTimezone string, defaultTimezone string) (string, error) {
	resolved := strings.TrimSpace(explicitTimezone)
	if resolved == "" {
		resolved = strings.TrimSpace(configTimezone)
	}
	if resolved == "" {
		resolved = strings.TrimSpace(defaultTimezone)
	}
	if resolved == "" {
		resolved = entity.DefaultTimezone
	}

	value, err := promptWithDefault(cscr.reader, i18n.T("prompt_timezone", map[string]interface{}{"Default": resolved}), resolved)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_timezone"), err)
	}

	return strings.TrimSpace(value), nil
}

func (cscr *CodespaceConfigResolver) ResolveServiceName(explicitServiceName string) (string, error) {
	defaultServiceName := strings.TrimSpace(explicitServiceName)
	if defaultServiceName == "" {
		defaultServiceName = "app"
	}
	value, err := promptWithDefault(cscr.reader, i18n.T("prompt_service_name", map[string]interface{}{"Default": defaultServiceName}), defaultServiceName)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_service_name"), err)
	}
	return strings.TrimSpace(value), nil
}
