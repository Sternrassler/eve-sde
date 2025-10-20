package generator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Schema repräsentiert ein analysiertes JSONL-Schema
type Schema struct {
	Fields map[string]*FieldInfo
}

// FieldInfo enthält Type-Informationen für ein Feld
type FieldInfo struct {
	GoType       string
	IsRequired   bool
	IsLocalized  bool
	SampleValues []interface{}
}

// AnalyzeJSONL analysiert eine JSONL-Datei und extrahiert Schema-Informationen
func AnalyzeJSONL(path string, maxLines int) (*Schema, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("konnte Datei nicht öffnen: %w", err)
	}
	defer file.Close()

	schema := &Schema{
		Fields: make(map[string]*FieldInfo),
	}

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() && lineCount < maxLines {
		lineCount++

		var data map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &data); err != nil {
			continue // Skip fehlerhafte Zeilen
		}

		// Analysiere jedes Feld
		for key, value := range data {
			field, exists := schema.Fields[key]
			if !exists {
				field = &FieldInfo{
					IsRequired:   key == "_key", // _key ist immer required
					SampleValues: make([]interface{}, 0, 3),
				}
				schema.Fields[key] = field
			}

			// Zähle Vorkommen (auch null)
			field.SampleValues = append(field.SampleValues[:0:0], field.SampleValues...) // Ensure capacity

			// Speichere Non-Null Sample-Werte (max 3)
			if value != nil {
				if len(field.SampleValues) < 3 {
					field.SampleValues = append(field.SampleValues, value)
				}

				// Inferiere Typ (ignoriere null)
				goType := inferGoType(value)

				// Typ-Logik: Ignoriere interface{} als Kandidat
				if goType == "interface{}" {
					continue
				}

				if field.GoType == "" || field.GoType == "interface{}" {
					// Erster echter Typ oder Override von interface{}
					field.GoType = goType
				} else if field.GoType != goType {
					// Spezialfall: int64 + float64 → float64 (JSON-Zahlen)
					if (field.GoType == "int64" && goType == "float64") ||
						(field.GoType == "float64" && goType == "int64") {
						field.GoType = "float64"
					} else {
						// Echter Typ-Konflikt → fallback zu interface{}
						field.GoType = "interface{}"
					}
				}

				// Prüfe auf LocalizedText
				if !field.IsLocalized {
					field.IsLocalized = isLocalizedText(value)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("fehler beim Lesen: %w", err)
	}

	// Markiere Felder als optional basierend auf Sample-Count
	// (Felder ohne Samples = nur null-Werte = optional)
	for key, field := range schema.Fields {
		if key == "_key" {
			field.IsRequired = true // _key immer required
			continue
		}
		// Wenn wir Sample-Werte haben und sie in allen Zeilen vorkommen → required
		if len(field.SampleValues) == lineCount && lineCount > 1 {
			field.IsRequired = true
		} else {
			field.IsRequired = false
		}
	}

	return schema, nil
}

// inferGoType bestimmt den Go-Typ aus einem JSON-Wert
func inferGoType(value interface{}) string {
	if value == nil {
		return "interface{}"
	}

	switch v := value.(type) {
	case bool:
		return "bool"
	case float64:
		// JSON dekodiert alle Zahlen als float64
		if v == float64(int64(v)) {
			return "int64"
		}
		return "float64"
	case string:
		return "string"
	case map[string]interface{}:
		if isLocalizedText(v) {
			return "LocalizedText"
		}
		return "map[string]interface{}"
	case []interface{}:
		if len(v) == 0 {
			return "[]interface{}"
		}
		// Typ des ersten Elements
		elemType := inferGoType(v[0])
		return "[]" + elemType
	default:
		return "interface{}"
	}
}

// isLocalizedText prüft, ob ein Objekt ein mehrsprachiges Text-Objekt ist
func isLocalizedText(value interface{}) bool {
	m, ok := value.(map[string]interface{})
	if !ok {
		return false
	}

	// Prüfe auf typische Sprach-Keys
	langKeys := []string{"de", "en", "es", "fr", "ja", "ko", "ru", "zh"}
	hasLang := false
	allStrings := true

	for k, v := range m {
		// Ist es ein Sprach-Key?
		isLangKey := false
		for _, lang := range langKeys {
			if k == lang {
				isLangKey = true
				break
			}
		}

		if isLangKey {
			hasLang = true
			if _, ok := v.(string); !ok {
				allStrings = false
			}
		} else {
			// Nicht-Sprach-Keys → kein LocalizedText
			return false
		}
	}

	return hasLang && allStrings
}
