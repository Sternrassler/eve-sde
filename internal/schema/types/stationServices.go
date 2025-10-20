// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/stationServices.jsonl

package types

// StationServices represents the schema for stationServices.jsonl
// This is a simplified struct - use actual schema documentation for production
type StationServices struct {
	Key int64 `json:"_key"`
	ServiceName map[string]interface{} `json:"serviceName"`
	De string `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru string `json:"ru"`
	Zh string `json:"zh"`
}
