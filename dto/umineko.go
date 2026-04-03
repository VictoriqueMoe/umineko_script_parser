package dto

type UminekoQuote struct {
	BaseQuote
	HasRedTruth    bool `json:"hasRedTruth,omitempty" example:"false"`
	HasBlueTruth   bool `json:"hasBlueTruth,omitempty" example:"false"`
	HasGoldTruth   bool `json:"hasGoldTruth,omitempty" example:"false"`
	HasPurpleTruth bool `json:"hasPurpleTruth,omitempty" example:"false"`
}
