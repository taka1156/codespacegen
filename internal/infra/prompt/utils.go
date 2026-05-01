package prompt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

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
