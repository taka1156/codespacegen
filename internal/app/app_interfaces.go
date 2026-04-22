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
	Resolve(cliConfig entity.CliConfig, defaultSetting entity.DefaultSetting, overrides map[string]json.RawMessage, defaultTimezone string, defaultImage string) (*entity.CodespaceConfig, error)
}

type artifactExecutor interface {
	Execute(config entity.CodespaceConfig, enableOverwriteFile bool, outputDir string) error
}
