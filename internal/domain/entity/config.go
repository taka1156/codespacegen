package entity

import "regexp"

// This file defines the CodespaceConfig struct and related constants.
const DefaultImage = "alpine:latest"

// CodespaceConfig holds values used to generate devcontainer artifacts.
const DefaultTimezone = "UTC"

const DefaultVersion = "dev"

type PortMappingPatterns struct {
	PortOnly    *regexp.Regexp
	PortMapping *regexp.Regexp
}

var DefaultPortMappingPatterns = PortMappingPatterns{
	PortOnly:    regexp.MustCompile(`^\d+$`),
	PortMapping: regexp.MustCompile(`^\d+:\d+$`),
}




type CliConfig struct {
	OutputDir       *string
	ContainerName   *string
	ServiceName     *string
	Language        *string
	WorkspaceFolder *string
	BaseImage       *string
	Timezone        *string
	ImageConfig     *string
	Port            *string
	ComposeFile     *string
	Overwrite       *bool
	Lang            *string
	ShowVersion     *bool
}

type LocaleConfig struct {
	Lang     string
	Language string
	LcAll    string
}

var DefaultLocale = LocaleConfig{
	Lang:     "ja_JP.UTF-8",
	Language: "ja_JP:ja",
	LcAll:    "ja_JP.UTF-8",
}

type JsonEntry struct {
	Image            string
	Install          string
	Locale           LocaleConfig
	Timezone         string
	VSCodeExtensions []string
}

type CodespaceConfig struct {
	ContainerName    string
	ServiceName      string
	WorkspaceFolder  string
	BaseImage        string
	Locale           LocaleConfig
	Timezone         string
	ComposeFileName  string
	PortMapping      string
	InstallCommand   string
	VSCodeExtensions []string
}

func (c CodespaceConfig) Validate() error {
	if c.ContainerName == "" {
		return ErrInvalidConfig("container name is required")
	}
	if c.ServiceName == "" {
		return ErrInvalidConfig("service name is required")
	}
	if c.WorkspaceFolder == "" {
		return ErrInvalidConfig("workspace folder is required")
	}
	if c.BaseImage == "" {
		return ErrInvalidConfig("base image is required")
	}
	if c.ComposeFileName == "" {
		return ErrInvalidConfig("compose file name is required")
	}
	return nil
}
