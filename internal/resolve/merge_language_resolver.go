package resolve

import (
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/i18n"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (cscr *CodeSpaceConfigResolver) MergeLanguageEntries(overrides map[string]json.RawMessage) (map[string]entity.JsonEntry, error) {
	mergedImages := make(map[string]entity.JsonEntry)

	common, err := parseCommonEntry(overrides)
	if err != nil {
		return nil, err
	}

	for k, v := range overrides {
		normalizedKey := strings.ToLower(strings.TrimSpace(k))
		if normalizedKey == "" || normalizedKey == "common" || normalizedKey == "$schema" {
			continue
		}
		entry, err := parseLanguageEntry(v)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", i18n.T("error_invalid_entry_for_key", map[string]interface{}{"Key": k}), err)
		}

		mergedImages[normalizedKey] = mergeLanguageEntries(common, entry)

	}

	return mergedImages, nil

}

func mergeLanguageEntries(base entity.JsonEntry, override entity.JsonEntry) entity.JsonEntry {
	locale := override.Locale
	if locale.Lang == "" {
		locale = base.Locale
	}

	merged := entity.JsonEntry{
		Image:    firstNonEmpty(override.Image, base.Image),
		Install:  firstNonEmpty(override.Install, base.Install),
		Timezone: firstNonEmpty(override.Timezone, base.Timezone),
		Locale:   locale,
	}

	merged.VSCodeExtensions = append(merged.VSCodeExtensions, base.VSCodeExtensions...)
	merged.VSCodeExtensions = append(merged.VSCodeExtensions, override.VSCodeExtensions...)

	return merged
}

func parseCommonEntry(overrides map[string]json.RawMessage) (entity.JsonEntry, error) {
	for k, v := range overrides {
		if strings.ToLower(strings.TrimSpace(k)) != "common" {
			continue
		}

		entry, err := parseLanguageEntry(v)
		if err != nil {
			return entity.JsonEntry{}, fmt.Errorf("%s: %w", i18n.T("error_invalid_entry_for_key", map[string]interface{}{"Key": k}), err)
		}
		return entry, nil
	}

	return entity.JsonEntry{}, nil
}

func parseLanguageEntry(raw json.RawMessage) (entity.JsonEntry, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return entity.JsonEntry{Image: strings.TrimSpace(s)}, nil
	}

	var setting struct {
		Image    string `json:"image"`
		Install  string `json:"install"`
		Timezone string `json:"timezone"`
		Locale   struct {
			Lang     string `json:"lang"`
			Language string `json:"language"`
			LcAll    string `json:"lcAll"`
		} `json:"locale"`
		VSCodeExtensions []string `json:"vscodeExtensions"`
	}

	if err := json.Unmarshal(raw, &setting); err != nil {
		return entity.JsonEntry{}, fmt.Errorf("%s: %w", i18n.T("error_must_be_string_or_object"), err)
	}

	image := strings.TrimSpace(setting.Image)
	install := strings.TrimSpace(setting.Install)
	timezone := strings.TrimSpace(setting.Timezone)
	locale := entity.LocaleConfig{
		Lang:     strings.TrimSpace(setting.Locale.Lang),
		Language: strings.TrimSpace(setting.Locale.Language),
		LcAll:    strings.TrimSpace(setting.Locale.LcAll),
	}
	vscodeExtensions := make([]string, 0, len(setting.VSCodeExtensions))
	for _, ext := range setting.VSCodeExtensions {
		trimmed := strings.TrimSpace(ext)
		if trimmed != "" {
			vscodeExtensions = append(vscodeExtensions, trimmed)
		}
	}
	if image == "" && install != "" {
		return entity.JsonEntry{}, errors.New(i18n.T("error_image_required_when_install"))
	}

	return entity.JsonEntry{Image: image, Install: install, Locale: locale, Timezone: timezone, VSCodeExtensions: vscodeExtensions}, nil
}
