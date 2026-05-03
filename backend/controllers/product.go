package controllers

import (
	"communityGroupBuying/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProductList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID, _ := strconv.Atoi(c.Query("category_id"))
	keyword := c.Query("keyword")
	isRecommend := c.Query("is_recommend")
	isNew := c.Query("is_new")

	var products []models.Product
	var total int64

	query := models.DB.Model(&models.Product{}).Where("status = 1")

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if isRecommend == "1" {
		query = query.Where("is_recommend = 1")
	}
	if isNew == "1" {
		query = query.Where("is_new = 1")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Category").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      products,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func ProductDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var product models.Product
	if models.DB.Preload("Category").First(&product, id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "商品不存在",
		})
		return
	}

	var comments []models.Comment
	models.DB.Where("product_id = ? AND status = 1", id).Preload("User").Order("created_at DESC").Limit(5).Find(&comments)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"product":  product,
			"comments": comments,
		},
	})
}

func RecommendProducts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	var products []models.Product
	models.DB.Where("status = 1 AND is_recommend = 1").Order("sort ASC, created_at DESC").Limit(limit).Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": products,
	})
}

func CategoryList(c *gin.Context) {
	var categories []models.Category
	models.DB.Where("status = 1").Order("sort ASC").Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": categories,
	})
}

func BannerList(c *gin.Context) {
	typ := c.DefaultQuery("type", "1")
	
	var banners []models.Banner
	models.DB.Where("status = 1 AND type = ?", typ).Order("sort ASC").Find(&banners)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": banners,
	})
}

func getPageInfo(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func buildProductQuery(c *gin.Context) *gorm.DB {
	query := models.DB.Model(&models.Product{})

	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	return query
}
