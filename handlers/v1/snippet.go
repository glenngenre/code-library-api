package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/glenngenre/code-paste-service/models"
	"github.com/glenngenre/code-paste-service/services"
	"github.com/glenngenre/code-paste-service/validation"
)

// GetSnippets godoc
// @Summary Get all snippets
// @Description Retrieve a cursor-paginated list of all public code snippets
// @Tags snippets
// @Accept json
// @Produce json
// @Param search query string false "Search term to filter snippets by title, description, or language"
// @Param language query string false "Filter snippets by language"
// @Param cursor query string false "Pagination cursor (opaque string from previous response)"
// @Param page_size query int false "Number of items per page (default: 20, max: 100)"
// @Param direction query string false "Pagination direction: 'next' or 'prev' (default: next)"
// @Success 200 {object} models.SnippetCursorListResponse
// @Failure 400 {object} models.EmptyResponse
// @Failure 500 {object} models.EmptyResponse
// @Router /api/snippets [get]
func GetSnippets(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := services.GetSnippets(services.GetSnippetsParams{
		Search:    c.Query("search"),
		Language:  c.Query("language"),
		Cursor:    c.Query("cursor"),
		Direction: c.DefaultQuery("direction", "next"),
		PageSize:  pageSize,
	})
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrInvalidCursor) {
			status = http.StatusBadRequest
		}
		c.JSON(status, models.Response{Success: false, Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.SnippetCursorListResponse{
		Success: true,
		Data:    result.Snippets,
		Meta: models.CursorMeta{
			PageSize:   pageSize,
			NextCursor: result.NextCursor,
			PrevCursor: result.PrevCursor,
			HasMore:    result.HasMore,
		},
	})
}

// GetSnippetByID godoc
// @Summary Get a snippet by ID
// @Description Retrieve a single public code snippet by its unique ID
// @Tags snippets
// @Accept json
// @Produce json
// @Param id path string true "Snippet UUID"
// @Success 200 {object} models.SnippetResponse
// @Failure 404 {object} models.EmptyResponse
// @Failure 500 {object} models.EmptyResponse
// @Router /api/snippet/{id} [get]
func GetSnippetByID(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=60")

	snippet, err := services.FetchPublicSnippet(c.Param("id"))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrSnippetNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, models.Response{Success: false, Data: nil})
		return
	}

	c.JSON(http.StatusOK, models.Response{Success: true, Data: snippet})
}

// GetRawSnippetByID godoc
// @Summary Get raw snippet code by ID
// @Description Retrieve the raw code content of a public snippet by its unique ID. Returns plain text.
// @Tags snippets
// @Accept json
// @Produce text/plain
// @Param id path string true "Snippet UUID"
// @Success 200 {string} string "Raw code content"
// @Failure 404 {object} models.EmptyResponse
// @Failure 500 {object} models.EmptyResponse
// @Router /api/snippet/{id}/raw [get]
func GetRawSnippetByID(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=60")

	snippet, err := services.FetchPublicSnippet(c.Param("id"))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrSnippetNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, models.Response{Success: false, Data: nil})
		return
	}

	contentType := "text/plain; charset=utf-8"
	if v, ok := mimeTypes[strings.ToLower(snippet.Language)]; ok {
		contentType = v
	}

	c.Data(http.StatusOK, contentType, []byte(snippet.Code))
}

// CreateSnippet godoc
// @Summary Create a new snippet
// @Description Create a new code snippet with the provided details
// @Tags snippets
// @Accept json
// @Produce json
// @Param snippet body models.CreateSnippetRequest true "Snippet payload"
// @Success 201 {object} models.SnippetResponse
// @Failure 400 {object} models.EmptyResponse
// @Failure 500 {object} models.EmptyResponse
// @Router /api/snippet [post]
func CreateSnippet(c *gin.Context) {
	var req models.CreateSnippetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fieldErrors := validation.Format(err)
		if fieldErrors != nil {
			c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: fieldErrors})
		} else {
			c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: err.Error()})
		}
		return
	}

	snippet, err := services.CreateSnippet(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil})
		return
	}

	c.JSON(http.StatusCreated, models.Response{Success: true, Data: snippet})
}

