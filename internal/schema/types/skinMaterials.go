// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/skinMaterials.jsonl

package types

// SkinMaterials represents the schema for skinMaterials.jsonl
// This is a simplified struct - use actual schema documentation for production
type SkinMaterials struct {
	Key int64 `json:"_key"`
	DisplayName map[string]interface{} `json:"displayName"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
	MaterialSetID int64 `json:"materialSetID"`
}
