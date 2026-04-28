package input

import (
	"codespacegen/internal/domain/entity"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// --- フェイク実装 ---

type fakeLoader struct {
	data []byte
	err  error
}

func (f *fakeLoader) Load(_ string) ([]byte, error) {
	return f.data, f.err
}

// newJsonInputWithFakes は httpsLoader と fileLoader にフェイクを注入した JsonInput を返す。
func newJsonInputWithFakes(httpsLoader, fileLoader baseImageConfigLoader) *JsonInput {
	return &JsonInput{
		httpsLoader: httpsLoader,
		fileLoader:  fileLoader,
	}
}

// --- LoadLanguageImages ---

func TestLoadLanguageImages_ReturnsNilWhenSourceIsEmpty(t *testing.T) {
	ji := newJsonInputWithFakes(&fakeLoader{}, &fakeLoader{})
	got, err := ji.LoadLanguageImages("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestLoadLanguageImages_ParsesValidJSONFromFileLoader(t *testing.T) {
	raw := []byte(`{"python": {"image": "python:3.12"},"node":{"image":"node:20"}}`)
	expected := &entity.JsonConfig{
		Langs: map[string]*entity.LangEntry{
			"python": {Image: "python:3.12"},
			"node":   {Image: "node:20"},
		},
	}
	ji := newJsonInputWithFakes(
		&fakeLoader{},
		&fakeLoader{data: raw},
	)

	got, err := ji.LoadLanguageImages("/some/path/config.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Langs) != 2 {
		t.Errorf("expected 2 entries, got %d", len(got.Langs))
	}
	if _, ok := got.Langs["python"]; !ok {
		t.Error("expected python key in result")
	}
	if diff := cmp.Diff(got, expected); diff != "" {
		t.Errorf("mismatch (-got +expected):\n%s", diff)
	}
}

func TestLoadLanguageImages_ParsesValidJSONFromHTTPSLoader(t *testing.T) {
	raw := []byte(`{"rust": {"image": "rust:1.76"}}`)
	expected := &entity.JsonConfig{
		Langs: map[string]*entity.LangEntry{
			"rust": {Image: "rust:1.76"},
		},
	}
	ji := newJsonInputWithFakes(
		&fakeLoader{data: raw},
		&fakeLoader{},
	)

	got, err := ji.LoadLanguageImages("https://example.com/config.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got.Langs["rust"]; !ok {
		t.Error("expected rust key in result")
	}
	if diff := cmp.Diff(got, expected); diff != "" {
		t.Errorf("mismatch (-got +expected):\n%s", diff)
	}
}

func TestLoadLanguageImages_ReturnsErrorForHTTPSource(t *testing.T) {
	ji := newJsonInputWithFakes(&fakeLoader{}, &fakeLoader{})
	_, err := ji.LoadLanguageImages("http://example.com/config.json")
	if err == nil {
		t.Fatal("expected error for http:// source, got nil")
	}
}

func TestLoadLanguageImages_ReturnsErrorWhenLoaderFails(t *testing.T) {
	ji := newJsonInputWithFakes(
		&fakeLoader{},
		&fakeLoader{err: errors.New("read error")},
	)
	_, err := ji.LoadLanguageImages("/bad/path.json")
	if err == nil {
		t.Fatal("expected error from loader, got nil")
	}
}

func TestLoadLanguageImages_ReturnsErrorOnInvalidJSON(t *testing.T) {
	ji := newJsonInputWithFakes(
		&fakeLoader{},
		&fakeLoader{data: []byte(`{invalid json}`)},
	)
	_, err := ji.LoadLanguageImages("/some/path.json")
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestLoadLanguageImages_ReturnsNilWhenFileLoaderReturnsNil(t *testing.T) {
	// fileConfigLoader.Load は ErrNotExist の場合 nil を返す
	ji := newJsonInputWithFakes(
		&fakeLoader{},
		&fakeLoader{data: nil, err: nil},
	)
	got, err := ji.LoadLanguageImages("/nonexistent.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}
