package resolve

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

func (cscr *CodespaceConfigResolver) MergeLanguageEntries(commonEntry *entity.CommonEntry, langEntries map[string]*entity.LangEntry) (map[string]entity.LangEntry, error) {
	mergedImages := make(map[string]entity.LangEntry)

	switch {
	case commonEntry == nil && langEntries == nil:
		return mergedImages, nil
	case langEntries == nil:
		return mergedImages, nil
	case commonEntry == nil:
		for k, entry := range langEntries {
			normalizedKey := strings.ToLower(strings.TrimSpace(k))
			if normalizedKey == "" || normalizedKey == "common" || normalizedKey == "$schema" {
				continue
			}
			mergedImages[normalizedKey] = *entry
		}
		return mergedImages, nil
	default:
		for k, entry := range langEntries {
			normalizedKey := strings.ToLower(strings.TrimSpace(k))
			if normalizedKey == "" || normalizedKey == "common" || normalizedKey == "$schema" {
				continue
			}
			mergedImages[normalizedKey] = mergeLanguageEntries(*commonEntry, *entry)
		}
	}

	return mergedImages, nil

}

func mergeLanguageEntries(common entity.CommonEntry, LangEntry entity.LangEntry) entity.LangEntry {
	var baseLocale, overrideLocale entity.LocaleConfig
	if common.Locale != nil {
		baseLocale = *common.Locale
	} else {
		baseLocale = entity.DefaultLocale
	}
	if LangEntry.Locale != nil {
		overrideLocale = *LangEntry.Locale
	} else {
		overrideLocale = entity.DefaultLocale
	}
	merged := entity.LangEntry{
		Image:      LangEntry.Image,
		RunCommand: LangEntry.RunCommand,
		Timezone:   utils.Ptr(firstNonEmpty(LangEntry.Timezone, common.Timezone)),
		Locale:     utils.Ptr(mergeLocale(baseLocale, overrideLocale)),
	}

	switch {
	case LangEntry.VSCodeExtensions != nil && common.VSCodeExtensions != nil:
		merged.VSCodeExtensions = utils.Ptr(append(*LangEntry.VSCodeExtensions, *common.VSCodeExtensions...))
	case LangEntry.VSCodeExtensions != nil:
		merged.VSCodeExtensions = LangEntry.VSCodeExtensions
	case common.VSCodeExtensions != nil:
		merged.VSCodeExtensions = common.VSCodeExtensions
	}

	return merged
}

func mergeLocale(base entity.LocaleConfig, override entity.LocaleConfig) entity.LocaleConfig {
	if strings.TrimSpace(override.Lang) == "" {
		return base
	}

	return override
}
