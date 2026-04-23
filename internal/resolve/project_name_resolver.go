package resolve

import (
	"codespacegen/internal/i18n"
	"fmt"
	"strings"
)

func (cscr *CodespaceConfigResolver) ResolveProjectName(explicitProjectName string) (string, error) {
	defaultProjectName := strings.TrimSpace(explicitProjectName)
	return promptUntilResolved(
		cscr.reader,
		defaultProjectName,
		func() { printProjectNamePrompt(defaultProjectName) },
		i18n.T("error_failed_to_read_project_name"),
		func(line string, defaultValue string, isEOF bool) (string, bool, error) {
			if line == "" {
				if defaultValue != "" {
					return defaultValue, true, nil
				}
				if isEOF {
					return "", true, fmt.Errorf("%s", i18n.T("error_project_name_required"))
				}
				fmt.Println(i18n.T("msg_project_name_mandatory"))
				return "", false, nil
			}

			return line, true, nil
		},
	)
}

func printProjectNamePrompt(defaultProjectName string) {
	if defaultProjectName == "" {
		fmt.Print(i18n.T("prompt_project_name_required"))
		return
	}

	fmt.Print(i18n.T("prompt_project_name_with_default", map[string]interface{}{"Default": defaultProjectName}))
}
