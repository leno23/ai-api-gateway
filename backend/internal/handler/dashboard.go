package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/leno23/ai-api-gateway/internal/model"
	"github.com/leno23/ai-api-gateway/internal/repository"
)

type DashboardDeps struct {
	Repos *repository.Repos
	Nodes []PortalNode
}

type PortalNode struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Region string `json:"region"`
}

func DashboardStats(deps DashboardDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		data, err := deps.Repos.DashboardStats(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "stats"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
	}
}

func DashboardCharts(deps DashboardDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		rangeKey := c.DefaultQuery("range", "7d")
		data, err := deps.Repos.DashboardCharts(userID, rangeKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "charts"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
	}
}

func DashboardNodes(deps DashboardDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "data": deps.Nodes})
	}
}

func ListUsageLogs(deps DashboardDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		f := usageLogFilterFromQuery(c, userID)
		rows, total, err := deps.Repos.ListUsageLogs(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "logs"})
			return
		}
		items := make([]gin.H, 0, len(rows))
		for _, row := range rows {
			items = append(items, serializeUsageLog(row))
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"items":     items,
				"total":     total,
				"page":      f.Page,
				"page_size": f.PageSize,
			},
		})
	}
}

func ListTaskLogs(deps DashboardDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := ctxUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		f := taskLogFilterFromQuery(c, userID)
		rows, total, err := deps.Repos.ListTaskLogs(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "tasks"})
			return
		}
		items := make([]gin.H, 0, len(rows))
		for _, row := range rows {
			items = append(items, serializeTaskLog(row))
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"items":     items,
				"total":     total,
				"page":      f.Page,
				"page_size": f.PageSize,
			},
		})
	}
}

func serializeUsageLog(row model.RequestLog) gin.H {
	var billing any
	if len(row.BillingDetail) > 0 {
		_ = json.Unmarshal(row.BillingDetail, &billing)
	}
	ttft := any(nil)
	if row.TimeToFirstMs != nil {
		ttft = *row.TimeToFirstMs
	}
	latency := any(nil)
	if row.LatencyMs != nil {
		latency = *row.LatencyMs
	}
	status := any(nil)
	if row.StatusCode != nil {
		status = *row.StatusCode
	}
	return gin.H{
		"id":                 row.ID,
		"request_id":         row.RequestID,
		"token_name":         row.TokenName,
		"token_group":        row.TokenGroup,
		"model":              row.Model,
		"request_path":       row.RequestPath,
		"prompt_tokens":      row.PromptTokens,
		"completion_tokens":  row.CompletionTokens,
		"total_tokens":       row.TotalTokens,
		"cost_quota":         row.CostQuota,
		"latency_ms":         latency,
		"time_to_first_ms":   ttft,
		"billing_detail":     billing,
		"status_code":        status,
		"error_message":      row.ErrorMessage,
		"created_at":         row.CreatedAt,
	}
}

func serializeTaskLog(row model.TaskLog) gin.H {
	durationMs := any(nil)
	if row.FinishedAt != nil {
		d := row.FinishedAt.Sub(row.SubmittedAt).Milliseconds()
		durationMs = d
	}
	return gin.H{
		"id":            row.ID,
		"task_id":       row.TaskID,
		"platform":      row.Platform,
		"type":          row.TaskType,
		"status":        row.Status,
		"progress":      row.Progress,
		"detail":        row.Detail,
		"submitted_at":  row.SubmittedAt,
		"finished_at":   row.FinishedAt,
		"duration_ms":   durationMs,
	}
}

func usageLogFilterFromQuery(c *gin.Context, userID int64) repository.UsageLogFilter {
	page, pageSize := pageFromQuery(c)
	return repository.UsageLogFilter{
		UserID:     userID,
		Start:      parseTimeQuery(c.Query("start")),
		End:        parseTimeQuery(c.Query("end")),
		TokenName:  c.Query("token_name"),
		Model:      c.Query("model"),
		RequestID:  c.Query("request_id"),
		TokenGroup: c.Query("token_group"),
		Page:       page,
		PageSize:   pageSize,
	}
}

func taskLogFilterFromQuery(c *gin.Context, userID int64) repository.TaskLogFilter {
	page, pageSize := pageFromQuery(c)
	return repository.TaskLogFilter{
		UserID:   userID,
		Start:    parseTimeQuery(c.Query("start")),
		End:      parseTimeQuery(c.Query("end")),
		TaskID:   c.Query("task_id"),
		Page:     page,
		PageSize: pageSize,
	}
}

func parseTimeQuery(v string) *time.Time {
	if v == "" {
		return nil
	}
	if t, err := time.Parse(time.RFC3339, v); err == nil {
		return &t
	}
	if t, err := time.Parse("2006-01-02", v); err == nil {
		return &t
	}
	return nil
}

func pageFromQuery(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
