package config

import "codespacegen/internal/domain/entity"

type DefaultSetting struct {
	Timezone            string
	Image               string
	Version             string
	PortMappingPatterns entity.PortMappingPatterns
}

type DefaultConfig struct {
}

func NewDefaultConfig() *DefaultConfig {
	return &DefaultConfig{}
}

func (dc *DefaultConfig) GetDefaultSetting() DefaultSetting {
	ds := DefaultSetting{}
	ds.Image = entity.DefaultImage
	ds.Timezone = entity.DefaultTimezone
	ds.Version = entity.DefaultVersion
	ds.PortMappingPatterns = entity.DefaultPortMappingPatterns

	return ds
}
