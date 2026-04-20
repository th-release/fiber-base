package xss

import (
	"html"
	"slices"
	"strings"
	"sync"

	"cth.release/common"
	"github.com/microcosm-cc/bluemonday"
)

var (
	defaultSanitizer = NewSanitizer(Config{})
	defaultMu        sync.RWMutex
)

type Config struct {
	AllowElements []string
}

type Sanitizer struct {
	textPolicy *bluemonday.Policy
	htmlPolicy *bluemonday.Policy
}

func NewSanitizer(cfg Config) *Sanitizer {
	return &Sanitizer{
		textPolicy: bluemonday.StrictPolicy(),
		htmlPolicy: newHTMLPolicy(cfg),
	}
}

func NewSanitizerFromAppConfig(cfg *common.Config) *Sanitizer {
	if cfg == nil {
		return NewSanitizer(Config{})
	}

	return NewSanitizer(Config{
		AllowElements: cfg.XSS.AllowElements,
	})
}

func ConfigureFromAppConfig(cfg *common.Config) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	defaultSanitizer = NewSanitizerFromAppConfig(cfg)
}

func Escape(input string) string {
	return html.EscapeString(strings.TrimSpace(input))
}

func StripTags(input string) string {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultSanitizer.StripTags(input)
}

func SanitizeText(input string) string {
	return StripTags(input)
}

func SanitizeHTML(input string) string {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultSanitizer.SanitizeHTML(input)
}

func ContainsRisk(input string) bool {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultSanitizer.ContainsRisk(input)
}

func SanitizeMap(values map[string]string) map[string]string {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultSanitizer.SanitizeMap(values)
}

func SanitizeHTMLMap(values map[string]string) map[string]string {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultSanitizer.SanitizeHTMLMap(values)
}

func (s *Sanitizer) StripTags(input string) string {
	return strings.TrimSpace(s.textPolicy.Sanitize(input))
}

func (s *Sanitizer) SanitizeText(input string) string {
	return s.StripTags(input)
}

func (s *Sanitizer) SanitizeHTML(input string) string {
	return strings.TrimSpace(s.htmlPolicy.Sanitize(input))
}

func (s *Sanitizer) ContainsRisk(input string) bool {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return false
	}

	return trimmed != s.SanitizeHTML(trimmed)
}

func (s *Sanitizer) SanitizeMap(values map[string]string) map[string]string {
	sanitized := make(map[string]string, len(values))
	for key, value := range values {
		sanitized[key] = s.SanitizeText(value)
	}
	return sanitized
}

func (s *Sanitizer) SanitizeHTMLMap(values map[string]string) map[string]string {
	sanitized := make(map[string]string, len(values))
	for key, value := range values {
		sanitized[key] = s.SanitizeHTML(value)
	}
	return sanitized
}

func newHTMLPolicy(cfg Config) *bluemonday.Policy {
	policy := bluemonday.UGCPolicy()
	policy.AllowElements("p", "br", "hr", "span")
	policy.AllowAttrs("class").OnElements("code", "pre")
	allowConfiguredElements(policy, cfg.AllowElements)
	return policy
}

func allowConfiguredElements(policy *bluemonday.Policy, elements []string) {
	filtered := make([]string, 0, len(elements))
	for _, element := range elements {
		normalized := strings.ToLower(strings.TrimSpace(element))
		if normalized == "" || slices.Contains(filtered, normalized) {
			continue
		}
		filtered = append(filtered, normalized)
	}

	if len(filtered) > 0 {
		policy.AllowElements(filtered...)
	}
}
