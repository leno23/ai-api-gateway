package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/leno23/ai-api-gateway/internal/model"
	"github.com/leno23/ai-api-gateway/internal/repository"
)

type AdminUserDeps struct {
	Repos *repository.Repos
}

type patchUserStatusReq struct {
	Status int16 `json:"status" binding:"required"`
}

// PatchUserStatus sets users.status (e.g. ban with 0).
func PatchUserStatus(deps AdminUserDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id"})
			return
		}
		var req patchUserStatusReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Status != model.UserStatusActive && req.Status != model.UserStatusDisabled {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}
		if err := deps.Repos.UpdateUserStatus(id, req.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
