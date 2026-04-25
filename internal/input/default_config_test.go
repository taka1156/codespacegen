package input

import (
	"codespacegen/internal/domain/entity"
	"testing"
)

func TestDefaultConfig_GetDefaultSetting_ReturnsExpectedValues(t *testing.T) {
	dc := NewDefaultConfig()
	got := dc.GetDefaultSetting()

	if got.Image != entity.DefaultImage {
		t.Errorf("Image: got %q, want %q", got.Image, entity.DefaultImage)
	}
	if got.Timezone != entity.DefaultTimezone {
		t.Errorf("Timezone: got %q, want %q", got.Timezone, entity.DefaultTimezone)
	}
	if got.VscSchema != entity.DefaultVscSchema {
		t.Errorf("VscSchema: got %q, want %q", got.VscSchema, entity.DefaultVscSchema)
	}
	if got.SettingJsonFileName != entity.DefaultTemplateJsonPath {
		t.Errorf("SettingJsonFileName: got %q, want %q", got.SettingJsonFileName, entity.DefaultTemplateJsonPath)
	}
}

func TestDefaultConfig_GetDefaultSetting_AlpineModulesContainTzdata(t *testing.T) {
	dc := NewDefaultConfig()
	got := dc.GetDefaultSetting()

	found := false
	for _, m := range got.OsModules.AlpineModules {
		if m == "tzdata" {
			found = true
			break
		}
	}
	if !found {
		t.Error("AlpineModules should contain 'tzdata'")
	}
}

func TestDefaultConfig_GetDefaultSetting_DebianLikeModulesContainLocales(t *testing.T) {
	dc := NewDefaultConfig()
	got := dc.GetDefaultSetting()

	found := false
	for _, m := range got.OsModules.DebianLikeModules {
		if m == "locales" {
			found = true
			break
		}
	}
	if !found {
		t.Error("DebianLikeModules should contain 'locales'")
	}
}
