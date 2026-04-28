package resolve

import (
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/i18n"
	"codespacegen/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type parsedLanguageSetting struct {
	Image      string `json:"image"`
	RunCommand string `json:"runCommand"`
	Timezone   string `json:"timezone"`
	Locale     struct {
		Lang     string `json:"lang"`
		Language string `json:"language"`
		LcAll    string `json:"lcAll"`
	} `json:"locale"`
	VSCodeExtensions []string `json:"vscodeExtensions"`
}

func (cscr *CodespaceConfigResolver) MergeLanguageEntries(overrides map[string]json.RawMessage) (map[string]entity.LangEntry, error) {
	mergedImages := make(map[string]entity.LangEntry)

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

func mergeLanguageEntries(base entity.LangEntry, override entity.LangEntry) entity.LangEntry {
	var baseLocale, overrideLocale entity.LocaleConfig
	if base.Locale != nil {
		baseLocale = *base.Locale
	} else {
		baseLocale = entity.DefaultLocale
	}
	if override.Locale != nil {
		overrideLocale = *override.Locale
	} else {
		overrideLocale = entity.DefaultLocale
	}
	merged := entity.LangEntry{
		Image:      firstNonEmpty(override.Image, base.Image),
		RunCommand: firstNonEmpty(override.RunCommand, base.RunCommand),
		Timezone:   firstNonEmpty(override.Timezone, base.Timezone),
		Locale:     utils.Ptr(mergeLocale(baseLocale, overrideLocale)),
	}

	merged.VSCodeExtensions = append(merged.VSCodeExtensions, base.VSCodeExtensions...)
	merged.VSCodeExtensions = append(merged.VSCodeExtensions, override.VSCodeExtensions...)

	return merged
}

func mergeLocale(base entity.LocaleConfig, override entity.LocaleConfig) entity.LocaleConfig {
	if strings.TrimSpace(override.Lang) == "" {
		return base
	}

	return override
}

func parseCommonEntry(overrides map[string]json.RawMessage) (entity.LangEntry, error) {
	for k, v := range overrides {
		if strings.ToLower(strings.TrimSpace(k)) != "common" {
			continue
		}

		entry, err := parseLanguageEntry(v)
		if err != nil {
			return entity.LangEntry{}, fmt.Errorf("%s: %w", i18n.T("error_invalid_entry_for_key", map[string]interface{}{"Key": k}), err)
		}
		return entry, nil
	}

	return entity.LangEntry{}, nil
}

func parseLanguageEntry(raw json.RawMessage) (entity.LangEntry, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return entity.LangEntry{Image: strings.TrimSpace(s)}, nil
	}

	var setting parsedLanguageSetting
	if err := json.Unmarshal(raw, &setting); err != nil {
		return entity.LangEntry{}, fmt.Errorf("%s: %w", i18n.T("error_must_be_string_or_object"), err)
	}

	entry := toJsonEntry(setting)
	if entry.Image == "" && entry.RunCommand != "" {
		return entity.LangEntry{}, errors.New(i18n.T("error_image_required_when_RunCommand"))
	}

	return entry, nil
}

func toJsonEntry(setting parsedLanguageSetting) entity.LangEntry {
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

	return entity.LangEntry{
		Image:            strings.TrimSpace(setting.Image),
		RunCommand:       strings.TrimSpace(setting.RunCommand),
		Locale:           utils.Ptr(locale),
		Timezone:         strings.TrimSpace(setting.Timezone),
		VSCodeExtensions: vscodeExtensions,
	}
}
