package app

import (
	"encoding/json"

	"codespacegen/internal/domain/entity"
	"codespacegen/internal/workflow/collect"
)

type inputCollector interface {
	CollectConfig() (*collect.CollectedInputs, error)
}

type configAssembler interface {
	Resolve(clientConfig entity.ClientConfig, defaultSetting entity.DefaultSetting, jsonConfig map[string]json.RawMessage) (*entity.CodespaceConfig, error)
}

type artifactExecutor interface {
	Execute(config entity.CodespaceConfig, enableOverwriteFile bool, outputDir string) error
}
