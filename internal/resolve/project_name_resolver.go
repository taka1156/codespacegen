package resolve

import (
	"bufio"
	"codespacegen/internal/i18n"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func (cscr *CodeSpaceConfigResolver) ResolveProjectName(explicitProjectName string) (string, error) {
	defaultProjectName := strings.TrimSpace(explicitProjectName)
	reader := bufio.NewReader(os.Stdin)

	for {
		printProjectNamePrompt(defaultProjectName)

		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return resolveProjectNameOnEOF(line, defaultProjectName)
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

func printProjectNamePrompt(defaultProjectName string) {
	if defaultProjectName == "" {
		fmt.Print(i18n.T("prompt_project_name_required"))
		return
	}

	fmt.Print(i18n.T("prompt_project_name_with_default", map[string]interface{}{"Default": defaultProjectName}))
}

func resolveProjectNameOnEOF(line string, defaultProjectName string) (string, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		if defaultProjectName != "" {
			return defaultProjectName, nil
		}
		return "", fmt.Errorf("%s", i18n.T("error_project_name_required"))
	}

	return line, nil
}
