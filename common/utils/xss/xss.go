package xss

import (
	"html"
	"regexp"
	"strings"
)

var (
	tagPattern        = regexp.MustCompile(`(?is)<[^>]*>`)
	scriptTagPattern  = regexp.MustCompile(`(?is)<\s*script[^>]*>.*?<\s*/\s*script\s*>`)
	eventAttrPattern  = regexp.MustCompile(`(?i)\bon[a-z]+\s*=`)
	jsProtocolPattern = regexp.MustCompile(`(?i)(javascript:|data:text/html)`)
)

func Escape(input string) string {
	return html.EscapeString(strings.TrimSpace(input))
}

func StripTags(input string) string {
	withoutScripts := scriptTagPattern.ReplaceAllString(input, "")
	return strings.TrimSpace(tagPattern.ReplaceAllString(withoutScripts, ""))
}

func SanitizeText(input string) string {
	return Escape(StripTags(input))
}

func ContainsRisk(input string) bool {
	return scriptTagPattern.MatchString(input) ||
		eventAttrPattern.MatchString(input) ||
		jsProtocolPattern.MatchString(input)
}

func SanitizeMap(values map[string]string) map[string]string {
	sanitized := make(map[string]string, len(values))
	for key, value := range values {
		sanitized[key] = SanitizeText(value)
	}
	return sanitized
}
