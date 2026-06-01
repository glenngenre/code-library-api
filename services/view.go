package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/glenngenre/code-paste-service/database"
	"github.com/glenngenre/code-paste-service/models"
	"gorm.io/gorm"
)

func ViewHash(ip, snippetID string) string {
	key := []byte(os.Getenv("VIEW_SALT"))
	h := hmac.New(sha256.New, key)
	h.Write([]byte(ip + snippetID))
	return hex.EncodeToString(h.Sum(nil))
}

// Returns (alreadyViewed bool, err error)
func IncrementViews(snippetID, ip string) (bool, error) {
	hash := ViewHash(ip, snippetID)

	var view models.SnippetView
	if err := database.DB.Where("snippet_id = ? AND hash = ?", snippetID, hash).First(&view).Error; err == nil {
		return true, nil
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&models.SnippetView{SnippetID: snippetID, Hash: hash}).Error; err != nil {
			return err
		}
		return tx.Model(&models.Snippet{}).
			Where("snippet_id = ?", snippetID).
			UpdateColumn("views", gorm.Expr("views + 1")).Error
	})

	return false, err
}
