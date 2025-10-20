// Code generated from JSONL data analysis
// DO NOT EDIT manually - regenerate with scripts/fetch-schemas.sh --refresh
// Source: data/jsonl/skinLicenses.jsonl

package types

// SkinLicenses represents the schema for skinLicenses.jsonl
// This is a simplified struct - use actual schema documentation for production
type SkinLicenses struct {
	Key int64 `json:"_key"`
	Duration interface{} `json:"duration"`
	LicenseTypeID int64 `json:"licenseTypeID"`
	SkinID int64 `json:"skinID"`
}
