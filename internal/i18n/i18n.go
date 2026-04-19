package i18n

import (
	"embed"
	"os"
	"strings"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localeFS embed.FS

var (
	bundle    *goi18n.Bundle
	localizer *goi18n.Localizer
)

func init() {
	bundle = goi18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	for _, name := range []string{"locales/en.yaml", "locales/ja.yaml"} {
		_, _ = bundle.LoadMessageFileFS(localeFS, name)
	}

	SetLang(detectLanguage())
}

// SetLang switches the active locale at runtime (e.g. "en" or "ja").
func SetLang(lang string) {
	localizer = goi18n.NewLocalizer(bundle, lang, "en")
}

func detectLanguage() string {
	for _, env := range []string{"LANGUAGE", "LC_ALL", "LC_MESSAGES", "LANG"} {
		val := os.Getenv(env)
		if val == "" {
			continue
		}
		val = strings.Split(val, ".")[0]
		val = strings.Split(val, "_")[0]
		if val != "" && val != "C" && val != "POSIX" {
			return val
		}
	}
	return "en"
}

// T translates a message by ID with optional template data.
func T(messageID string, templateData ...map[string]interface{}) string {
	cfg := &goi18n.LocalizeConfig{MessageID: messageID}
	if len(templateData) > 0 {
		cfg.TemplateData = templateData[0]
	}
	msg, err := localizer.Localize(cfg)
	if err != nil {
		return messageID
	}
	return msg
}
