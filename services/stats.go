package services

import (
	"github.com/glenngenre/code-paste-service/database"
	"github.com/glenngenre/code-paste-service/models"
)

func GetGlobalViewsCount() (int64, error) {
	var totalViews int64
	result := database.DB.Model(&models.Snippet{}).Select("SUM(views)").Scan(&totalViews)
	return totalViews, result.Error
}

func GetGlobalSnippetsCount() (int64, error) {
	var totalSnippets int64
	result := database.DB.Model(&models.Snippet{}).Count(&totalSnippets)
	return totalSnippets, result.Error
}

type UniqueLanguagesResult struct {
	Languages []string
	Total     int
}

func GetGlobalUniqueLanguages() (*UniqueLanguagesResult, error) {
	var languages []string
	result := database.DB.Model(&models.Snippet{}).Distinct("language").Pluck("language", &languages)
	if result.Error != nil {
		return nil, result.Error
	}
	return &UniqueLanguagesResult{
		Languages: languages,
		Total:     len(languages),
	}, nil
}
