package i18n

import (
	"testing"
)

func TestT_ReturnsTranslatedMessage(t *testing.T) {
	SetLang("en")
	got := T("error_project_name_required")
	if got == "" || got == "error_project_name_required" {
		t.Errorf("T returned no translation: %q", got)
	}
}

func TestT_ReturnsMessageIDWhenKeyNotFound(t *testing.T) {
	SetLang("en")
	got := T("this_key_does_not_exist")
	if got != "this_key_does_not_exist" {
		t.Errorf("expected message ID fallback, got %q", got)
	}
}

func TestT_InterpolatesTemplateData(t *testing.T) {
	SetLang("en")
	got := T("error_unsupported_language", map[string]interface{}{"Language": "cobol"})
	want := "unsupported language: cobol"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSetLang_SwitchesToJapanese(t *testing.T) {
	SetLang("ja")
	defer SetLang("en")

	got := T("error_project_name_required")
	if got == "" || got == "error_project_name_required" {
		t.Errorf("T returned no Japanese translation: %q", got)
	}
}

func TestDetectLanguage_ReturnsEnByDefault(t *testing.T) {
	t.Setenv("LANGUAGE", "")
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "")
	got := detectLanguage()
	if got != "en" {
		t.Errorf("got %q, want %q", got, "en")
	}
}

func TestDetectLanguage_ReturnsLangFromEnvVar(t *testing.T) {
	t.Setenv("LANGUAGE", "")
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "ja_JP.UTF-8")
	got := detectLanguage()
	if got != "ja" {
		t.Errorf("got %q, want %q", got, "ja")
	}
}

func TestDetectLanguage_SkipsPosixAndC(t *testing.T) {
	t.Setenv("LANGUAGE", "C")
	t.Setenv("LC_ALL", "POSIX")
	t.Setenv("LC_MESSAGES", "")
	t.Setenv("LANG", "")
	got := detectLanguage()
	if got != "en" {
		t.Errorf("got %q, want %q", got, "en")
	}
}
