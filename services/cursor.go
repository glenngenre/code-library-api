package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/glenngenre/code-paste-service/models"
)

var ErrInvalidCursor = errors.New("invalid cursor")

func EncodeCursor(createdAt time.Time, snippetID string) string {
	payload := models.CursorPayload{CreatedAt: createdAt, SnippetID: snippetID}
	b, _ := json.Marshal(payload)
	return base64.RawURLEncoding.EncodeToString(b)
}

func DecodeCursor(cursor string) (models.CursorPayload, error) {
	b, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return models.CursorPayload{}, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
	}
	var payload models.CursorPayload
	if err := json.Unmarshal(b, &payload); err != nil {
		return models.CursorPayload{}, fmt.Errorf("%w: %v", ErrInvalidCursor, err)
	}
	return payload, nil
}
