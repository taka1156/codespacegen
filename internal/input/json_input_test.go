package input

import (
	"errors"
	"testing"

	"github.com/taka1156/codespacegen/internal/domain/entity"

	"github.com/google/go-cmp/cmp"
)

type fakeLoader struct {
	data []byte
	err  error
}

func (f *fakeLoader) Load(_ string) ([]byte, error) {
	return f.data, f.err
}

func newJsonInputWithFakes(httpsLoader, fileLoader baseImageConfigLoader) *JsonInput {
	return &JsonInput{
		httpsLoader: httpsLoader,
		fileLoader:  fileLoader,
	}
}

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
	raw := []byte(`{
		"langs": [{"profileName": "python", "image": "python:3.12"},{"profileName": "node", "image":"node:20"}]
	}`)
	expected := &entity.JsonConfig{
		Langs: []*entity.LangEntry{
			{ProfileName: "python", Image: "python:3.12"},
			{ProfileName: "node", Image: "node:20"},
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
	found := false
	for _, entry := range got.Langs {
		if entry.ProfileName == "python" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected python key in result")
	}
	if diff := cmp.Diff(got, expected); diff != "" {
		t.Errorf("mismatch (-got +expected):\n%s", diff)
	}
}

func TestLoadLanguageImages_ParsesValidJSONFromHTTPSLoader(t *testing.T) {
	raw := []byte(`{ "langs": [{"profileName": "rust", "image": "rust:1.76"}]}`)
	expected := &entity.JsonConfig{
		Langs: []*entity.LangEntry{
			{ProfileName: "rust", Image: "rust:1.76"},
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
	found := false
	for _, entry := range got.Langs {
		if entry.ProfileName == "rust" {
			found = true
			break
		}
	}
	if !found {
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
