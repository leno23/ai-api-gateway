package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/model"
)

type AnnouncementDeps struct {
	DB *gorm.DB
}

func ListAnnouncements(deps AnnouncementDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		placement := c.DefaultQuery("placement", "home")
		now := time.Now()
		var rows []model.Announcement
		q := deps.DB.Where("status = ? AND placement = ?", 1, placement)
		q = q.Where("(starts_at IS NULL OR starts_at <= ?)", now).
			Where("(ends_at IS NULL OR ends_at >= ?)", now)
		if err := q.Order("id desc").Limit(20).Find(&rows).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "announcements"})
			return
		}
		items := make([]gin.H, 0, len(rows))
		for _, a := range rows {
			items = append(items, gin.H{
				"id":      a.ID,
				"title":   a.Title,
				"content": a.Content,
				"level":   a.Level,
			})
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
	}
}
