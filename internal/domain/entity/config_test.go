package entity

import (
	"testing"
)

func TestClientConfig_StringAccessors_ReturnsEmptyWhenNil(t *testing.T) {
	c := ClientConfig{}
	if got := c.OutputDirValue(); got != "" {
		t.Errorf("OutputDirValue: got %q, want %q", got, "")
	}
	if got := c.ContainerNameValue(); got != "" {
		t.Errorf("ContainerNameValue: got %q, want %q", got, "")
	}
	if got := c.ServiceNameValue(); got != "" {
		t.Errorf("ServiceNameValue: got %q, want %q", got, "")
	}
	if got := c.LanguageValue(); got != "" {
		t.Errorf("LanguageValue: got %q, want %q", got, "")
	}
	if got := c.WorkspaceFolderValue(); got != "" {
		t.Errorf("WorkspaceFolderValue: got %q, want %q", got, "")
	}
	if got := c.TimezoneValue(); got != "" {
		t.Errorf("TimezoneValue: got %q, want %q", got, "")
	}
	if got := c.ImageConfigValue(); got != "" {
		t.Errorf("ImageConfigValue: got %q, want %q", got, "")
	}
	if got := c.PortValue(); got != "" {
		t.Errorf("PortValue: got %q, want %q", got, "")
	}
	if got := c.ComposeFileValue(); got != "" {
		t.Errorf("ComposeFileValue: got %q, want %q", got, "")
	}
	if got := c.LangValue(); got != "" {
		t.Errorf("LangValue: got %q, want %q", got, "")
	}
}

func TestClientConfig_StringAccessors_ReturnsValueWhenSet(t *testing.T) {
	str := func(s string) *string { return &s }

	c := ClientConfig{
		OutputDir:       str("/out"),
		ContainerName:   str("mycontainer"),
		ServiceName:     str("myservice"),
		Language:        str("python"),
		WorkspaceFolder: str("/workspace"),
		Timezone:        str("Asia/Tokyo"),
		ImageConfig:     str("https://example.com/config.json"),
		Port:            str("3000"),
		ComposeFile:     str("docker-compose.yaml"),
		Lang:            str("ja"),
	}

	if got := c.OutputDirValue(); got != "/out" {
		t.Errorf("OutputDirValue: got %q, want %q", got, "/out")
	}
	if got := c.ContainerNameValue(); got != "mycontainer" {
		t.Errorf("ContainerNameValue: got %q, want %q", got, "mycontainer")
	}
	if got := c.ServiceNameValue(); got != "myservice" {
		t.Errorf("ServiceNameValue: got %q, want %q", got, "myservice")
	}
	if got := c.LanguageValue(); got != "python" {
		t.Errorf("LanguageValue: got %q, want %q", got, "python")
	}
	if got := c.WorkspaceFolderValue(); got != "/workspace" {
		t.Errorf("WorkspaceFolderValue: got %q, want %q", got, "/workspace")
	}
	if got := c.TimezoneValue(); got != "Asia/Tokyo" {
		t.Errorf("TimezoneValue: got %q, want %q", got, "Asia/Tokyo")
	}
	if got := c.ImageConfigValue(); got != "https://example.com/config.json" {
		t.Errorf("ImageConfigValue: got %q, want %q", got, "https://example.com/config.json")
	}
	if got := c.PortValue(); got != "3000" {
		t.Errorf("PortValue: got %q, want %q", got, "3000")
	}
	if got := c.ComposeFileValue(); got != "docker-compose.yaml" {
		t.Errorf("ComposeFileValue: got %q, want %q", got, "docker-compose.yaml")
	}
	if got := c.LangValue(); got != "ja" {
		t.Errorf("LangValue: got %q, want %q", got, "ja")
	}
}

func TestClientConfig_BoolAccessors_ReturnsFalseWhenNil(t *testing.T) {
	c := ClientConfig{}
	if got := c.EnableOverwriteFileValue(); got != false {
		t.Errorf("EnableOverwriteFileValue: got %v, want false", got)
	}
	if got := c.ShowVersionValue(); got != false {
		t.Errorf("ShowVersionValue: got %v, want false", got)
	}
	if got := c.InitializeValue(); got != false {
		t.Errorf("InitializeValue: got %v, want false", got)
	}
}

func TestClientConfig_BoolAccessors_ReturnsTrueWhenSet(t *testing.T) {
	b := func(v bool) *bool { return &v }

	c := ClientConfig{
		EnableOverwriteFile: b(true),
		ShowVersion:         b(true),
		Initialize:          b(true),
	}

	if got := c.EnableOverwriteFileValue(); got != true {
		t.Errorf("EnableOverwriteFileValue: got %v, want true", got)
	}
	if got := c.ShowVersionValue(); got != true {
		t.Errorf("ShowVersionValue: got %v, want true", got)
	}
	if got := c.InitializeValue(); got != true {
		t.Errorf("InitializeValue: got %v, want true", got)
	}
}

func validCodespaceConfig() CodespaceConfig {
	return CodespaceConfig{
		ContainerName:   "mycontainer",
		ServiceName:     "myservice",
		WorkspaceFolder: "/workspace",
		BaseImage:       "alpine:latest",
		ComposeFileName: "docker-compose.yaml",
	}
}

func TestCodespaceConfig_Validate_SucceedsWithAllRequiredFields(t *testing.T) {
	c := validCodespaceConfig()
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCodespaceConfig_Validate_ErrorWhenContainerNameEmpty(t *testing.T) {
	c := validCodespaceConfig()
	c.ContainerName = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty ContainerName, got nil")
	}
}

func TestCodespaceConfig_Validate_ErrorWhenServiceNameEmpty(t *testing.T) {
	c := validCodespaceConfig()
	c.ServiceName = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty ServiceName, got nil")
	}
}

func TestCodespaceConfig_Validate_ErrorWhenWorkspaceFolderEmpty(t *testing.T) {
	c := validCodespaceConfig()
	c.WorkspaceFolder = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty WorkspaceFolder, got nil")
	}
}

func TestCodespaceConfig_Validate_ErrorWhenBaseImageEmpty(t *testing.T) {
	c := validCodespaceConfig()
	c.BaseImage = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty BaseImage, got nil")
	}
}

func TestCodespaceConfig_Validate_ErrorWhenComposeFileNameEmpty(t *testing.T) {
	c := validCodespaceConfig()
	c.ComposeFileName = ""
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty ComposeFileName, got nil")
	}
}
