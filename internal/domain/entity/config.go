package entity

const DefaultImage = "alpine:latest"

const DefaultTimezone = "UTC"

const DefaultVscSchema = "https://raw.githubusercontent.com/microsoft/vscode/main/extensions/configuration-editing/schemas/devContainer.vscode.schema.json"

type OsModules struct {
	AlpineModules     []string
	DebianLikeModules []string
}

type DefaultSetting struct {
	Timezone            string
	Image               string
	Version             string
	VscSchema           string
	OsModules           OsModules
	SettingJsonFileName string
}

type ClientConfig struct {
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
	OutputTemplateJson  *bool
	Initialize          *bool
}

func (c ClientConfig) OutputDirValue() string {
	return stringValue(c.OutputDir)
}

func (c ClientConfig) ContainerNameValue() string {
	return stringValue(c.ContainerName)
}

func (c ClientConfig) ServiceNameValue() string {
	return stringValue(c.ServiceName)
}

func (c ClientConfig) LanguageValue() string {
	return stringValue(c.Language)
}

func (c ClientConfig) WorkspaceFolderValue() string {
	return stringValue(c.WorkspaceFolder)
}

func (c ClientConfig) BaseImageValue() string {
	return stringValue(c.BaseImage)
}

func (c ClientConfig) TimezoneValue() string {
	return stringValue(c.Timezone)
}

func (c ClientConfig) ImageConfigValue() string {
	return stringValue(c.ImageConfig)
}

func (c ClientConfig) PortValue() string {
	return stringValue(c.Port)
}

func (c ClientConfig) ComposeFileValue() string {
	return stringValue(c.ComposeFile)
}

func (c ClientConfig) EnableOverwriteFileValue() bool {
	return boolValue(c.EnableOverwriteFile)
}

func (c ClientConfig) LangValue() string {
	return stringValue(c.Lang)
}

func (c ClientConfig) ShowVersionValue() bool {
	return boolValue(c.ShowVersion)
}

func (c ClientConfig) InitializeValue() bool {
	return boolValue(c.Initialize)
}

type LocaleConfig struct {
	Lang     string `json:"lang,omitempty"`
	Language string `json:"language,omitempty"`
	LcAll    string `json:"lcAll,omitempty"`
}

var DefaultLocale = LocaleConfig{
	Lang:     "ja_JP.UTF-8",
	Language: "ja_JP:ja",
	LcAll:    "ja_JP.UTF-8",
}

type JsonEntry struct {
	Image            string        `json:"image,omitempty"`
	RunCommand       string        `json:"runCommand,omitempty"`
	Locale           *LocaleConfig `json:"locale,omitempty"`
	Timezone         string        `json:"timezone,omitempty"`
	VSCodeExtensions []string      `json:"vscodeExtensions,omitempty"`
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
	RunCommand       string
	VSCodeExtensions []string
	OsModules        OsModules
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
