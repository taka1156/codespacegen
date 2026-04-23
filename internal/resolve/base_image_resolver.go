package resolve

import (
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/i18n"
	"errors"
	"strings"
)

func (cscr *CodespaceConfigResolver) ResolveBaseImage(language string, explicitBaseImage string, imageConfig string, jsonEntries map[string]entity.JsonEntry, defaultImage string) (entity.JsonEntry, error) {
	if explicitBaseImage != "" {
		return entity.JsonEntry{Image: explicitBaseImage}, nil
	}

	if strings.TrimSpace(language) == "" {
		image := strings.TrimSpace(defaultImage)
		if image == "" {
			image = entity.DefaultImage
		}
		return entity.JsonEntry{Image: image}, nil
	}

	key := strings.ToLower(strings.TrimSpace(language))
	entry, ok := jsonEntries[key]
	if !ok {
		return entity.JsonEntry{}, errors.New(i18n.T("error_unsupported_language", map[string]interface{}{"Language": language}))
	}

	if entry.Image == "" {
		return entity.JsonEntry{}, errors.New(i18n.T("error_image_required_for_language", map[string]interface{}{"Language": language}))
	}

	return entry, nil
}
