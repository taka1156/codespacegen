package assemble

import (
	"errors"
	"testing"

	"github.com/taka1156/codespacegen/internal/domain/entity"
	"github.com/taka1156/codespacegen/internal/utils"
)

type fakeCodespacePromptResolver struct {
	projectName     string
	language        string
	workspaceFolder string
	serviceName     string
	portMapping     string
	timezone        string

	errProjectName     error
	errLanguage        error
	errWorkspaceFolder error
	errServiceName     error
	errPortMapping     error
	errTimezone        error
}

func (f *fakeCodespacePromptResolver) PromptProjectName(_ string) (string, error) {
	return f.projectName, f.errProjectName
}

func (f *fakeCodespacePromptResolver) PromptLanguage(_ string) (string, error) {
	return f.language, f.errLanguage
}

func (f *fakeCodespacePromptResolver) PromptWorkspaceFolder(_ string) (string, error) {
	return f.workspaceFolder, f.errWorkspaceFolder
}

func (f *fakeCodespacePromptResolver) PromptServiceName(_ string) (string, error) {
	return f.serviceName, f.errServiceName
}

func (f *fakeCodespacePromptResolver) PromptPortMapping(_ string) (string, error) {
	return f.portMapping, f.errPortMapping
}

func (f *fakeCodespacePromptResolver) PromptTimezone(_ string) (string, error) {
	return f.timezone, f.errTimezone
}

var (
	defaultTestSetting = entity.DefaultSetting{Image: "alpine:latest", Timezone: "UTC"}
	defaultJsonConfig  = entity.JsonConfig{
		Common: &entity.CommonEntry{},
		Langs: []*entity.LangEntry{
			{
				ProfileName:      "python",
				Image:            "python:3.12",
				RunCommand:       utils.Ptr("pip install -r requirements.txt"),
				VSCodeExtensions: utils.Ptr([]string{"ms-python.python"}),
			},
		},
	}
)

func defaultFakeResolver() *fakeCodespacePromptResolver {
	return &fakeCodespacePromptResolver{
		projectName:     "myproject",
		language:        "python",
		workspaceFolder: "/workspace",
		serviceName:     "app",
		timezone:        "Asia/Tokyo",
	}
}

func TestAssembleCodespaceConfig_Resolve_BuildsConfigCorrectly(t *testing.T) {
	resolver := defaultFakeResolver()
	composeName := "docker-compose.yaml"
	clientConfig := entity.ClientConfig{
		ComposeFile: utils.Ptr(composeName),
	}

	rcc := NewAssembleCodespaceConfig(resolver)
	got, err := rcc.Resolve(clientConfig, defaultTestSetting, defaultJsonConfig)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ContainerName != "myproject" {
		t.Errorf("ContainerName: got %q, want %q", got.ContainerName, "myproject")
	}
	if got.ServiceName != "app" {
		t.Errorf("ServiceName: got %q, want %q", got.ServiceName, "app")
	}
	if got.WorkspaceFolder != "/workspace" {
		t.Errorf("WorkspaceFolder: got %q, want %q", got.WorkspaceFolder, "/workspace")
	}
	if got.BaseImage != "python:3.12" {
		t.Errorf("BaseImage: got %q, want %q", got.BaseImage, "python:3.12")
	}
	if got.Timezone != "Asia/Tokyo" {
		t.Errorf("Timezone: got %q, want %q", got.Timezone, "Asia/Tokyo")
	}
	if got.ComposeFileName != composeName {
		t.Errorf("ComposeFileName: got %q, want %q", got.ComposeFileName, composeName)
	}
	if got.RunCommand != "pip install -r requirements.txt" {
		t.Errorf("RunCommand: got %q, want %q", got.RunCommand, "pip install -r requirements.txt")
	}
	if len(got.VSCodeExtensions) != 1 || got.VSCodeExtensions[0] != "ms-python.python" {
		t.Errorf("VSCodeExtensions: got %v, want [ms-python.python]", got.VSCodeExtensions)
	}
}

func TestAssembleCodespaceConfig_Resolve_PortMappingIsSet(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.portMapping = "3000:3000"

	rcc := NewAssembleCodespaceConfig(resolver)
	got, err := rcc.Resolve(entity.ClientConfig{}, defaultTestSetting, defaultJsonConfig)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.PortMapping != "3000:3000" {
		t.Errorf("PortMapping: got %q, want %q", got.PortMapping, "3000:3000")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromPromptProjectName(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errProjectName = errors.New("project name error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, defaultTestSetting, defaultJsonConfig)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromPromptLanguage(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errLanguage = errors.New("language error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, defaultTestSetting, defaultJsonConfig)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromPromptWorkspaceFolder(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errWorkspaceFolder = errors.New("workspace folder error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, defaultTestSetting, defaultJsonConfig)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromPromptServiceName(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errServiceName = errors.New("service name error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, defaultTestSetting, defaultJsonConfig)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromPromptPortMapping(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errPortMapping = errors.New("port mapping error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, defaultTestSetting, defaultJsonConfig)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromPromptTimezone(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errTimezone = errors.New("timezone error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, defaultTestSetting, defaultJsonConfig)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromUnknownLanguage(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.language = "unknown-language"

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, defaultTestSetting, defaultJsonConfig)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_HeadlessMode(t *testing.T) {
	// Prepare ClientConfig for headless mode
	clientConfig := entity.ClientConfig{
		Headless:        utils.Ptr(true),
		ContainerName:   utils.Ptr("headless-project"),
		Language:        utils.Ptr("python"),
		WorkspaceFolder: utils.Ptr("/workspace"),
		ServiceName:     utils.Ptr("svc"),
		Port:            utils.Ptr("8080:8080"),
		Timezone:        utils.Ptr("Asia/Tokyo"),
		ComposeFile:     utils.Ptr("docker-compose.yml"),
	}

	rcc := NewAssembleCodespaceConfig(defaultFakeResolver())
	got, err := rcc.Resolve(clientConfig, defaultTestSetting, defaultJsonConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ContainerName != "headless-project" {
		t.Errorf("ContainerName: got %q, want %q", got.ContainerName, "headless-project")
	}
	if got.ServiceName != "svc" {
		t.Errorf("ServiceName: got %q, want %q", got.ServiceName, "svc")
	}
	if got.WorkspaceFolder != "/workspace" {
		t.Errorf("WorkspaceFolder: got %q, want %q", got.WorkspaceFolder, "/workspace")
	}
	if got.BaseImage != "python:3.12" {
		t.Errorf("BaseImage: got %q, want %q", got.BaseImage, "python:3.12")
	}
	if got.Timezone != "Asia/Tokyo" {
		t.Errorf("Timezone: got %q, want %q", got.Timezone, "Asia/Tokyo")
	}
	if got.ComposeFileName != "docker-compose.yml" {
		t.Errorf("ComposeFileName: got %q, want %q", got.ComposeFileName, "docker-compose.yml")
	}
	if got.PortMapping != "8080:8080" {
		t.Errorf("PortMapping: got %q, want %q", got.PortMapping, "8080:8080")
	}
	if got.RunCommand != "pip install -r requirements.txt" {
		t.Errorf("RunCommand: got %q, want %q", got.RunCommand, "pip install -r requirements.txt")
	}
	if len(got.VSCodeExtensions) != 1 || got.VSCodeExtensions[0] != "ms-python.python" {
		t.Errorf("VSCodeExtensions: got %v, want [ms-python.python]", got.VSCodeExtensions)
	}
}
