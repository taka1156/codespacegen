package initialize

import (
	"fmt"

	"github.com/taka1156/codespacegen/internal/domain/entity"
	"github.com/taka1156/codespacegen/internal/domain/service"
	"github.com/taka1156/codespacegen/internal/utils"
)

type InitializeSettingJson struct {
	settingTemplateGenerator service.SettingTemplateGenerator
	writer                   service.LocalFileWriter
}

func NewInitializeSettingJson(
	settingTemplateGenerator service.SettingTemplateGenerator,
	writer service.LocalFileWriter,
) *InitializeSettingJson {
	return &InitializeSettingJson{
		settingTemplateGenerator: settingTemplateGenerator,
		writer:                   writer,
	}
}

func (isj *InitializeSettingJson) Execute(
	templateJson entity.JsonConfig,
	settingJsonFileName string,
	outputDir string,
) error {
	file, err := isj.settingTemplateGenerator.Generate(templateJson)
	if err != nil {
		return fmt.Errorf("failed to generate template JSON: %w", err)
	}

	outputPath, err := utils.ResolveOutputPath(outputDir, settingJsonFileName)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	err = isj.writer.Write(outputPath, file, false)
	if err != nil {
		return fmt.Errorf("failed to write template JSON: %w", err)
	}

	return nil
}
