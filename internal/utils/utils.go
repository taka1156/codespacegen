package utils

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func Ptr[T any](v T) *T {
	return &v
}

var (
	portOnlyPattern    = regexp.MustCompile(`^\d+$`)
	portMappingPattern = regexp.MustCompile(`^\d+:\d+$`)
)

func NormalizePortMapping(value string) (string, error) {
	v := strings.TrimSpace(value)
	if portOnlyPattern.MatchString(v) {
		return fmt.Sprintf("%s:%s", v, v), nil
	}
	if portMappingPattern.MatchString(v) {
		return v, nil
	}

	return "", errors.New("invalid port format")
}

// ResolveOutputPath joins outputDir and relativePath after sanitizing and validating
// that the result does not escape outputDir.
func ResolveOutputPath(outputDir string, relativePath string) (string, error) {
	cleanRelativePath, err := sanitizeRelativePath(relativePath)
	if err != nil {
		return "", err
	}

	joinedPath := filepath.Join(outputDir, cleanRelativePath)
	if err := validatePathWithinOutputDir(outputDir, joinedPath, relativePath); err != nil {
		return "", err
	}

	return joinedPath, nil
}

func sanitizeRelativePath(relativePath string) (string, error) {
	cleanRelativePath := filepath.Clean(relativePath)
	if cleanRelativePath == "." || cleanRelativePath == "" {
		return "", fmt.Errorf("invalid file path: %s", relativePath)
	}
	if filepath.IsAbs(cleanRelativePath) {
		return "", fmt.Errorf("absolute path is not allowed: %s", relativePath)
	}

	return cleanRelativePath, nil
}

func validatePathWithinOutputDir(outputDir string, joinedPath string, relativePath string) error {
	cleanOutputDir := filepath.Clean(outputDir)
	relativeToOutputDir, err := filepath.Rel(cleanOutputDir, joinedPath)
	if err != nil {
		return fmt.Errorf("failed to calculate relative path: %w", err)
	}

	if relativeToOutputDir == ".." || strings.HasPrefix(relativeToOutputDir, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path escapes output directory: %s", relativePath)
	}

	return nil
}
