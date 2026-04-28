package prompt

import (
	"bufio"
	"errors"
	"io"

	"codespacegen/internal/i18n"
	"fmt"
	"strings"
)

type CodespacegenPrompter struct {
	reader *bufio.Reader
}

func NewCodespacegenPrompter(reader io.Reader) *CodespacegenPrompter {
	return &CodespacegenPrompter{
		reader: bufio.NewReader(reader),
	}
}

func (cp *CodespacegenPrompter) PromptProjectName(explicitProjectName string) (string, error) {
	defaultProjectName := strings.TrimSpace(explicitProjectName)
	for {
		if defaultProjectName == "" {
			fmt.Print(i18n.T("prompt_project_name_required"))
		} else {
			fmt.Print(i18n.T("prompt_project_name_with_default", map[string]interface{}{"Default": defaultProjectName}))
		}

		line, err := cp.reader.ReadString('\n')
		isEOF := errors.Is(err, io.EOF)
		if err != nil && !isEOF {
			return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_project_name"), err)
		}
		line = strings.TrimSpace(line)
		if line != "" {
			return line, nil
		}
		if defaultProjectName != "" {
			return defaultProjectName, nil
		}
		if isEOF {
			return "", fmt.Errorf("%s", i18n.T("error_project_name_required"))
		}
		fmt.Println(i18n.T("msg_project_name_mandatory"))
	}
}

func (cp *CodespacegenPrompter) PromptLanguage(explicitLanguage string) (string, error) {
	defaultLanguage := strings.TrimSpace(explicitLanguage)
	value, err := promptWithDefault(cp.reader, i18n.T("prompt_language"), defaultLanguage)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_language"), err)
	}
	return strings.ToLower(strings.TrimSpace(value)), nil
}

func (cp *CodespacegenPrompter) PromptWorkspaceFolder(explicitWorkspaceFolder string) (string, error) {
	defaultWorkspaceFolder := strings.TrimSpace(explicitWorkspaceFolder)
	if defaultWorkspaceFolder == "" {
		defaultWorkspaceFolder = "/workspace"
	}
	value, err := promptWithDefault(cp.reader, i18n.T("prompt_workspace_folder", map[string]interface{}{"Default": defaultWorkspaceFolder}), defaultWorkspaceFolder)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_workspace_folder"), err)
	}
	return strings.TrimSpace(value), nil
}

func (cp *CodespacegenPrompter) PromptTimezone(defaultTimezone string) (string, error) {
	value, err := promptWithDefault(cp.reader, i18n.T("prompt_timezone", map[string]interface{}{"Default": defaultTimezone}), defaultTimezone)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_timezone"), err)
	}

	return strings.TrimSpace(value), nil
}

func (cp *CodespacegenPrompter) PromptServiceName(explicitServiceName string) (string, error) {
	defaultServiceName := strings.TrimSpace(explicitServiceName)
	if defaultServiceName == "" {
		defaultServiceName = "app"
	}
	value, err := promptWithDefault(cp.reader, i18n.T("prompt_service_name", map[string]interface{}{"Default": defaultServiceName}), defaultServiceName)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_service_name"), err)
	}
	return strings.TrimSpace(value), nil
}

func (cp *CodespacegenPrompter) PromptPortMapping(explicitPort string) (string, error) {
	defaultPort := strings.TrimSpace(explicitPort)
	for {
		if defaultPort == "" {
			fmt.Print(i18n.T("prompt_port_empty"))
		} else {
			fmt.Print(i18n.T("prompt_port_with_default", map[string]interface{}{"Default": defaultPort}))
		}

		line, err := cp.reader.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_port_input"), err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			if defaultPort == "" {
				return "", nil
			}
			normalized, normErr := normalizePortMapping(defaultPort)
			if normErr == nil {
				return normalized, nil
			}
			fmt.Println(i18n.T("error_invalid_port_format"))
			continue
		}
		normalized, normErr := normalizePortMapping(line)
		if normErr == nil {
			return normalized, nil
		}
		fmt.Println(i18n.T("error_invalid_port_format"))
	}
}
