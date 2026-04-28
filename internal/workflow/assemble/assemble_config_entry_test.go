package assemble

import (
	"testing"

	"codespacegen/internal/domain/entity"
	"codespacegen/internal/utils"
)

// --- resolveMergedEntry ---

func TestResolveMergedEntry_ReturnsEmptyWhenLangsNil(t *testing.T) {
	acc := NewAssembleCodespaceConfig(nil)
	got, err := acc.resolveMergedEntry(entity.JsonConfig{Common: &entity.CommonEntry{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("len: got %d, want 0", len(got))
	}
}

func TestResolveMergedEntry_ReturnsEmptyWhenBothNil(t *testing.T) {
	acc := NewAssembleCodespaceConfig(nil)
	got, err := acc.resolveMergedEntry(entity.JsonConfig{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("len: got %d, want 0", len(got))
	}
}

func TestResolveMergedEntry_CopiesLangsWhenCommonNil(t *testing.T) {
	acc := NewAssembleCodespaceConfig(nil)
	jsonConfig := entity.JsonConfig{
		Langs: map[string]*entity.LangEntry{
			"python": {Image: "python:3.12"},
		},
	}
	got, err := acc.resolveMergedEntry(jsonConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["python"]; !ok {
		t.Error("expected 'python' key in result")
	}
}

func TestResolveMergedEntry_SkipsReservedKeys(t *testing.T) {
	acc := NewAssembleCodespaceConfig(nil)
	jsonConfig := entity.JsonConfig{
		Common: &entity.CommonEntry{},
		Langs: map[string]*entity.LangEntry{
			"python":  {Image: "python:3.12"},
			"common":  {Image: "should-be-skipped"},
			"$schema": {Image: "should-be-skipped"},
			"":        {Image: "should-be-skipped"},
		},
	}
	got, err := acc.resolveMergedEntry(jsonConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("len: got %d, want 1", len(got))
	}
	if _, ok := got["python"]; !ok {
		t.Error("expected 'python' key in result")
	}
}

func TestResolveMergedEntry_NormalizesKeyToLower(t *testing.T) {
	acc := NewAssembleCodespaceConfig(nil)
	jsonConfig := entity.JsonConfig{
		Common: &entity.CommonEntry{},
		Langs: map[string]*entity.LangEntry{
			"RUST": {Image: "rust:latest"},
		},
	}
	got, err := acc.resolveMergedEntry(jsonConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["rust"]; !ok {
		t.Error("expected 'rust' key (lowercase) in result")
	}
}

// --- mergeLanguageEntries ---

func TestMergeLanguageEntries_LangExtensionsAppendCommon(t *testing.T) {
	common := entity.CommonEntry{
		VSCodeExtensions: utils.Ptr([]string{"common.ext"}),
	}
	lang := entity.LangEntry{
		Image:            "python:3.12",
		VSCodeExtensions: utils.Ptr([]string{"lang.ext"}),
	}
	got := mergeLanguageEntries(common, lang)
	if len(*got.VSCodeExtensions) != 2 {
		t.Errorf("VSCodeExtensions len: got %d, want 2", len(*got.VSCodeExtensions))
	}
}

func TestMergeLanguageEntries_OnlyCommonExtensionsWhenLangHasNone(t *testing.T) {
	common := entity.CommonEntry{
		VSCodeExtensions: utils.Ptr([]string{"common.ext"}),
	}
	lang := entity.LangEntry{Image: "python:3.12"}
	got := mergeLanguageEntries(common, lang)
	if got.VSCodeExtensions == nil || len(*got.VSCodeExtensions) != 1 {
		t.Errorf("VSCodeExtensions: got %v, want [common.ext]", got.VSCodeExtensions)
	}
}

func TestMergeLanguageEntries_OnlyLangExtensionsWhenCommonHasNone(t *testing.T) {
	common := entity.CommonEntry{}
	lang := entity.LangEntry{
		Image:            "python:3.12",
		VSCodeExtensions: utils.Ptr([]string{"lang.ext"}),
	}
	got := mergeLanguageEntries(common, lang)
	if got.VSCodeExtensions == nil || len(*got.VSCodeExtensions) != 1 {
		t.Errorf("VSCodeExtensions: got %v, want [lang.ext]", got.VSCodeExtensions)
	}
}

func TestMergeLanguageEntries_LangTimezoneTakesPriorityOverCommon(t *testing.T) {
	common := entity.CommonEntry{Timezone: utils.Ptr("UTC")}
	lang := entity.LangEntry{
		Image:    "python:3.12",
		Timezone: utils.Ptr("Asia/Tokyo"),
	}
	got := mergeLanguageEntries(common, lang)
	if got.Timezone == nil || *got.Timezone != "Asia/Tokyo" {
		t.Errorf("Timezone: got %v, want Asia/Tokyo", got.Timezone)
	}
}

func TestMergeLanguageEntries_CommonTimezoneUsedWhenLangTimezoneEmpty(t *testing.T) {
	common := entity.CommonEntry{Timezone: utils.Ptr("UTC")}
	lang := entity.LangEntry{Image: "python:3.12"}
	got := mergeLanguageEntries(common, lang)
	if got.Timezone == nil || *got.Timezone != "UTC" {
		t.Errorf("Timezone: got %v, want UTC", got.Timezone)
	}
}

func TestMergeLanguageEntries_LangLocaleOverridesCommon(t *testing.T) {
	common := entity.CommonEntry{
		Locale: &entity.LocaleConfig{Lang: "en_US.UTF-8", Language: "en_US:en", LcAll: "en_US.UTF-8"},
	}
	lang := entity.LangEntry{
		Image:  "python:3.12",
		Locale: &entity.LocaleConfig{Lang: "ja_JP.UTF-8", Language: "ja_JP:ja", LcAll: "ja_JP.UTF-8"},
	}
	got := mergeLanguageEntries(common, lang)
	if got.Locale == nil || got.Locale.Lang != "ja_JP.UTF-8" {
		t.Errorf("Locale.Lang: got %v, want ja_JP.UTF-8", got.Locale)
	}
}

func TestMergeLanguageEntries_CommonLocaleUsedWhenLangLocaleNil(t *testing.T) {
	commonLocale := entity.LocaleConfig{Lang: "en_US.UTF-8", Language: "en_US:en", LcAll: "en_US.UTF-8"}
	common := entity.CommonEntry{
		Locale: &commonLocale,
	}
	lang := entity.LangEntry{Image: "python:3.12"}
	got := mergeLanguageEntries(common, lang)
	// LangEntry.Locale が nil の場合、common.Locale がそのまま使われる
	if got.Locale == nil || got.Locale.Lang != commonLocale.Lang {
		t.Errorf("Locale.Lang: got %v, want %q", got.Locale, commonLocale.Lang)
	}
}
