package v1

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Snippet routes
	r.GET("/api/snippets", GetSnippets)
	r.GET("/api/snippet/:id", GetSnippetByID)
	r.GET("/api/snippet/:id/raw", GetRawSnippetByID)
	r.POST("/api/snippet", CreateSnippet)
	r.POST("/api/snippets", BatchCreateSnippets)
	r.POST("/api/snippet/:id/view", IncrementViews)

	// Stats routes
	r.GET("/api/stats/snippets", GetGlobalSnippetsCount)
	r.GET("/api/stats/views", GetGlobalViewsCount)
	r.GET("/api/stats/languages", GetGlobalUniqueLanguages)
}
