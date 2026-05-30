package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type NodePingResult struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Region    string `json:"region"`
	OK        bool   `json:"ok"`
	LatencyMs int64  `json:"latency_ms"`
	Status    string `json:"status,omitempty"`
}

func NodesPing(deps DashboardDeps) gin.HandlerFunc {
	client := &http.Client{Timeout: 8 * time.Second}

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		results := make([]NodePingResult, 0, len(deps.Nodes))

		for _, node := range deps.Nodes {
			target := strings.TrimRight(node.URL, "/") + "/health"
			start := time.Now()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
			if err != nil {
				results = append(results, NodePingResult{
					Name: node.Name, URL: node.URL, Region: node.Region,
					OK: false, LatencyMs: 0, Status: "bad url",
				})
				continue
			}

			resp, err := client.Do(req)
			latency := time.Since(start).Milliseconds()
			row := NodePingResult{
				Name:      node.Name,
				URL:       node.URL,
				Region:    node.Region,
				LatencyMs: latency,
			}
			if err != nil {
				row.OK = false
				row.Status = err.Error()
				results = append(results, row)
				continue
			}
			_ = resp.Body.Close()
			row.OK = resp.StatusCode >= 200 && resp.StatusCode < 300
			row.Status = strconv.Itoa(resp.StatusCode)
			results = append(results, row)
		}

		c.JSON(http.StatusOK, gin.H{"success": true, "data": results})
	}
}
