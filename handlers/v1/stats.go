package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/glenngenre/code-paste-service/models"
	"github.com/glenngenre/code-paste-service/services"
)

// GetGlobalViewsCount godoc
// @Summary Get total view count
// @Description Retrieve the sum of views across all snippets
// @Tags stats
// @Produce json
// @Success 200 {object} models.GlobalViewsResponse
// @Failure 500 {object} models.EmptyResponse
// @Router /api/stats/views [get]
func GetGlobalViewsCount(c *gin.Context) {
	totalViews, err := services.GetGlobalViewsCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    gin.H{"total_views": totalViews},
	})
}

// GetGlobalSnippetsCount godoc
// @Summary Get total snippets count
// @Description Retrieve the total number of snippets created
// @Tags stats
// @Produce json
// @Success 200 {object} models.GlobalSnippetsResponse
// @Failure 500 {object} models.EmptyResponse
// @Router /api/stats/snippets [get]
func GetGlobalSnippetsCount(c *gin.Context) {
	totalSnippets, err := services.GetGlobalSnippetsCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    gin.H{"total_snippets": totalSnippets},
	})
}

// GetGlobalUniqueLanguages godoc
// @Summary Get unique programming languages
// @Description Retrieve a list of unique programming languages used in snippets
// @Tags stats
// @Produce json
// @Success 200 {object} models.GlobalLanguagesResponse
// @Failure 500 {object} models.EmptyResponse
// @Router /api/stats/languages [get]
func GetGlobalUniqueLanguages(c *gin.Context) {
	result, err := services.GetGlobalUniqueLanguages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil})
		return
	}
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: gin.H{
			"unique_languages": result.Languages,
			"total_languages":  result.Total,
		},
	})
}
