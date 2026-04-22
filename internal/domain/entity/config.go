package entity

// This file defines the CodespaceConfig struct and related constants.
const DefaultImage = "alpine:latest"

// DefaultSetting holds the resolved default values used across the application.
type DefaultSetting struct {
	Timezone  string
	Image     string
	Version   string
	VscSchema string
}

// CodespaceConfig holds values used to generate devcontainer artifacts.
const DefaultTimezone = "UTC"

// DefaultVersion is the default version of the generated configuration,
// which can be used for future enhancements or versioning of the generated files.
const DefaultVersion = "dev"

// DefaultVscSchema is the default schema URL for the devcontainer.json file.
const DefaultVscSchema = "https://raw.githubusercontent.com/microsoft/vscode/main/extensions/configuration-editing/schemas/devContainer.vscode.schema.json"

type CliConfig struct {
	OutputDir           *string
	ContainerName       *string
	ServiceName         *string
	Language            *string
	WorkspaceFolder     *string
	BaseImage           *string
	Timezone            *string
	ImageConfig         *string
	Port                *string
	ComposeFile         *string
	EnableOverwriteFile *bool
	Lang                *string
	ShowVersion         *bool
}

func (c CliConfig) OutputDirValue() string {
	return stringValue(c.OutputDir)
}

func (c CliConfig) ContainerNameValue() string {
	return stringValue(c.ContainerName)
}

func (c CliConfig) ServiceNameValue() string {
	return stringValue(c.ServiceName)
}

func (c CliConfig) LanguageValue() string {
	return stringValue(c.Language)
}

func (c CliConfig) WorkspaceFolderValue() string {
	return stringValue(c.WorkspaceFolder)
}

func (c CliConfig) BaseImageValue() string {
	return stringValue(c.BaseImage)
}

func (c CliConfig) TimezoneValue() string {
	return stringValue(c.Timezone)
}

func (c CliConfig) ImageConfigValue() string {
	return stringValue(c.ImageConfig)
}

func (c CliConfig) PortValue() string {
	return stringValue(c.Port)
}

func (c CliConfig) ComposeFileValue() string {
	return stringValue(c.ComposeFile)
}

func (c CliConfig) EnableOverwriteFileValue() bool {
	return boolValue(c.EnableOverwriteFile)
}

func (c CliConfig) LangValue() string {
	return stringValue(c.Lang)
}

func (c CliConfig) ShowVersionValue() bool {
	return boolValue(c.ShowVersion)
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
	Schema           string
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
