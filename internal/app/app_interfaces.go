package app

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
	"codespacegen/internal/workflow/collect"
)

type inputCollector interface {
	CollectConfig(args []string) (*collect.CollectedInputs, error)
}

type initializeSettingJson interface {
	Execute(templateJson entity.TemplateJson, settingJsonFileName string) error
}

type assembleConfigResolver interface {
	Resolve(clientConfig entity.ClientConfig, defaultSetting entity.DefaultSetting, jsonConfig map[string]json.RawMessage) (*entity.CodespaceConfig, error)
}

type generateCodespaceArtifacts interface {
	Execute(config entity.CodespaceConfig, enableOverwriteFile bool, outputDir string) error
}
