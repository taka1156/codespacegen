package resolve

import (
	"strings"
	"testing"

	"codespacegen/internal/domain/entity"
	"codespacegen/internal/utils"
)

// newResolver は strings.NewReader を使って入力を注入したリゾルバを返す。
func newResolver(input string) *CodespaceConfigResolver {
	return NewCodespaceConfigResolver(strings.NewReader(input))
}

// --- ResolveLanguage ---

func TestResolveLanguage_UsesExplicitValueWhenUserAccepts(t *testing.T) {
	r := newResolver("\n") // Enter キーのみ → デフォルト採用
	got, err := r.ResolveLanguage("Python")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "python" {
		t.Errorf("got %q, want %q", got, "python")
	}
}

func TestResolveLanguage_UserOverridesWithInput(t *testing.T) {
	r := newResolver("Rust\n")
	got, err := r.ResolveLanguage("python")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "rust" {
		t.Errorf("got %q, want %q", got, "rust")
	}
}

func TestResolveLanguage_NoExplicitUserInputsValue(t *testing.T) {
	r := newResolver("go\n")
	got, err := r.ResolveLanguage("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "go" {
		t.Errorf("got %q, want %q", got, "go")
	}
}

// --- ResolveWorkspaceFolder ---

func TestResolveWorkspaceFolder_UsesExplicitValue(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolveWorkspaceFolder("/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/app" {
		t.Errorf("got %q, want %q", got, "/app")
	}
}

func TestResolveWorkspaceFolder_DefaultsToWorkspaceWhenEmpty(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolveWorkspaceFolder("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/workspace" {
		t.Errorf("got %q, want %q", got, "/workspace")
	}
}

func TestResolveWorkspaceFolder_UserOverrides(t *testing.T) {
	r := newResolver("/custom\n")
	got, err := r.ResolveWorkspaceFolder("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/custom" {
		t.Errorf("got %q, want %q", got, "/custom")
	}
}

// --- ResolveServiceName ---

func TestResolveServiceName_UsesExplicitValue(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolveServiceName("myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "myapp" {
		t.Errorf("got %q, want %q", got, "myapp")
	}
}

func TestResolveServiceName_DefaultsToAppWhenEmpty(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolveServiceName("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "app" {
		t.Errorf("got %q, want %q", got, "app")
	}
}

// --- ResolveTimezone ---

func TestResolveTimezone_UsesExplicitTimezone(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolveTimezone("Asia/Tokyo", "", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Asia/Tokyo" {
		t.Errorf("got %q, want %q", got, "Asia/Tokyo")
	}
}

func TestResolveTimezone_FallsBackToConfigTimezone(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolveTimezone("", "America/New_York", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "America/New_York" {
		t.Errorf("got %q, want %q", got, "America/New_York")
	}
}

func TestResolveTimezone_FallsBackToDefaultTimezone(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolveTimezone("", "", "UTC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "UTC" {
		t.Errorf("got %q, want %q", got, "UTC")
	}
}

func TestResolveTimezone_FallsBackToEntityDefaultWhenAllEmpty(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolveTimezone("", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != entity.DefaultTimezone {
		t.Errorf("got %q, want %q", got, entity.DefaultTimezone)
	}
}

// --- ResolveProjectName ---

func TestResolveProjectName_UsesExplicitValue(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolveProjectName("myproject")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "myproject" {
		t.Errorf("got %q, want %q", got, "myproject")
	}
}

func TestResolveProjectName_UserInputsValue(t *testing.T) {
	r := newResolver("hello\n")
	got, err := r.ResolveProjectName("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
}

func TestResolveProjectName_EOFWithNoInputReturnsError(t *testing.T) {
	r := newResolver("") // EOF immediately
	_, err := r.ResolveProjectName("")
	if err == nil {
		t.Fatal("expected error for empty project name at EOF, got nil")
	}
}

// --- ResolvePortMapping ---

func TestResolvePortMapping_EmptyPorts(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolvePortMapping("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestResolvePortMapping_NormalizesPortOnlyToMapping(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolvePortMapping("3000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "3000:3000" {
		t.Errorf("got %q, want %q", got, "3000:3000")
	}
}

func TestResolvePortMapping_AcceptsFullMappingFormat(t *testing.T) {
	r := newResolver("\n")
	got, err := r.ResolvePortMapping("8080:9090")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "8080:9090" {
		t.Errorf("got %q, want %q", got, "8080:9090")
	}
}

func TestResolvePortMapping_UserInputsPort(t *testing.T) {
	r := newResolver("4000\n")
	got, err := r.ResolvePortMapping("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "4000:4000" {
		t.Errorf("got %q, want %q", got, "4000:4000")
	}
}

func TestResolvePortMapping_RetriesOnInvalidThenAcceptsValid(t *testing.T) {
	r := newResolver("bad\n3000\n")
	got, err := r.ResolvePortMapping("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "3000:3000" {
		t.Errorf("got %q, want %q", got, "3000:3000")
	}
}

func TestResolvePortMapping_InvalidDefaultPortRetriesUntilValid(t *testing.T) {
	// explicitPort が不正フォーマットのとき、Enter を押すとリトライし、次の有効な入力を受け付ける
	r := newResolver("\n5000\n")
	got, err := r.ResolvePortMapping("bad")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "5000:5000" {
		t.Errorf("got %q, want %q", got, "5000:5000")
	}
}

// --- ResolveBaseImage ---

func TestResolveBaseImage_ExplicitImageTakesPriority(t *testing.T) {
	r := newResolver("")
	got, err := r.ResolveBaseImage("python", "custom:latest", nil, "alpine:latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Image != "custom:latest" {
		t.Errorf("got %q, want %q", got.Image, "custom:latest")
	}
}

func TestResolveBaseImage_EmptyLanguageReturnsDefaultImage(t *testing.T) {
	r := newResolver("")
	got, err := r.ResolveBaseImage("", "", nil, "ubuntu:22.04")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Image != "ubuntu:22.04" {
		t.Errorf("got %q, want %q", got.Image, "ubuntu:22.04")
	}
}

func TestResolveBaseImage_EmptyLanguageAndDefaultImageFallsBackToConstant(t *testing.T) {
	r := newResolver("")
	got, err := r.ResolveBaseImage("", "", nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Image != entity.DefaultImage {
		t.Errorf("got %q, want %q", got.Image, entity.DefaultImage)
	}
}

func TestResolveBaseImage_LanguageFoundInEntries(t *testing.T) {
	r := newResolver("")
	entries := map[string]entity.LangEntry{
		"python": {Image: "python:3.12", RunCommand: utils.Ptr("pip install -r requirements.txt")},
	}
	got, err := r.ResolveBaseImage("python", "", entries, "alpine:latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Image != "python:3.12" {
		t.Errorf("got %q, want %q", got.Image, "python:3.12")
	}
	if got.RunCommand == nil || *got.RunCommand != "pip install -r requirements.txt" {
		t.Errorf("RunCommand: got %q, want %q", *got.RunCommand, "pip install -r requirements.txt")
	}
}

func TestResolveBaseImage_LanguageNotFoundReturnsError(t *testing.T) {
	r := newResolver("")
	entries := map[string]entity.LangEntry{}
	_, err := r.ResolveBaseImage("python", "", entries, "alpine:latest")
	if err == nil {
		t.Fatal("expected error for unsupported language, got nil")
	}
}

func TestResolveBaseImage_LanguageCaseInsensitive(t *testing.T) {
	r := newResolver("")
	entries := map[string]entity.LangEntry{
		"python": {Image: "python:3.12"},
	}
	got, err := r.ResolveBaseImage("Python", "", entries, "alpine:latest")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Image != "python:3.12" {
		t.Errorf("got %q, want %q", got.Image, "python:3.12")
	}
}

func TestResolveBaseImage_EntryWithEmptyImageReturnsError(t *testing.T) {
	r := newResolver("")
	entries := map[string]entity.LangEntry{
		"python": {Image: ""},
	}
	_, err := r.ResolveBaseImage("python", "", entries, "alpine:latest")
	if err == nil {
		t.Fatal("expected error for entry with empty image, got nil")
	}
}

// --- MergeLanguageEntries ---

func TestMergeLanguageEntries_BothNilReturnsEmpty(t *testing.T) {
	r := newResolver("")
	got, err := r.MergeLanguageEntries(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestMergeLanguageEntries_EmptyOverrides(t *testing.T) {
	r := newResolver("")
	got, err := r.MergeLanguageEntries(&entity.CommonEntry{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestMergeLanguageEntries_NilCommonPassesThroughLanguages(t *testing.T) {
	r := newResolver("")
	langEntries := map[string]*entity.LangEntry{
		"python": {Image: "python:3.12"},
	}
	got, err := r.MergeLanguageEntries(nil, langEntries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry, ok := got["python"]
	if !ok {
		t.Fatal("expected python key in result")
	}
	if entry.Image != "python:3.12" {
		t.Errorf("Image: got %q, want %q", entry.Image, "python:3.12")
	}
}

func TestMergeLanguageEntries_CommonAppliedAsBase(t *testing.T) {
	r := newResolver("")
	commonEntry := &entity.CommonEntry{
		VSCodeExtensions: utils.Ptr([]string{"GitHub.copilot"}),
	}
	langEntries := map[string]*entity.LangEntry{
		"python": {Image: "python:3.12"},
	}
	got, err := r.MergeLanguageEntries(commonEntry, langEntries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry, ok := got["python"]
	if !ok {
		t.Fatal("expected python key in result")
	}
	if entry.VSCodeExtensions == nil || len(*entry.VSCodeExtensions) != 1 || (*entry.VSCodeExtensions)[0] != "GitHub.copilot" {
		t.Errorf("VSCodeExtensions: got %v, want [GitHub.copilot]", entry.VSCodeExtensions)
	}
}

func TestMergeLanguageEntries_LanguageTimezoneOverridesCommon(t *testing.T) {
	r := newResolver("")
	commonEntry := &entity.CommonEntry{
		Timezone: utils.Ptr("UTC"),
	}
	langEntries := map[string]*entity.LangEntry{
		"python": {Image: "python:3.12", Timezone: utils.Ptr("Asia/Tokyo")},
	}
	got, err := r.MergeLanguageEntries(commonEntry, langEntries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := got["python"]
	if entry.Timezone == nil || *entry.Timezone != "Asia/Tokyo" {
		t.Errorf("Timezone: got %v, want Asia/Tokyo", entry.Timezone)
	}
}

func TestMergeLanguageEntries_LanguageOverridesCommon(t *testing.T) {
	r := newResolver("")
	commonEntry := &entity.CommonEntry{
		Timezone: utils.Ptr("UTC"),
	}
	langEntries := map[string]*entity.LangEntry{
		"python": {Image: "python:3.12", Timezone: utils.Ptr("Asia/Tokyo"), RunCommand: utils.Ptr("pip install -r requirements.txt")},
	}
	got, err := r.MergeLanguageEntries(commonEntry, langEntries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := got["python"]
	if entry.Image != "python:3.12" {
		t.Errorf("Image: got %q, want %q", entry.Image, "python:3.12")
	}
	if entry.RunCommand == nil || *entry.RunCommand != "pip install -r requirements.txt" {
		t.Errorf("RunCommand: got %q, want %q", *entry.RunCommand, "pip install -r requirements.txt")
	}
}

func TestMergeLanguageEntries_VSCodeExtensionsMergedFromBoth(t *testing.T) {
	r := newResolver("")
	commonEntry := &entity.CommonEntry{
		VSCodeExtensions: utils.Ptr([]string{"GitHub.copilot"}),
	}
	langEntries := map[string]*entity.LangEntry{
		"python": {Image: "python:3.12", VSCodeExtensions: utils.Ptr([]string{"ms-python.python"})},
	}
	got, err := r.MergeLanguageEntries(commonEntry, langEntries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := got["python"]
	if entry.VSCodeExtensions == nil || len(*entry.VSCodeExtensions) != 2 {
		t.Errorf("VSCodeExtensions: got %v, want 2 items", entry.VSCodeExtensions)
	}
}

func TestMergeLanguageEntries_OnlyCommonVSCodeExtensions(t *testing.T) {
	r := newResolver("")
	commonEntry := &entity.CommonEntry{
		VSCodeExtensions: utils.Ptr([]string{"GitHub.copilot"}),
	}
	langEntries := map[string]*entity.LangEntry{
		"python": {Image: "python:3.12"},
	}
	got, err := r.MergeLanguageEntries(commonEntry, langEntries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := got["python"]
	if entry.VSCodeExtensions == nil || len(*entry.VSCodeExtensions) != 1 || (*entry.VSCodeExtensions)[0] != "GitHub.copilot" {
		t.Errorf("VSCodeExtensions: got %v, want [GitHub.copilot]", entry.VSCodeExtensions)
	}
}

func TestMergeLanguageEntries_OnlyLangVSCodeExtensions(t *testing.T) {
	r := newResolver("")
	commonEntry := &entity.CommonEntry{}
	langEntries := map[string]*entity.LangEntry{
		"python": {Image: "python:3.12", VSCodeExtensions: utils.Ptr([]string{"ms-python.python"})},
	}
	got, err := r.MergeLanguageEntries(commonEntry, langEntries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := got["python"]
	if entry.VSCodeExtensions == nil || len(*entry.VSCodeExtensions) != 1 || (*entry.VSCodeExtensions)[0] != "ms-python.python" {
		t.Errorf("VSCodeExtensions: got %v, want [ms-python.python]", entry.VSCodeExtensions)
	}
}

func TestMergeLanguageEntries_KeyNormalizedToLowercase(t *testing.T) {
	r := newResolver("")
	langEntries := map[string]*entity.LangEntry{
		"Python": {Image: "python:3.12"},
	}
	got, err := r.MergeLanguageEntries(&entity.CommonEntry{}, langEntries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["python"]; !ok {
		t.Error("expected key normalized to lowercase 'python'")
	}
	if _, ok := got["Python"]; ok {
		t.Error("expected original 'Python' key to be absent")
	}
}

func TestMergeLanguageEntries_CommonAndSchemaKeysFiltered(t *testing.T) {
	r := newResolver("")
	langEntries := map[string]*entity.LangEntry{
		"python":  {Image: "python:3.12"},
		"common":  {Image: "should-be-filtered"},
		"$schema": {Image: "should-be-filtered"},
	}
	got, err := r.MergeLanguageEntries(&entity.CommonEntry{}, langEntries)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["common"]; ok {
		t.Error("expected 'common' key to be filtered out")
	}
	if _, ok := got["$schema"]; ok {
		t.Error("expected '$schema' key to be filtered out")
	}
	if _, ok := got["python"]; !ok {
		t.Error("expected 'python' key to be present")
	}
}
