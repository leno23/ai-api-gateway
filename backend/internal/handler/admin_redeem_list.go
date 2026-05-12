package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/leno23/ai-api-gateway/internal/repository"
)

type AdminRedeemListDeps struct {
	Repos *repository.Repos
}

// ListRedeemCodes paginates redeem_codes for admin (optional ?status=).
func ListRedeemCodes(deps AdminRedeemListDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		ps, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
		if page < 1 {
			page = 1
		}
		if ps < 1 {
			ps = 50
		}
		if ps > 200 {
			ps = 200
		}
		var st *int16
		if q := c.Query("status"); q != "" {
			v, err := strconv.ParseInt(q, 10, 16)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
				return
			}
			s := int16(v)
			st = &s
		}
		offset := (page - 1) * ps
		rows, total, err := deps.Repos.ListRedeemCodes(repository.RedeemListFilter{
			Status: st,
			Offset: offset,
			Limit:  ps,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "list redeem codes"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"items":     rows,
			"total":     total,
			"page":      page,
			"page_size": ps,
		})
	}
}

// RedeemStats returns counts per redeem_codes.status.
func RedeemStats(deps AdminRedeemListDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := deps.Repos.RedeemCodeStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "stats"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"by_status": stats})
	}
}
