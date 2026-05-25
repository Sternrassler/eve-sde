// Package registry hält die Mapping-Registry zwischen JSONL-Datensätzen und
// Go-Typen für den SDE-Import. Die konkreten Mappings (Mappings) werden von
// sde-schema-gen in mappings_gen.go generiert.
package registry

import "reflect"

// SchemaMapping definiert das Mapping zwischen einer JSONL-Datei und einem Go-Typ.
type SchemaMapping struct {
	Name       string
	JSONLFile  string
	StructType reflect.Type
	Indices    []string
}
