package setting

import (
	"codespacegen/internal/domain/entity"
	"encoding/json"
	"testing"
)

func TestSettingTemplateGenerator_Generate_ReturnsValidJSON(t *testing.T) {
	g := NewSettingTemplateGenerator()
	input := entity.TemplateJson{
		Schema: "https://example.com/schema.json",
		Go: entity.JsonEntry{
			Image: "golang:1.24-alpine",
		},
	}

	got, err := g.Generate(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantBytes, _ := json.MarshalIndent(input, "", "  ")
	if got != string(wantBytes) {
		t.Errorf("got:\n%s\nwant:\n%s", got, string(wantBytes))
	}
}

func TestSettingTemplateGenerator_Generate_EmptyTemplateReturnsJSON(t *testing.T) {
	g := NewSettingTemplateGenerator()
	got, err := g.Generate(entity.TemplateJson{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == "" {
		t.Error("expected non-empty JSON output, got empty string")
	}
}
