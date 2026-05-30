package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/leno23/ai-api-gateway/internal/repository"
	"github.com/leno23/ai-api-gateway/internal/service"
)

type CatalogDeps struct {
	Repos *repository.Repos
}

func ListTokenGroups(deps CatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := deps.Repos.ListTokenGroups()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token groups"})
			return
		}
		items := make([]gin.H, 0, len(rows))
		for _, g := range rows {
			items = append(items, gin.H{
				"slug":       g.Slug,
				"name":       g.Name,
				"multiplier": g.Multiplier,
			})
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
	}
}

func ListCatalogModels(deps CatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		var billingType *int16
		if bt := c.Query("billing_type"); bt != "" {
			if v, err := strconv.ParseInt(bt, 10, 16); err == nil {
				v16 := int16(v)
				billingType = &v16
			}
		}
		filter := repository.CatalogFilter{
			Provider:     c.Query("provider"),
			EndpointType: c.Query("endpoint_type"),
			BillingType:  billingType,
			Tag:          c.Query("tag"),
			Query:        c.Query("q"),
			Page:         page,
			PageSize:     pageSize,
		}
		rows, total, err := deps.Repos.ListCatalogModels(filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "models"})
			return
		}
		groupSlug := c.DefaultQuery("token_group", "default")
		multiplier := 1.0
		if g, err := deps.Repos.TokenGroupBySlug(groupSlug); err == nil {
			multiplier = g.Multiplier
		}
		items := make([]gin.H, 0, len(rows))
		for _, r := range rows {
			name := r.DisplayName
			if name == "" {
				name = r.Model
			}
			prices := service.ModelPriceBreakdownFrom(r, 1)
			pricesWithMult := service.ModelPriceBreakdownFrom(r, multiplier)
			items = append(items, gin.H{
				"model":          r.Model,
				"display_name":   name,
				"provider":       r.Provider,
				"endpoint_type":  r.EndpointType,
				"billing_type":   r.BillingType,
				"billing_label":  service.BillingLabel(r.BillingType),
				"tags":           []string(r.Tags),
				"prices":         prices,
				"prices_applied": pricesWithMult,
			})
		}
		providers, _ := deps.Repos.DistinctProviders()
		groups, _ := deps.Repos.ListTokenGroups()
		groupItems := make([]gin.H, 0, len(groups))
		for _, g := range groups {
			groupItems = append(groupItems, gin.H{
				"slug":       g.Slug,
				"name":       g.Name,
				"multiplier": g.Multiplier,
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"items":      items,
				"total":      total,
				"page":       page,
				"page_size":  pageSize,
				"multiplier": multiplier,
				"token_group": groupSlug,
			},
			"meta": gin.H{
				"providers":    providers,
				"token_groups": groupItems,
			},
		})
	}
}

func GetModelPrice(deps CatalogDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		modelName := c.Param("model")
		if modelName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "model required"})
			return
		}
		mp, err := deps.Repos.ModelPriceByName(modelName)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "lookup"})
			return
		}
		groupSlug := c.DefaultQuery("token_group", "default")
		multiplier := 1.0
		if g, err := deps.Repos.TokenGroupBySlug(groupSlug); err == nil {
			multiplier = g.Multiplier
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"model":         mp.Model,
				"display_name":  mp.DisplayName,
				"base":          service.ModelPriceBreakdownFrom(*mp, 1),
				"with_multiplier": service.ModelPriceBreakdownFrom(*mp, multiplier),
				"token_group":   groupSlug,
			},
		})
	}
}
