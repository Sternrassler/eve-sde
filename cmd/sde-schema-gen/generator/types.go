package generator

import (
	"regexp"
	"strings"
)

// FileNameToTypeName konvertiert JSONL-Dateinamen zu Go-Typ-Namen
// blueprints.jsonl → Blueprints
func FileNameToTypeName(fileName string) string {
	name := strings.TrimSuffix(fileName, ".jsonl")
	return ToCamelCase(name, true)
}

// TypeNameToFileName konvertiert Typ-Namen zu Dateinamen
// Blueprints → blueprints
func TypeNameToFileName(typeName string) string {
	return toSnakeCase(typeName)
}

// ToCamelCase konvertiert snake_case oder kebab-case zu CamelCase
func ToCamelCase(s string, capitalizeFirst bool) string {
	// Entferne _key Prefix
	s = strings.TrimPrefix(s, "_")

	// Bereits CamelCase? (enthält Klein-dann-Groß)
	if regexp.MustCompile(`[a-z][A-Z]`).MatchString(s) {
		// Capitalize first letter if needed
		if capitalizeFirst && len(s) > 0 {
			return strings.ToUpper(s[:1]) + s[1:]
		}
		return s
	}

	// Split bei _ oder -
	words := regexp.MustCompile(`[_-]+`).Split(s, -1)
	result := make([]string, 0, len(words))

	for i, word := range words {
		if word == "" {
			continue
		}

		// Spezial-Behandlung für Abkürzungen
		upper := strings.ToUpper(word)
		if isAbbreviation(upper) {
			result = append(result, upper)
			continue
		}

		// Capitalize first letter
		if i == 0 && !capitalizeFirst {
			result = append(result, word)
		} else if len(word) > 0 {
			result = append(result, strings.ToUpper(word[:1])+strings.ToLower(word[1:]))
		}
	}

	return strings.Join(result, "")
}

// isAbbreviation prüft, ob ein Wort eine bekannte Abkürzung ist
func isAbbreviation(word string) bool {
	abbrevs := []string{
		"ID", "NPC", "CEO", "URL", "API", "HP", "EHP",
		"SDE", "UI", "AI", "XP", "DPS", "ETA",
	}
	for _, abbr := range abbrevs {
		if word == abbr {
			return true
		}
	}
	return false
}

// toSnakeCase konvertiert CamelCase zu snake_case
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
