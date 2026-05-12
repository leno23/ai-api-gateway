package service

import (
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/model"
)

// EstimateMaxCost returns conservative max quota for prompt + max completion tokens.
func EstimateMaxCost(db *gorm.DB, modelName string, promptTokens, maxOutputTokens int64) (int64, error) {
	var mp model.ModelPrice
	if err := db.Where("model = ?", modelName).First(&mp).Error; err != nil {
		return 0, err
	}
	in := (promptTokens * mp.PromptPrice) / 1000
	out := (maxOutputTokens * mp.CompletionPrice) / 1000
	unit := (promptTokens * mp.UnitPrice) / 1000
	return in + out + unit, nil
}

// ActualCost computes settled cost from actual token counts.
func ActualCost(db *gorm.DB, modelName string, inTok, outTok int64) (int64, error) {
	var mp model.ModelPrice
	if err := db.Where("model = ?", modelName).First(&mp).Error; err != nil {
		return 0, err
	}
	in := (inTok * mp.PromptPrice) / 1000
	out := (outTok * mp.CompletionPrice) / 1000
	return in + out, nil
}
