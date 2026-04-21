package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"codespacegen/internal/domain/entity"
)

// ─── resolveTimezone ───────────────────────────────────────────────────────────

func TestResolveTimezone_UsesDefaultWhenEmpty(t *testing.T) {
	timezone, err := resolveTimezone("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if timezone != entity.DefaultTimezone {
		t.Fatalf("expected %s but got: %s", entity.DefaultTimezone, timezone)
	}
}

func TestResolveTimezone_TrimsExplicitValue(t *testing.T) {
	timezone, err := resolveTimezone("  Europe/Berlin  ", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if timezone != "Europe/Berlin" {
		t.Fatalf("unexpected timezone: %s", timezone)
	}
}

func TestResolveTimezone_UsesConfigTimezoneWhenFlagEmpty(t *testing.T) {
	timezone, err := resolveTimezone("", "  America/Los_Angeles  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if timezone != "America/Los_Angeles" {
		t.Fatalf("unexpected timezone: %s", timezone)
	}
}

// 明示フラグはコンフィグのタイムゾーンより優先される。
func TestResolveTimezone_ExplicitFlagTakesPriorityOverConfig(t *testing.T) {
	timezone, err := resolveTimezone("Asia/Tokyo", "Europe/Berlin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if timezone != "Asia/Tokyo" {
		t.Fatalf("expected Asia/Tokyo but got: %s", timezone)
	}
}

// ─── normalizePortMapping ──────────────────────────────────────────────────────

func TestNormalizePortMapping_PortOnlyExpandsToBothSides(t *testing.T) {
	got, err := normalizePortMapping("3000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "3000:3000" {
		t.Fatalf("expected 3000:3000 but got: %s", got)
	}
}

func TestNormalizePortMapping_PortMappingPassesThrough(t *testing.T) {
	got, err := normalizePortMapping("8080:3000")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "8080:3000" {
		t.Fatalf("expected 8080:3000 but got: %s", got)
	}
}

func TestNormalizePortMapping_InvalidValueReturnsError(t *testing.T) {
	_, err := normalizePortMapping("abc")
	if err == nil {
		t.Fatal("expected error for invalid port mapping")
	}
}

func TestNormalizePortMapping_TrimsWhitespace(t *testing.T) {
	got, err := normalizePortMapping("  3000  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "3000:3000" {
		t.Fatalf("expected 3000:3000 but got: %s", got)
	}
}

// ─── resolveBaseImage ──────────────────────────────────────────────────────────

func TestResolveBaseImage_MoonbitFromConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{
  "moonbit": {
    "image": "ubuntu:24.04",
    "install": "curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash",
    "timezone": "Europe/Berlin"
  }
}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	entry, err := resolveBaseImage("moonbit", "", cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Image != "ubuntu:24.04" {
		t.Fatalf("unexpected image: %s", entry.Image)
	}
	if entry.Install != "curl -fsSL https://cli.moonbitlang.com/install/unix.sh | bash" {
		t.Fatalf("unexpected install command: %s", entry.Install)
	}
	if entry.Timezone != "Europe/Berlin" {
		t.Fatalf("unexpected timezone: %s", entry.Timezone)
	}
	if len(entry.VSCodeExtensions) != 0 {
		t.Fatalf("unexpected extensions: %#v", entry.VSCodeExtensions)
	}
	if entry.Locale != (entity.LocaleConfig{}) {
		t.Fatalf("expected empty locale, got: %#v", entry.Locale)
	}
}

func TestResolveBaseImage_ReturnsLocaleFromConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{
  "go": {
    "image": "golang:1.24-alpine",
    "locale": {
      "lang": "en_US.UTF-8",
      "language": "en_US:en",
      "lcAll": "en_US.UTF-8"
    }
  }
}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	entry, err := resolveBaseImage("go", "", cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := entity.LocaleConfig{Lang: "en_US.UTF-8", Language: "en_US:en", LcAll: "en_US.UTF-8"}
	if entry.Locale != want {
		t.Fatalf("unexpected locale: %#v", entry.Locale)
	}
}

// explicitBaseImage が指定された場合は config ファイルを読まずに早期リターンする。
func TestResolveBaseImage_ExplicitBaseImageOverridesConfig(t *testing.T) {
	entry, err := resolveBaseImage("moonbit", "ubuntu:latest", "codespacegen.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Image != "ubuntu:latest" {
		t.Fatalf("unexpected image: %s", entry.Image)
	}
	if entry.Install != "" {
		t.Fatalf("expected empty install command, got: %s", entry.Install)
	}
	if entry.Timezone != "" {
		t.Fatalf("expected empty timezone, got: %s", entry.Timezone)
	}
	if len(entry.VSCodeExtensions) != 0 {
		t.Fatalf("expected empty extensions, got: %#v", entry.VSCodeExtensions)
	}
	if entry.Locale != (entity.LocaleConfig{}) {
		t.Fatalf("expected empty locale, got: %#v", entry.Locale)
	}
}

// language が空の場合は DefaultImage を返し、config ファイルを参照しない。
func TestResolveBaseImage_EmptyLanguageReturnsDefaultImage(t *testing.T) {
	entry, err := resolveBaseImage("", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Image != entity.DefaultImage {
		t.Fatalf("expected %s but got: %s", entity.DefaultImage, entry.Image)
	}
	if entry.Install != "" || entry.Timezone != "" {
		t.Fatalf("expected empty install and timezone for default image")
	}
	if len(entry.VSCodeExtensions) != 0 {
		t.Fatalf("expected empty extensions, got: %#v", entry.VSCodeExtensions)
	}
}

// config に存在しないキーは "unsupported language" エラーを返す。
func TestResolveBaseImage_UnsupportedLanguageReturnsError(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{"go": {"image": "golang:1.24-alpine"}}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	_, err := resolveBaseImage("unknownlang", "", cfgPath)
	if err == nil {
		t.Fatal("expected error for unsupported language")
	}
	if !strings.Contains(err.Error(), "unsupported language") {
		t.Fatalf("expected 'unsupported language' in error, got: %v", err)
	}
}

func TestResolveBaseImage_ReturnsErrorWhenImageIsMissing(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{"go": {"timezone": "UTC"}}`
	os.WriteFile(cfgPath, []byte(cfg), 0o644)

	_, err := resolveBaseImage("go", "", cfgPath)
	if err == nil {
		t.Fatal("expected error when image is missing")
	}
}

// ─── loadLanguageImages ────────────────────────────────────────────────────────

// common のフィールドは各言語エントリにマージされ、タイムゾーンは言語側で上書きされ、
// vscodeExtensions は [common...] + [language...] の順で結合される。
func TestLoadLanguageImages_MergesCommonAndLanguageSettings(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{
  "common": {
    "timezone": "Asia/Tokyo",
    "vscodeExtensions": [
      "MS-CEINTL.vscode-language-pack-ja"
    ]
  },
  "go": {
    "timezone": "UTC",
    "vscodeExtensions": [
      "golang.Go",
      "MS-CEINTL.vscode-language-pack-ja"
    ]
  }
}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	entries, err := loadLanguageImages(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	goEntry, ok := entries["go"]
	if !ok {
		t.Fatal("expected go entry")
	}
	if goEntry.Timezone != "UTC" {
		t.Fatalf("unexpected timezone: %s", goEntry.Timezone)
	}
	want := []string{"MS-CEINTL.vscode-language-pack-ja", "golang.Go", "MS-CEINTL.vscode-language-pack-ja"}
	if !reflect.DeepEqual(goEntry.VSCodeExtensions, want) {
		t.Fatalf("unexpected extensions: %#v", goEntry.VSCodeExtensions)
	}
}

// common に locale が設定されていて言語側に locale がない場合、common の locale が引き継がれる。
func TestLoadLanguageImages_MergesLocaleFromCommon(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{
  "common": {
    "locale": {
      "lang": "ja_JP.UTF-8",
      "language": "ja_JP:ja",
      "lcAll": "ja_JP.UTF-8"
    }
  },
  "go": {
    "image": "golang:1.24-alpine"
  }
}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	entries, err := loadLanguageImages(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	goEntry := entries["go"]
	want := entity.LocaleConfig{Lang: "ja_JP.UTF-8", Language: "ja_JP:ja", LcAll: "ja_JP.UTF-8"}
	if goEntry.Locale != want {
		t.Fatalf("expected locale from common, got: %#v", goEntry.Locale)
	}
}

// 言語側に locale が設定されている場合は common の locale より優先される。
func TestLoadLanguageImages_LanguageLocaleOverridesCommon(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{
  "common": {
    "locale": {
      "lang": "ja_JP.UTF-8",
      "language": "ja_JP:ja",
      "lcAll": "ja_JP.UTF-8"
    }
  },
  "go": {
    "image": "golang:1.24-alpine",
    "locale": {
      "lang": "en_US.UTF-8",
      "language": "en_US:en",
      "lcAll": "en_US.UTF-8"
    }
  }
}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	entries, err := loadLanguageImages(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	goEntry := entries["go"]
	want := entity.LocaleConfig{Lang: "en_US.UTF-8", Language: "en_US:en", LcAll: "en_US.UTF-8"}
	if goEntry.Locale != want {
		t.Fatalf("expected language locale to override common, got: %#v", goEntry.Locale)
	}
}

// 言語エントリを文字列で記述した場合（{"go": "golang:latest"}）もサポートする。
func TestLoadLanguageImages_StringEntryIsSupported(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{"go": "golang:1.24-alpine"}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	entries, err := loadLanguageImages(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	goEntry, ok := entries["go"]
	if !ok {
		t.Fatal("expected go entry")
	}
	if goEntry.Image != "golang:1.24-alpine" {
		t.Fatalf("unexpected image: %s", goEntry.Image)
	}
}

// "$schema" キーは言語エントリとして扱われない。
func TestLoadLanguageImages_SchemaKeyIsExcluded(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{
  "$schema": "https://example.com/schema.json",
  "go": {"image": "golang:1.24-alpine"}
}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	entries, err := loadLanguageImages(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := entries["$schema"]; ok {
		t.Fatal("$schema should not appear as a language entry")
	}
}

// "common" キーは言語エントリとして扱われない。
func TestLoadLanguageImages_CommonKeyIsExcluded(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "images.json")
	cfg := `{
  "common": {"timezone": "UTC"},
  "go": {"image": "golang:1.24-alpine"}
}`
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	entries, err := loadLanguageImages(cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := entries["common"]; ok {
		t.Fatal("common should not appear as a language entry")
	}
}

// config ファイルが存在しない場合はエラーを返す。
// (fetchBaseImageConfig は nil を返し、json.Unmarshal(nil) がエラーになる)
func TestLoadLanguageImages_FileNotFoundReturnsError(t *testing.T) {
	_, err := loadLanguageImages(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err == nil {
		t.Fatal("expected error when config file does not exist")
	}
}

// ─── fetchBaseImageConfig ──────────────────────────────────────────────────────

// http:// スキームは明示的にエラーで拒否される（https:// のみ許可）。
func TestFetchBaseImageConfig_HttpUrlReturnsError(t *testing.T) {
	_, err := fetchBaseImageConfig("http://example.com/config.json")
	if err == nil {
		t.Fatal("expected error for http:// URL")
	}
	if !strings.Contains(err.Error(), "http://") {
		t.Fatalf("expected error to mention http://, got: %v", err)
	}
}

// ファイルが存在しない場合は nil を返し、エラーにはならない。
func TestFetchBaseImageConfig_FileNotFoundReturnsNil(t *testing.T) {
	raw, err := fetchBaseImageConfig(filepath.Join(t.TempDir(), "nonexistent.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if raw != nil {
		t.Fatalf("expected nil for missing file, got: %s", raw)
	}
}

// ─── parseLanguageEntry ────────────────────────────────────────────────────────

// 文字列形式 ("golang:latest") は Image フィールドにセットされる。
func TestParseLanguageEntry_StringFormat(t *testing.T) {
	raw := json.RawMessage(`"golang:1.24-alpine"`)
	entry, err := parseLanguageEntry(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Image != "golang:1.24-alpine" {
		t.Fatalf("unexpected image: %s", entry.Image)
	}
}

// オブジェクト形式ですべてのフィールドが正しくパースされる。
func TestParseLanguageEntry_ObjectFormat(t *testing.T) {
	raw := json.RawMessage(`{
		"image": "ubuntu:24.04",
		"install": "apt-get install -y curl",
		"timezone": "Asia/Tokyo",
		"vscodeExtensions": ["ms-python.python"]
	}`)
	entry, err := parseLanguageEntry(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Image != "ubuntu:24.04" {
		t.Fatalf("unexpected image: %s", entry.Image)
	}
	if entry.Install != "apt-get install -y curl" {
		t.Fatalf("unexpected install: %s", entry.Install)
	}
	if entry.Timezone != "Asia/Tokyo" {
		t.Fatalf("unexpected timezone: %s", entry.Timezone)
	}
	if !reflect.DeepEqual(entry.VSCodeExtensions, []string{"ms-python.python"}) {
		t.Fatalf("unexpected extensions: %#v", entry.VSCodeExtensions)
	}
}

// locale フィールドが正しくパースされ entity.LocaleConfig に変換される。
func TestParseLanguageEntry_LocaleIsParsed(t *testing.T) {
	raw := json.RawMessage(`{
		"image": "ubuntu:24.04",
		"locale": {
			"lang": "en_US.UTF-8",
			"language": "en_US:en",
			"lcAll": "en_US.UTF-8"
		}
	}`)
	entry, err := parseLanguageEntry(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := entity.LocaleConfig{Lang: "en_US.UTF-8", Language: "en_US:en", LcAll: "en_US.UTF-8"}
	if entry.Locale != want {
		t.Fatalf("unexpected locale: %#v", entry.Locale)
	}
}

// locale を省略した場合は空の LocaleConfig になる。
func TestParseLanguageEntry_LocaleIsEmptyWhenOmitted(t *testing.T) {
	raw := json.RawMessage(`{"image": "golang:1.24-alpine"}`)
	entry, err := parseLanguageEntry(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Locale != (entity.LocaleConfig{}) {
		t.Fatalf("expected empty locale, got: %#v", entry.Locale)
	}
}

// image なしで install を指定するとエラーを返す。
func TestParseLanguageEntry_ImageMissingWithInstallReturnsError(t *testing.T) {
	raw := json.RawMessage(`{"install": "apt-get install -y curl"}`)
	_, err := parseLanguageEntry(raw)
	if err == nil {
		t.Fatal("expected error when image is missing but install is present")
	}
}

// 各フィールドの前後の空白は除去される。
func TestParseLanguageEntry_WhitespaceIsTrimmed(t *testing.T) {
	raw := json.RawMessage(`{"image": "  ubuntu:24.04  ", "timezone": "  Asia/Tokyo  "}`)
	entry, err := parseLanguageEntry(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Image != "ubuntu:24.04" {
		t.Fatalf("expected trimmed image, got: %q", entry.Image)
	}
	if entry.Timezone != "Asia/Tokyo" {
		t.Fatalf("expected trimmed timezone, got: %q", entry.Timezone)
	}
}

// vscodeExtensions 内の空文字列・空白のみのエントリはスキップされる。
func TestParseLanguageEntry_EmptyExtensionsAreSkipped(t *testing.T) {
	raw := json.RawMessage(`{"image": "ubuntu:24.04", "vscodeExtensions": ["ms-python.python", "  ", "golang.Go"]}`)
	entry, err := parseLanguageEntry(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"ms-python.python", "golang.Go"}
	if !reflect.DeepEqual(entry.VSCodeExtensions, want) {
		t.Fatalf("expected empty extensions to be skipped, got: %#v", entry.VSCodeExtensions)
	}
}
