package controllers

import (
	"communityGroupBuying/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FavoriteList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, pageSize := getPageInfo(c)

	var favorites []models.Favorite
	var total int64

	query := models.DB.Model(&models.Favorite{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Product").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&favorites)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      favorites,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func AddFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")

	var data struct {
		ProductID uint `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var product models.Product
	if models.DB.First(&product, data.ProductID).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "商品不存在",
		})
		return
	}

	var existing models.Favorite
	if models.DB.Where("user_id = ? AND product_id = ?", userID, data.ProductID).First(&existing).RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "已收藏",
		})
		return
	}

	favorite := models.Favorite{
		UserID:    userID,
		ProductID: data.ProductID,
	}

	if err := models.DB.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "收藏失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "收藏成功",
	})
}

func DeleteFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	result := models.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Favorite{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "收藏不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "取消收藏成功",
	})
}
