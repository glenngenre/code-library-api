package services

import (
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/glenngenre/code-paste-service/database"
	"github.com/glenngenre/code-paste-service/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrSnippetNotFound = errors.New("snippet not found")

// FetchPublicSnippet retrieves a single public snippet by its snippet_id.
//
// Returns [ErrSnippetNotFound] if the snippet does not exist or is not public.
func FetchPublicSnippet(snippetID string) (*models.Snippet, error) {
	var snippet models.Snippet
	result := database.DB.
		Where("snippet_id = ? AND is_public = ?", snippetID, true).
		First(&snippet)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrSnippetNotFound
	}
	return &snippet, result.Error
}

type GetSnippetsParams struct {
	Search    string
	Language  string
	Cursor    string
	Direction string
	PageSize  int
}

type GetSnippetsResult struct {
	Snippets   []models.Snippet
	NextCursor *string
	PrevCursor *string
	HasMore    bool
}

func GetSnippets(params GetSnippetsParams) (*GetSnippetsResult, error) {
	var cursorCreatedAt time.Time
	var cursorSnippetID string
	hasCursor := false

	if params.Cursor != "" {
		decoded, err := DecodeCursor(params.Cursor)
		if err != nil {
			return nil, err
		}
		cursorCreatedAt = decoded.CreatedAt
		cursorSnippetID = decoded.SnippetID
		hasCursor = true
	}

	query := database.DB.Model(&models.Snippet{}).Where("is_public = ?", true)

	if params.Search != "" {
		like := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(language) LIKE ?",
			like, like, like,
		)
	}
	if params.Language != "" {
		query = query.Where("LOWER(language) = ?", strings.ToLower(params.Language))
	}

	if hasCursor {
		if params.Direction == "prev" {
			query = query.Where(
				"(created_at > ?) OR (created_at = ? AND snippet_id > ?)",
				cursorCreatedAt, cursorCreatedAt, cursorSnippetID,
			).Order("created_at ASC, snippet_id ASC")
		} else {
			query = query.Where(
				"(created_at < ?) OR (created_at = ? AND snippet_id < ?)",
				cursorCreatedAt, cursorCreatedAt, cursorSnippetID,
			).Order("created_at DESC, snippet_id DESC")
		}
	} else {
		query = query.Order("created_at DESC, snippet_id DESC")
	}

	var snippets []models.Snippet
	if err := query.Limit(params.PageSize + 1).Find(&snippets).Error; err != nil {
		return nil, err
	}

	if params.Direction == "prev" && hasCursor {
		slices.Reverse(snippets)
	}

	hasMore := len(snippets) > params.PageSize
	if hasMore {
		snippets = snippets[:params.PageSize]
	}

	var nextCursor, prevCursor *string
	if len(snippets) > 0 {
		if hasMore || params.Direction == "prev" {
			encoded := EncodeCursor(snippets[len(snippets)-1].CreatedAt, snippets[len(snippets)-1].SnippetID)
			nextCursor = &encoded
		}

		if hasCursor && params.Direction == "next" {
			encoded := EncodeCursor(snippets[0].CreatedAt, snippets[0].SnippetID)
			prevCursor = &encoded
		} else if params.Direction == "prev" && hasMore {
			encoded := EncodeCursor(snippets[0].CreatedAt, snippets[0].SnippetID)
			prevCursor = &encoded
		}
	}
	return &GetSnippetsResult{
		Snippets:   snippets,
		NextCursor: nextCursor,
		PrevCursor: prevCursor,
		HasMore:    hasMore,
	}, nil
}

// CreateSnippet generates a UUID and persists a new snippet.
func CreateSnippet(req models.CreateSnippetRequest) (*models.Snippet, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	snippet := models.Snippet{
		SnippetID:   id.String(),
		Title:       req.Title,
		Code:        req.Code,
		Language:    req.Language,
		Filename:    req.Filename,
		Description: req.Description,
		Category:    req.Category,
		IsPublic:    req.IsPublic,
	}

	if err := database.DB.Create(&snippet).Error; err != nil {
		return nil, err
	}

	return &snippet, nil
}

// BatchCreateSnippets generates UUIDs and bulk-inserts multiple snippets.
func BatchCreateSnippets(reqs []models.CreateSnippetRequest) error {
	snippets := make([]models.Snippet, 0, len(reqs))

	for _, req := range reqs {
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		snippets = append(snippets, models.Snippet{
			SnippetID:   id.String(),
			Title:       req.Title,
			Code:        req.Code,
			Language:    req.Language,
			Filename:    req.Filename,
			Description: req.Description,
			Category:    req.Category,
			IsPublic:    req.IsPublic,
		})
	}

	return database.DB.Create(&snippets).Error
}
