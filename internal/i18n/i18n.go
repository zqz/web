package i18n

import (
	"encoding/json"
	"html/template"
	"strings"
)

// DefaultLocale is the fallback language.
const DefaultLocale = "en"

// Messages holds per-locale key -> message. Use En for English.
var Messages = map[string]map[string]string{
	"en": En,
}

// T returns the message for the given locale and key. Falls back to DefaultLocale, then to key.
func T(locale, key string) string {
	if locale == "" {
		locale = DefaultLocale
	}
	if msgs, ok := Messages[locale]; ok {
		if s, ok := msgs[key]; ok {
			return s
		}
	}
	if msgs, ok := Messages[DefaultLocale]; ok {
		if s, ok := msgs[key]; ok {
			return s
		}
	}
	return key
}

// TFunc returns a function suitable for template FuncMap: t(key) using the given locale.
func TFunc(locale string) func(string) string {
	return func(key string) string {
		return T(locale, key)
	}
}

// QuoteJS escapes s for safe use inside a JavaScript single-quoted string.
func QuoteJS(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '\'':
			b.WriteString(`\'`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// JSON returns the message map for the locale as JSON for use in script tags (e.g. window.I18N).
func JSON(locale string) template.JS {
	if locale == "" {
		locale = DefaultLocale
	}
	msgs, ok := Messages[locale]
	if !ok {
		msgs = En
	}
	raw, _ := json.Marshal(msgs)
	return template.JS(raw)
}
