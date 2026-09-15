package middleware

import (
	"bytes"
	"embed"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/glenngenre/code-paste-service/models"
)

//go:embed templates/rate_limit.html
var templateFS embed.FS
var tmpl *template.Template

func init() {
	tmpl = template.Must(template.ParseFS(templateFS, "templates/rate_limit.html"))
}

const blockDuration = 30 * time.Minute
const visitorLifespan = 3 * time.Minute
const cleanupInterval = 5 * time.Minute

type visitor struct {
	readLimiter  *rate.Limiter
	writeLimiter *rate.Limiter
	lastSeen     time.Time
}

var (
	visitors   = make(map[string]*visitor)
	blockedIPs = make(map[string]time.Time)
	mu         sync.Mutex
	blockMu    sync.RWMutex
	db         *gorm.DB
)

func Init(database *gorm.DB) {
	db = database

	var records []models.RateLimit
	if err := db.Where("unblock_at > ?", time.Now()).Find(&records).Error; err == nil {
		for _, r := range records {
			blockedIPs[r.IP] = r.UnblockAt
		}
	}

	go cleanupVisitors()
}

func cleanupVisitors() {
	for {
		time.Sleep(cleanupInterval)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > visitorLifespan {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func isBlocked(ip string) (bool, time.Duration) {
	blockMu.RLock()
	unblockAt, exists := blockedIPs[ip]
	blockMu.RUnlock()

	if !exists {
		return false, 0
	}

	now := time.Now()
	if now.Before(unblockAt) {
		return true, unblockAt.Sub(now)
	}

	blockMu.Lock()
	delete(blockedIPs, ip)
	blockMu.Unlock()

	return false, 0
}

func blockIP(ip, reason string) error {
	now := time.Now()
	unblockAt := now.Add(blockDuration)

	record := models.RateLimit{
		IP:        ip,
		BlockedAt: now,
		UnblockAt: unblockAt,
		Reason:    reason,
	}

	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ip"}},
		DoUpdates: clause.AssignmentColumns([]string{"blocked_at", "unblock_at", "reason"}),
	}).Create(&record).Error

	if err == nil {
		blockMu.Lock()
		blockedIPs[ip] = unblockAt
		blockMu.Unlock()
	}

	return err
}

func unblockIP(ip string) error {
	return db.Where("ip = ?", ip).Delete(&models.RateLimit{}).Error
}

func getVisitor(ip string) *visitor {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		v = &visitor{
			readLimiter:  rate.NewLimiter(10, 30),
			writeLimiter: rate.NewLimiter(2, 5),
			lastSeen:     time.Now(),
		}
		visitors[ip] = v
		return v
	}

	v.lastSeen = time.Now()
	return v
}

func resetVisitor(ip string) {
	mu.Lock()
	defer mu.Unlock()
	delete(visitors, ip)
}

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/static/") {
			c.Next()
			return
		}

		ip := c.ClientIP()

		if blocked, remaining := isBlocked(ip); blocked {
			c.Header("Retry-After", fmt.Sprintf("%d", int(remaining.Seconds())))
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusTooManyRequests, generateRateLimitHTML(ip, remaining))
			c.Abort()
			return
		}

		v := getVisitor(ip)

		limiter := v.writeLimiter
		if c.Request.Method == http.MethodGet {
			limiter = v.readLimiter
		}
		if !limiter.Allow() {
			if err := blockIP(ip, "rate_limit"); err != nil {
				_ = c.Error(fmt.Errorf("failed to block ip %s: %w", ip, err))
			}
			resetVisitor(ip)

			c.Header("Retry-After", fmt.Sprintf("%d", int(blockDuration.Seconds())))
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusTooManyRequests, generateRateLimitHTML(ip, blockDuration))
			c.Abort()
			return
		}

		c.Next()
	}
}
func generateRateLimitHTML(ip string, remaining time.Duration) string {
	unblockTime := time.Now().Add(remaining)

	data := struct {
		IP            string
		ResumeTime    string
		UnblockTimeMs int64
	}{
		IP:            ip,
		ResumeTime:    unblockTime.Format("3:04 PM"),
		UnblockTimeMs: unblockTime.UnixNano() / 1e6,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "Rate Limit Exceeded. Please try again later."
	}

	return buf.String()
}
