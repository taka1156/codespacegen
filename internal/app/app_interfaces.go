package app

import (
	"github.com/taka1156/codespacegen/internal/domain/entity"
	"github.com/taka1156/codespacegen/internal/workflow/collect"
)

type inputCollector interface {
	CollectConfig(args []string) (*collect.CollectedInputs, error)
}

type initializeSettingJson interface {
	Execute(templateJson entity.JsonConfig, settingJsonFileName string) error
}

type assembleConfigResolver interface {
	Resolve(clientConfig entity.ClientConfig, defaultSetting entity.DefaultSetting, jsonConfig entity.JsonConfig) (*entity.CodespaceConfig, error)
}

type generateCodespaceArtifacts interface {
	Execute(config entity.CodespaceConfig, enableOverwriteFile bool, outputDir string) error
}

type updateCommandline interface {
	Update(currentVersion string) error
}
