package config

import (
	"fmt"

	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

// LanguageTag is a custom type that wraps language.Tag for YAML marshaling
type LanguageTag struct {
	language.Tag
}

// UnmarshalYAML implements custom unmarshaling for language tags
func (lt *LanguageTag) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}

	tag, err := language.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid language tag %q: %w", s, err)
	}

	lt.Tag = tag
	return nil
}

// MarshalYAML implements custom marshaling for language tags
func (lt LanguageTag) MarshalYAML() (interface{}, error) {
	return lt.Tag.String(), nil
}
