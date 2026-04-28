package resolve

import (
	"codespacegen/internal/i18n"
	"fmt"
	"strings"
)

func (cscr *CodespaceConfigResolver) ResolvePortMapping(explicitPort string) (string, error) {
	defaultPort := strings.TrimSpace(explicitPort)
	return promptUntilResolved(
		cscr.reader,
		defaultPort,
		func() { printPortPrompt(defaultPort) },
		i18n.T("error_failed_to_read_port_input"),
		func(line string, defaultValue string, _ bool) (string, bool, error) {
			if line == "" {
				if defaultValue != "" {
					normalizedDefault, err := normalizePortMapping(defaultValue)
					if err != nil {
						fmt.Println(i18n.T("error_invalid_port_format"))
						return "", false, nil
					}
					return normalizedDefault, true, nil
				}
				return "", true, nil
			}

			normalized, err := normalizePortMapping(line)
			if err == nil {
				return normalized, true, nil
			}

			fmt.Println(i18n.T("error_invalid_port_format"))
			return "", false, nil
		},
	)
}

func printPortPrompt(defaultPort string) {
	if defaultPort == "" {
		fmt.Print(i18n.T("prompt_port_empty"))
		return
	}

	fmt.Print(i18n.T("prompt_port_with_default", map[string]interface{}{"Default": defaultPort}))
}
