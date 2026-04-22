package resolve

import (
	"bufio"
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/i18n"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type CodeSpaceConfigResolver struct {
	
}

func NewCodeSpaceConfigResolver() *CodeSpaceConfigResolver {
	return &CodeSpaceConfigResolver{}
}

func (cscr *CodeSpaceConfigResolver) ResolveProjectName(explicitProjectName string) (string, error) {
	defaultProjectName := strings.TrimSpace(explicitProjectName)
	reader := bufio.NewReader(os.Stdin)

	for {
		if defaultProjectName == "" {
			fmt.Print(i18n.T("prompt_project_name_required"))
		} else {
			fmt.Print(i18n.T("prompt_project_name_with_default", map[string]interface{}{"Default": defaultProjectName}))
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				line = strings.TrimSpace(line)
				if line == "" {
					if defaultProjectName != "" {
						return defaultProjectName, nil
					}
					return "", fmt.Errorf("%s", i18n.T("error_project_name_required"))
				}
				return line, nil
			}
			return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_project_name"), err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			if defaultProjectName != "" {
				return defaultProjectName, nil
			}
			fmt.Println(i18n.T("msg_project_name_mandatory"))
			continue
		}

		return line, nil
	}
}

func (cscr *CodeSpaceConfigResolver) ResolveLanguage(explicitLanguage string) (string, error) {
	defaultLanguage := strings.TrimSpace(explicitLanguage)
	value, err := promptWithDefault(i18n.T("prompt_language"), defaultLanguage)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_language"), err)
	}
	return strings.ToLower(strings.TrimSpace(value)), nil
}

func (cscr *CodeSpaceConfigResolver) ResolveWorkspaceFolder(explicitWorkspaceFolder string) (string, error) {
	defaultWorkspaceFolder := strings.TrimSpace(explicitWorkspaceFolder)
	if defaultWorkspaceFolder == "" {
		defaultWorkspaceFolder = "/workspace"
	}
	value, err := promptWithDefault(i18n.T("prompt_workspace_folder", map[string]interface{}{"Default": defaultWorkspaceFolder}), defaultWorkspaceFolder)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_workspace_folder"), err)
	}
	return strings.TrimSpace(value), nil
}

func (cscr *CodeSpaceConfigResolver) ResolveTimezone(explicitTimezone string, configTimezone string) (string, error) {
	defaultTimezone := strings.TrimSpace(explicitTimezone)
	if defaultTimezone == "" {
		defaultTimezone = strings.TrimSpace(configTimezone)
	}
	if defaultTimezone == "" {
		defaultTimezone = entity.DefaultTimezone
	}

	value, err := promptWithDefault(i18n.T("prompt_timezone", map[string]interface{}{"Default": defaultTimezone}), defaultTimezone)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_timezone"), err)
	}

	return strings.TrimSpace(value), nil
}

func (cscr *CodeSpaceConfigResolver) ResolveServiceName(explicitServiceName string) (string, error) {
	defaultServiceName := strings.TrimSpace(explicitServiceName)
	if defaultServiceName == "" {
		defaultServiceName = "app"
	}
	value, err := promptWithDefault(i18n.T("prompt_service_name", map[string]interface{}{"Default": defaultServiceName}), defaultServiceName)
	if err != nil {
		return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_service_name"), err)
	}
	return strings.TrimSpace(value), nil
}

func (cscr *CodeSpaceConfigResolver) ResolvePortMapping(explicitPort string) (string, error) {
	defaultPort := strings.TrimSpace(explicitPort)
	reader := bufio.NewReader(os.Stdin)
	for {
		if defaultPort == "" {
			fmt.Print(i18n.T("prompt_port_empty"))
		} else {
			fmt.Print(i18n.T("prompt_port_with_default", map[string]interface{}{"Default": defaultPort}))
		}
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				line = strings.TrimSpace(line)
				if line == "" {
					if defaultPort != "" {
						return normalizePortMapping(defaultPort)
					}
					return "", nil
				}
				return normalizePortMapping(line)
			}
			return "", fmt.Errorf("%s: %w", i18n.T("error_failed_to_read_port_input"), err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			if defaultPort != "" {
				return normalizePortMapping(defaultPort)
			}
			return "", nil
		}

		normalized, err := normalizePortMapping(line)
		if err == nil {
			return normalized, nil
		}

		fmt.Println(i18n.T("error_invalid_port_format"))
	}
}

func (cscr *CodeSpaceConfigResolver) ResolveBaseImage(language string, explicitBaseImage string, imageConfig string, jsonEntries map[string]entity.JsonEntry) (entity.JsonEntry, error) {
	if explicitBaseImage != "" {
		return entity.JsonEntry{Image: explicitBaseImage}, nil
	}

	if strings.TrimSpace(language) == "" {
		return entity.JsonEntry{Image: entity.DefaultImage}, nil
	}

	key := strings.ToLower(strings.TrimSpace(language))
	entry, ok := jsonEntries[key]
	if !ok {
		return entity.JsonEntry{}, errors.New(i18n.T("error_unsupported_language", map[string]interface{}{"Language": language}))
	}

	if entry.Image == "" {
		return entity.JsonEntry{}, errors.New(i18n.T("error_image_required_for_language", map[string]interface{}{"Language": language}))
	}

	return entry, nil
}
