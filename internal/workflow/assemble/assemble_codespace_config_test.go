package assemble

import (
	"encoding/json"
	"errors"
	"testing"

	"codespacegen/internal/domain/entity"
)

// fakeConfigResolver は ConfigResolver のテスト用実装。
type fakeConfigResolver struct {
	projectName     string
	language        string
	workspaceFolder string
	serviceName     string
	portMapping     string
	timezone        string
	mergedEntries   map[string]entity.LangEntry
	baseImageEntry  entity.LangEntry

	errProjectName     error
	errLanguage        error
	errWorkspaceFolder error
	errServiceName     error
	errPortMapping     error
	errTimezone        error
	errMergeEntries    error
	errBaseImage       error
}

func (f *fakeConfigResolver) ResolveProjectName(_ string) (string, error) {
	return f.projectName, f.errProjectName
}

func (f *fakeConfigResolver) ResolveLanguage(_ string) (string, error) {
	return f.language, f.errLanguage
}

func (f *fakeConfigResolver) ResolveWorkspaceFolder(_ string) (string, error) {
	return f.workspaceFolder, f.errWorkspaceFolder
}

func (f *fakeConfigResolver) ResolveServiceName(_ string) (string, error) {
	return f.serviceName, f.errServiceName
}

func (f *fakeConfigResolver) ResolvePortMapping(_ string) (string, error) {
	return f.portMapping, f.errPortMapping
}

func (f *fakeConfigResolver) ResolveTimezone(_, _, _ string) (string, error) {
	return f.timezone, f.errTimezone
}

func (f *fakeConfigResolver) MergeLanguageEntries(_ map[string]json.RawMessage) (map[string]entity.LangEntry, error) {
	return f.mergedEntries, f.errMergeEntries
}

func (f *fakeConfigResolver) ResolveBaseImage(_, _, _ string, _ map[string]entity.LangEntry, _ string) (entity.LangEntry, error) {
	return f.baseImageEntry, f.errBaseImage
}

// defaultFakeResolver は正常系で使う最小限の fakeConfigResolver を返す。
func defaultFakeResolver() *fakeConfigResolver {
	return &fakeConfigResolver{
		projectName:     "myproject",
		language:        "python",
		workspaceFolder: "/workspace",
		serviceName:     "app",
		portMapping:     "",
		timezone:        "UTC",
		mergedEntries:   map[string]entity.LangEntry{},
		baseImageEntry:  entity.LangEntry{Image: "python:3.12", RunCommand: "pip install -r requirements.txt"},
	}
}

func strPtr(s string) *string { return &s }

// --- Resolve 正常系 ---

func TestAssembleCodespaceConfig_Resolve_BuildsConfigCorrectly(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.baseImageEntry = entity.LangEntry{
		Image:            "python:3.12",
		RunCommand:       "pip install -r requirements.txt",
		Timezone:         "Asia/Tokyo",
		VSCodeExtensions: []string{"ms-python.python"},
	}
	resolver.timezone = "Asia/Tokyo"

	composeName := "docker-compose.yaml"
	clientConfig := entity.ClientConfig{
		ComposeFile: strPtr(composeName),
	}

	rcc := NewAssembleCodespaceConfig(resolver)
	got, err := rcc.Resolve(clientConfig, entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}, nil)

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
	got, err := rcc.Resolve(entity.ClientConfig{}, entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.PortMapping != "3000:3000" {
		t.Errorf("PortMapping: got %q, want %q", got.PortMapping, "3000:3000")
	}
}

// --- Resolve エラー伝播 ---

func TestAssembleCodespaceConfig_Resolve_ErrorFromResolveProjectName(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errProjectName = errors.New("project name error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromResolveLanguage(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errLanguage = errors.New("language error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromResolveWorkspaceFolder(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errWorkspaceFolder = errors.New("workspace folder error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromResolveServiceName(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errServiceName = errors.New("service name error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromResolvePortMapping(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errPortMapping = errors.New("port mapping error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromMergeLanguageEntries(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errMergeEntries = errors.New("merge error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromResolveBaseImage(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errBaseImage = errors.New("base image error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAssembleCodespaceConfig_Resolve_ErrorFromResolveTimezone(t *testing.T) {
	resolver := defaultFakeResolver()
	resolver.errTimezone = errors.New("timezone error")

	rcc := NewAssembleCodespaceConfig(resolver)
	_, err := rcc.Resolve(entity.ClientConfig{}, entity.DefaultSetting{
		Image:    "alpine:latest",
		Timezone: "UTC",
	}, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