// BatchCreateSnippets godoc
// @Summary Batch create snippets
// @Description Create multiple code snippets in a single request. Accepts an array of snippet objects.
// @Tags snippets
// @Accept json
// @Produce json
// @Param snippets body []models.CreateSnippetRequest true "Array of snippet payloads"
// @Success 201 {object} models.EmptyResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/snippets [post]
func BatchCreateSnippets(c *gin.Context) {
	var reqs []models.CreateSnippetRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		fieldErrors := validation.Format(err)
		if fieldErrors != nil {
			c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: fieldErrors})
		} else {
			c.JSON(http.StatusBadRequest, models.Response{Success: false, Data: err.Error()})
		}
		return
	}

	if err := services.BatchCreateSnippets(reqs); err != nil {
		c.JSON(http.StatusInternalServerError, models.EmptyResponse{Success: false})
		return
	}

	c.JSON(http.StatusCreated, models.EmptyResponse{Success: true})
}

// IncrementViews godoc
// @Summary Increment snippet view count
// @Description Increment the view count of a snippet by its ID. Deduplicates by hashing the client IP and snippet ID, so repeated calls from the same IP are ignored.
// @Tags snippets
// @Accept json
// @Produce json
// @Param id path string true "Snippet UUID"
// @Success 200 {object} models.EmptyResponse
// @Failure 404 {object} models.EmptyResponse
// @Failure 500 {object} models.EmptyResponse
// @Router /api/snippet/{id}/view [post]
func IncrementViews(c *gin.Context) {
	_, err := services.IncrementViews(c.Param("id"), c.ClientIP())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{Success: false, Data: nil})
		return
	}
	c.JSON(http.StatusOK, models.Response{Success: true, Data: nil})
}

// I have no idea where to put this, so here it is.
// A map of common programming languages to their MIME types for syntax
// highlighting when serving raw code content.
//
// This is used in the GetRawSnippetByID handler to set the appropriate
// Content-Type header based on the snippet's language.
//
// I didn't want to put this in the database, so that we don't have
// to maintain a separate table of languages and MIME types. (also so we
// don't have to query the database every time we serve raw content) We
// normally don't add new languages very often, so this should be sufficient
// for now.
var mimeTypes = map[string]string{
	"go":         "text/x-go; charset=utf-8",
	"python":     "text/x-python; charset=utf-8",
	"javascript": "text/javascript; charset=utf-8",
	"js":         "text/javascript; charset=utf-8",
	"typescript": "text/typescript; charset=utf-8",
	"ts":         "text/typescript; charset=utf-8",
	"jsx":        "text/jsx; charset=utf-8",
	"tsx":        "text/tsx; charset=utf-8",
	"react":      "text/jsx; charset=utf-8",
	"vue":        "text/javascript; charset=utf-8",
	"html":       "text/html; charset=utf-8",
	"css":        "text/css; charset=utf-8",
	"json":       "application/json; charset=utf-8",
	"yaml":       "application/yaml; charset=utf-8",
	"yml":        "application/yaml; charset=utf-8",
	"xml":        "application/xml; charset=utf-8",
	"markdown":   "text/markdown; charset=utf-8",
	"md":         "text/markdown; charset=utf-8",
	"shell":      "text/x-shellscript; charset=utf-8",
	"bash":       "text/x-shellscript; charset=utf-8",
	"sh":         "text/x-shellscript; charset=utf-8",
	"c":          "text/x-c; charset=utf-8",
	"cpp":        "text/x-c++; charset=utf-8",
	"csharp":     "text/x-csharp; charset=utf-8",
	"java":       "text/x-java-source; charset=utf-8",
	"rust":       "text/x-rust; charset=utf-8",
	"php":        "application/x-httpd-php; charset=utf-8",
	"ruby":       "text/x-ruby; charset=utf-8",
	"swift":      "text/x-swift; charset=utf-8",
	"kotlin":     "text/x-kotlin; charset=utf-8",
}
