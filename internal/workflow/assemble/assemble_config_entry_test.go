package assemble

import (
	"testing"

	"github.com/taka1156/codespacegen/internal/domain/entity"
	"github.com/taka1156/codespacegen/internal/utils"
)

func TestResolveMergedEntry_ReturnsEmptyWhenLangsNil(t *testing.T) {
	acc := NewAssembleCodespaceConfig(nil)
	got := acc.resolveMergedEntry(entity.JsonConfig{Common: &entity.CommonEntry{}})
	if len(got) != 0 {
		t.Errorf("len: got %d, want 0", len(got))
	}
}

func TestResolveMergedEntry_ReturnsEmptyWhenBothNil(t *testing.T) {
	acc := NewAssembleCodespaceConfig(nil)
	got := acc.resolveMergedEntry(entity.JsonConfig{})
	if len(got) != 0 {
		t.Errorf("len: got %d, want 0", len(got))
	}
}

func TestResolveMergedEntry_CopiesLangsWhenCommonNil(t *testing.T) {
	acc := NewAssembleCodespaceConfig(nil)
	jsonConfig := entity.JsonConfig{
		Langs: []*entity.LangEntry{
			{ProfileName: "python", Image: "python:3.12"},
		},
	}
	got := acc.resolveMergedEntry(jsonConfig)
	found := false
	for _, entry := range got {
		if entry.ProfileName == "python" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'python' key in result")
	}
}

func TestResolveMergedEntry_SkipsReservedKeys(t *testing.T) {
	acc := NewAssembleCodespaceConfig(nil)
	jsonConfig := entity.JsonConfig{
		Common: &entity.CommonEntry{},
		Langs: []*entity.LangEntry{
			{ProfileName: "python", Image: "python:3.12"},
			{ProfileName: "common", Image: "should-be-skipped"},
			{ProfileName: "$schema", Image: "should-be-skipped"},
			{ProfileName: "", Image: "should-be-skipped"},
		},
	}
	got := acc.resolveMergedEntry(jsonConfig)
	if len(got) != 1 {
		t.Errorf("len: got %d, want 1", len(got))
	}
	found := false
	for _, entry := range got {
		if entry.ProfileName == "python" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'python' key in result")
	}
}

func TestResolveMergedEntry_NormalizesKeyToLower(t *testing.T) {
	acc := NewAssembleCodespaceConfig(nil)
	jsonConfig := entity.JsonConfig{
		Common: &entity.CommonEntry{},
		Langs: []*entity.LangEntry{
			{ProfileName: "rust", Image: "rust:latest"},
		},
	}
	got := acc.resolveMergedEntry(jsonConfig)
	found := false
	for _, entry := range got {
		if entry.ProfileName == "rust" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'rust' key (lowercase) in result")
	}
}

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

func TestMergeLanguageEntries_DeduplicatesExtensions(t *testing.T) {
	common := entity.CommonEntry{
		VSCodeExtensions: utils.Ptr([]string{"common.ext", "shared.ext"}),
	}
	lang := entity.LangEntry{
		Image:            "python:3.12",
		VSCodeExtensions: utils.Ptr([]string{"shared.ext", "lang.ext"}),
	}
	got := mergeLanguageEntries(common, lang)
	if len(*got.VSCodeExtensions) != 3 {
		t.Errorf("VSCodeExtensions len: got %d, want 3 (common.ext, shared.ext, lang.ext)", len(*got.VSCodeExtensions))
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
