// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/dogmaAttributes.jsonl

package types

// DogmaAttributes represents the schema for dogmaAttributes.jsonl
// This is a simplified struct - use actual schema documentation for production
type DogmaAttributes struct {
	Key int64 `json:"_key"`
	AttributeCategoryID int64 `json:"attributeCategoryID"`
	DataType int64 `json:"dataType"`
	DefaultValue int64 `json:"defaultValue"`
	Description string `json:"description"`
	DisplayWhenZero bool `json:"displayWhenZero"`
	HighIsGood bool `json:"highIsGood"`
	Name string `json:"name"`
	Published bool `json:"published"`
	Stackable bool `json:"stackable"`
}
