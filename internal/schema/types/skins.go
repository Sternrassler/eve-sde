// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/skins.jsonl

package types

// Skins represents the schema for skins.jsonl
// This is a simplified struct - use actual schema documentation for production
type Skins struct {
	Key int64 `json:"_key"`
	AllowCCPDevs bool `json:"allowCCPDevs"`
	InternalName string `json:"internalName"`
	SkinMaterialID int64 `json:"skinMaterialID"`
	Types []interface{} `json:"types"`
	VisibleSerenity bool `json:"visibleSerenity"`
	VisibleTranquility bool `json:"visibleTranquility"`
}
