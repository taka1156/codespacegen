package input

import "codespacegen/internal/domain/entity"

const DefaultImage = "alpine:latest"
const DefaultTimezone = "UTC"
const DefaultVersion = "dev"
const DefaultVscSchema = "https://raw.githubusercontent.com/microsoft/vscode/main/extensions/configuration-editing/schemas/devContainer.vscode.schema.json"

type DefaultConfig struct {
}

func NewDefaultConfig() *DefaultConfig {
	return &DefaultConfig{}
}

func (dc *DefaultConfig) GetDefaultSetting() entity.DefaultSetting {
	return entity.DefaultSetting{
		Image:     DefaultImage,
		Timezone:  DefaultTimezone,
		Version:   DefaultVersion,
		VscSchema: DefaultVscSchema,
	}
}
