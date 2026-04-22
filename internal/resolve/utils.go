package resolve

import (
	"bufio"
	"codespacegen/internal/i18n"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	portOnlyPattern    = regexp.MustCompile(`^\d+$`)
	portMappingPattern = regexp.MustCompile(`^\d+:\d+$`)
)

func normalizePortMapping(value string) (string, error) {
	v := strings.TrimSpace(value)
	if portOnlyPattern.MatchString(v) {
		return fmt.Sprintf("%s:%s", v, v), nil
	}
	if portMappingPattern.MatchString(v) {
		return v, nil
	}

	return "", errors.New(i18n.T("error_invalid_port_mapping", map[string]interface{}{"Value": value}))
}

func promptWithDefault(reader *bufio.Reader, prompt string, defaultValue string) (string, error) {
	fmt.Print(prompt)
	line, err := reader.ReadString('\n')
	if err != nil {
		if errors.Is(err, io.EOF) {
			line = strings.TrimSpace(line)
			if line == "" {
				return defaultValue, nil
			}
			return line, nil
		}
		return "", err
	}

	line = strings.TrimSpace(line)
	if line == "" {
		return defaultValue, nil
	}

	return line, nil
}

func promptUntilResolved(reader *bufio.Reader, defaultValue string, promptFn func(), readErrMessage string, handleLine func(line string, defaultValue string, isEOF bool) (string, bool, error)) (string, error) {
	for {
		promptFn()

		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				line = strings.TrimSpace(line)
				result, _, handleErr := handleLine(line, defaultValue, true)
				return result, handleErr
			}
			return "", fmt.Errorf("%s: %w", readErrMessage, err)
		}

		line = strings.TrimSpace(line)
		result, done, handleErr := handleLine(line, defaultValue, false)
		if handleErr != nil {
			return "", handleErr
		}
		if done {
			return result, nil
		}
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		trimmed := strings.TrimSpace(v)
		if trimmed != "" {
			return trimmed
		}
	}

	return ""
}
