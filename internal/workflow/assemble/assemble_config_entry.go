package assemble

import (
	"strings"

	"github.com/taka1156/codespacegen/internal/domain/entity"
	"github.com/taka1156/codespacegen/internal/utils"
)

func (acc *AssembleCodespaceConfig) resolveMergedEntry(jsonConfig entity.JsonConfig) (map[string]entity.LangEntry, error) {
	mergedImages := make(map[string]entity.LangEntry)

	switch {
	case jsonConfig.Common == nil && jsonConfig.Langs == nil:
		return mergedImages, nil
	case jsonConfig.Langs == nil:
		return mergedImages, nil
	case jsonConfig.Common == nil:
		for _, entry := range jsonConfig.Langs {
			normalizedKey := strings.ToLower(strings.TrimSpace(entry.ProfileName))
			if normalizedKey == "" || normalizedKey == "common" || normalizedKey == "$schema" {
				continue
			}
			mergedImages[normalizedKey] = *entry
		}
		return mergedImages, nil
	default:
		for _, entry := range jsonConfig.Langs {
			normalizedKey := strings.ToLower(strings.TrimSpace(entry.ProfileName))
			if normalizedKey == "" || normalizedKey == "common" || normalizedKey == "$schema" {
				continue
			}
			mergedImages[normalizedKey] = mergeLanguageEntries(*jsonConfig.Common, *entry)
		}
	}

	return mergedImages, nil

}

func mergeLanguageEntries(common entity.CommonEntry, langEntry entity.LangEntry) entity.LangEntry {
	merged := entity.LangEntry{
		ProfileName:   langEntry.ProfileName,
		Image:         langEntry.Image,
		LinuxPackages: langEntry.LinuxPackages,
		RunCommand:    langEntry.RunCommand,
	}

	switch {
	case common.VSCodeExtensions != nil && langEntry.VSCodeExtensions != nil:
		commonCopy := make([]string, len(*common.VSCodeExtensions))
		copy(commonCopy, *common.VSCodeExtensions)
		mergedExtensions := append(commonCopy, *langEntry.VSCodeExtensions...)
		merged.VSCodeExtensions = utils.Ptr(uniqueStringsPreserveOrder(mergedExtensions))
	case langEntry.VSCodeExtensions != nil:
		merged.VSCodeExtensions = langEntry.VSCodeExtensions
	case common.VSCodeExtensions != nil:
		merged.VSCodeExtensions = common.VSCodeExtensions
	}

	return merged
}
