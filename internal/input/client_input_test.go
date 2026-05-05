package input

import (
	"testing"

	"github.com/taka1156/codespacegen/internal/domain/entity"
)

func TestNewClientInput_ReturnsNonNil(t *testing.T) {
	ci := NewClientInput()
	if ci == nil {
		t.Fatal("expected non-nil ClientInput")
	}
}

func TestClientInput_GetInput_InitCommand(t *testing.T) {
	ci := NewClientInput()
	args := []string{"codespacegen", "init", "-output", "custom-devcontainer"}

	got := ci.GetInput(args)

	if got.Mode.CommandlineModeValue() != entity.Initialize.CommandlineModeValue() {
		t.Errorf("Mode: got %q, want %q", got.Mode.CommandlineModeValue(), entity.Initialize.CommandlineModeValue())
	}
	if got.OutputDirValue() != "custom-devcontainer" {
		t.Errorf("OutputDir: got %q, want %q", got.OutputDirValue(), "custom-devcontainer")
	}
}

func TestClientInput_GetInput_UpdateCommand(t *testing.T) {
	ci := NewClientInput()
	args := []string{"codespacegen", "update"}

	got := ci.GetInput(args)

	if got.Mode.CommandlineModeValue() != entity.Update.CommandlineModeValue() {
		t.Errorf("Mode: got %q, want %q", got.Mode.CommandlineModeValue(), entity.Update.CommandlineModeValue())
	}
}

func TestClientInput_GetInput_UpdateCommand_WithLang(t *testing.T) {
	ci := NewClientInput()
	args := []string{"codespacegen", "update", "-lang", "ja"}

	got := ci.GetInput(args)

	if got.Mode.CommandlineModeValue() != entity.Update.CommandlineModeValue() {
		t.Errorf("Mode: got %q, want %q", got.Mode.CommandlineModeValue(), entity.Update.CommandlineModeValue())
	}
	if got.LangValue() != "ja" {
		t.Errorf("Lang: got %q, want %q", got.LangValue(), "ja")
	}
}

func TestClientInput_GetInput_VersionCommand(t *testing.T) {
	ci := NewClientInput()
	args := []string{"codespacegen", "version"}

	got := ci.GetInput(args)

	if got.Mode.CommandlineModeValue() != entity.Version.CommandlineModeValue() {
		t.Errorf("Mode: got %q, want %q", got.Mode.CommandlineModeValue(), entity.Version.CommandlineModeValue())
	}
}

func TestClientInput_GetInput_DefaultCommandWithFlags(t *testing.T) {
	ci := NewClientInput()
	args := []string{
		"codespacegen",
		"-output", "custom-output",
		"-name", "my-project",
		"-service", "app",
		"-language", "go",
		"-workspace-folder", "/work",
		"-timezone", "Asia/Tokyo",
		"-image-config", "custom.json",
		"-port", "3000:3000",
		"-compose-file", "compose.yml",
		"-force",
		"-lang", "ja",
		"-headless",
	}

	got := ci.GetInput(args)

	if got.Mode.CommandlineModeValue() != "undefined" {
		t.Errorf("Mode: got %q, want %q", got.Mode.CommandlineModeValue(), "undefined")
	}
	if got.OutputDirValue() != "custom-output" {
		t.Errorf("OutputDir: got %q, want %q", got.OutputDirValue(), "custom-output")
	}
	if got.ContainerNameValue() != "my-project" {
		t.Errorf("ContainerName: got %q, want %q", got.ContainerNameValue(), "my-project")
	}
	if got.ServiceNameValue() != "app" {
		t.Errorf("ServiceName: got %q, want %q", got.ServiceNameValue(), "app")
	}
	if got.LanguageValue() != "go" {
		t.Errorf("Language: got %q, want %q", got.LanguageValue(), "go")
	}
	if got.WorkspaceFolderValue() != "/work" {
		t.Errorf("WorkspaceFolder: got %q, want %q", got.WorkspaceFolderValue(), "/work")
	}
	if got.TimezoneValue() != "Asia/Tokyo" {
		t.Errorf("Timezone: got %q, want %q", got.TimezoneValue(), "Asia/Tokyo")
	}
	if got.ImageConfigValue() != "custom.json" {
		t.Errorf("ImageConfig: got %q, want %q", got.ImageConfigValue(), "custom.json")
	}
	if got.PortValue() != "3000:3000" {
		t.Errorf("Port: got %q, want %q", got.PortValue(), "3000:3000")
	}
	if got.ComposeFileValue() != "compose.yml" {
		t.Errorf("ComposeFile: got %q, want %q", got.ComposeFileValue(), "compose.yml")
	}
	if !got.EnableOverwriteFileValue() {
		t.Error("EnableOverwriteFile: got false, want true")
	}
	if got.LangValue() != "ja" {
		t.Errorf("Lang: got %q, want %q", got.LangValue(), "ja")
	}
	if !got.HeadlessValue() {
		t.Error("Headless: got false, want true")
	}
}

func TestClientInput_GetInput_DefaultCommandWithoutFlags_UsesDefaults(t *testing.T) {
	ci := NewClientInput()
	args := []string{"codespacegen"}

	got := ci.GetInput(args)

	if got.Mode.CommandlineModeValue() != "undefined" {
		t.Errorf("Mode: got %q, want %q", got.Mode.CommandlineModeValue(), "undefined")
	}
	if got.OutputDirValue() != ".devcontainer" {
		t.Errorf("OutputDir: got %q, want %q", got.OutputDirValue(), ".devcontainer")
	}
	if got.ImageConfigValue() != "codespacegen.json" {
		t.Errorf("ImageConfig: got %q, want %q", got.ImageConfigValue(), "codespacegen.json")
	}
	if got.ComposeFileValue() != "docker-compose.yaml" {
		t.Errorf("ComposeFile: got %q, want %q", got.ComposeFileValue(), "docker-compose.yaml")
	}
	if got.WorkspaceFolderValue() != "/workspace" {
		t.Errorf("WorkspaceFolder: got %q, want %q", got.WorkspaceFolderValue(), "/workspace")
	}
	if got.EnableOverwriteFileValue() {
		t.Error("EnableOverwriteFile: got true, want false")
	}
	if got.HeadlessValue() {
		t.Error("Headless: got true, want false")
	}
}
