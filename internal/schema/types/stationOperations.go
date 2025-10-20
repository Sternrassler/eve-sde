// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/stationOperations.jsonl

package types

// StationOperations represents the schema for stationOperations.jsonl
// This is a simplified struct - use actual schema documentation for production
type StationOperations struct {
	Key int64 `json:"_key"`
	ActivityID int64 `json:"activityID"`
	Border int64 `json:"border"`
	Corridor int64 `json:"corridor"`
	Description map[string]interface{} `json:"description"`
	De interface{} `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
	Fr string `json:"fr"`
	Ja string `json:"ja"`
	Ko string `json:"ko"`
	Ru interface{} `json:"ru"`
	Zh string `json:"zh"`
	Fringe int64 `json:"fringe"`
	Hub int64 `json:"hub"`
	ManufacturingFactor int64 `json:"manufacturingFactor"`
	OperationName map[string]interface{} `json:"operationName"`
	De interface{} `json:"de"`
	En string `json:"en"`
	Es string `json:"es"`
}
