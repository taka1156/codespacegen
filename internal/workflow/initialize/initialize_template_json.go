package initialize

import (
	"codespacegen/internal/domain/entity"
	"codespacegen/internal/domain/service"
	"fmt"
)

type InitializeSettingJson struct {
	settingTemplateGenerator service.SettingTemplateGenerator
	workdirProvider          service.WorkdirProvider
	writer                   service.LocalFileWriter
}

func NewInitializeSettingJson(settingTemplateGenerator service.SettingTemplateGenerator, workdirProvider service.WorkdirProvider, writer service.LocalFileWriter) *InitializeSettingJson {
	return &InitializeSettingJson{
		settingTemplateGenerator: settingTemplateGenerator,
		workdirProvider:          workdirProvider,
		writer:                   writer,
	}
}

func (isj *InitializeSettingJson) Execute(templateJson entity.TemplateJson, settingJsonFileName string) error {
	file, err := isj.settingTemplateGenerator.Generate(templateJson)
	if err != nil {
		return fmt.Errorf("failed to generate template JSON: %w", err)
	}

	outputPath, err := isj.workdirProvider.GetConfigOutputPath()
	if err != nil {
		return fmt.Errorf("failed to get config output path: %w", err)
	}

	err = isj.writer.Write(outputPath+settingJsonFileName, file, false)
	if err != nil {
		return fmt.Errorf("failed to write template JSON: %w", err)
	}

	return nil
}
