package assemble

import (
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/utils"
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

func (acc *AssembleCodespaceConfig) resolveMergedEntry(jsonConfig entity.JsonConfig) (map[string]entity.LangEntry, error) {
	mergedImages := make(map[string]entity.LangEntry)

	switch {
	case jsonConfig.Common == nil && jsonConfig.Langs == nil:
		return mergedImages, nil
	case jsonConfig.Langs == nil:
		return mergedImages, nil
	case jsonConfig.Common == nil:
		for k, entry := range jsonConfig.Langs {
			normalizedKey := strings.ToLower(strings.TrimSpace(k))
			if normalizedKey == "" || normalizedKey == "common" || normalizedKey == "$schema" {
				continue
			}
			mergedImages[normalizedKey] = *entry
		}
		return mergedImages, nil
	default:
		for k, entry := range jsonConfig.Langs {
			normalizedKey := strings.ToLower(strings.TrimSpace(k))
			if normalizedKey == "" || normalizedKey == "common" || normalizedKey == "$schema" {
				continue
			}
			mergedImages[normalizedKey] = mergeLanguageEntries(*jsonConfig.Common, *entry)
		}
	}

	return mergedImages, nil

}

func mergeLanguageEntries(common entity.CommonEntry, LangEntry entity.LangEntry) entity.LangEntry {
	// priority: language-specific > common
	var baseLocale entity.LocaleConfig = entity.DefaultLocale

	if common.Locale != nil {
		baseLocale = *common.Locale
	}

	var resolvedLocale entity.LocaleConfig
	if LangEntry.Locale != nil {
		resolvedLocale = resolveLocale(baseLocale, *LangEntry.Locale)
	} else {
		resolvedLocale = baseLocale
	}

	merged := entity.LangEntry{
		Image:         LangEntry.Image,
		LinuxPackages: LangEntry.LinuxPackages,
		RunCommand:    LangEntry.RunCommand,
		Timezone:      utils.Ptr(firstNonEmpty(LangEntry.Timezone, common.Timezone)),
		Locale:        utils.Ptr(resolvedLocale),
	}

	switch {
	case common.VSCodeExtensions != nil && LangEntry.VSCodeExtensions != nil:
		commonCopy := make([]string, len(*common.VSCodeExtensions))
		copy(commonCopy, *common.VSCodeExtensions)
		merged.VSCodeExtensions = utils.Ptr(append(commonCopy, *LangEntry.VSCodeExtensions...))
	case LangEntry.VSCodeExtensions != nil:
		merged.VSCodeExtensions = LangEntry.VSCodeExtensions
	case common.VSCodeExtensions != nil:
		merged.VSCodeExtensions = common.VSCodeExtensions
	}

	return merged
}

func resolveLocale(base entity.LocaleConfig, override entity.LocaleConfig) entity.LocaleConfig {
	if strings.TrimSpace(override.Lang) == "" {
		return base
	}

	return override
}

func firstNonEmpty(values ...*string) string {
	for _, v := range values {
		if v != nil {
			trimmed := strings.TrimSpace(*v)
			if trimmed != "" {
				return trimmed
			}
		}
	}

	return ""
}
