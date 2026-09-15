package v1

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, basePath string) {
	base := r.Group(basePath)

	base.GET("/api/snippets", GetSnippets)
	base.GET("/api/snippet/:id", GetSnippetByID)
	base.GET("/api/snippet/:id/raw", GetRawSnippetByID)
	base.POST("/api/snippet", CreateSnippet)
	base.POST("/api/snippets", BatchCreateSnippets)
	base.POST("/api/snippet/:id/view", IncrementViews)

	base.GET("/api/stats/snippets", GetGlobalSnippetsCount)
	base.GET("/api/stats/views", GetGlobalViewsCount)
	base.GET("/api/stats/languages", GetGlobalUniqueLanguages)
}
