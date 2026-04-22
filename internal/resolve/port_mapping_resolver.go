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

func (cscr *CodeSpaceConfigResolver) ResolvePortMapping(explicitPort string) (string, error) {
	defaultPort := strings.TrimSpace(explicitPort)
	reader := bufio.NewReader(os.Stdin)

	for {
		printPortPrompt(defaultPort)

		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return resolvePortOnEOF(line, defaultPort)
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

func printPortPrompt(defaultPort string) {
	if defaultPort == "" {
		fmt.Print(i18n.T("prompt_port_empty"))
		return
	}

	fmt.Print(i18n.T("prompt_port_with_default", map[string]interface{}{"Default": defaultPort}))
}

func resolvePortOnEOF(line string, defaultPort string) (string, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		if defaultPort != "" {
			return normalizePortMapping(defaultPort)
		}
		return "", nil
	}

	return normalizePortMapping(line)
}
