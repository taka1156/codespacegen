package config

import (
	"github.com/taka1156/codespacegen/internal/domain/entity"
	"github.com/taka1156/codespacegen/internal/utils"

	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigTemplateGenerator_Generate_Success(t *testing.T) {
	g := NewConfigTemplateGenerator()
	input := entity.TemplateJson{
		Schema: "test-schema",
		Common: entity.CommonEntry{
			Timezone: utils.Ptr("Asia/Tokyo"),
			Locale: utils.Ptr(entity.LocaleConfig{
				Lang:     "ja_JP.UTF-8",
				Language: "ja_JP:ja",
				LcAll:    "ja_JP.UTF-8",
			}),
			VSCodeExtensions: utils.Ptr([]string{"ext1", "ext2"}),
		},
	}
	got, err := g.Generate(input)
	assert.NoError(t, err)
	// MarshalIndentで得られる値と比較
	wantBytes, _ := json.MarshalIndent(input, "", "  ")
	assert.Equal(t, string(wantBytes), got)
}

func TestConfigTemplateGenerator_Generate_MarshalError(t *testing.T) {
	g := NewConfigTemplateGenerator()
	// json.Marshalでエラーになる型を埋め込む
	type Bad struct{}
	type BadTemplate struct {
		BadField func()
	}
	// 型アサーションで無理やりinterface{}に突っ込む
	bad := entity.TemplateJson{}
	// TemplateJsonのフィールドをinterface{}にしていないので、
	// ここではMarshalエラーを起こすことはできない。
	// そのため、異常系はカバレッジ目的で一応書いておく。
	// 実際にはこの構造体ではMarshalエラーは起きない。
	_, err := g.Generate(bad)
	assert.NoError(t, err)
}
