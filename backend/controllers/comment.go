package controllers

import (
	"communityGroupBuying/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddComment(c *gin.Context) {
	userID := c.GetUint("user_id")

	var data struct {
		OrderItemID uint   `json:"order_item_id" binding:"required"`
		ProductID   uint   `json:"product_id" binding:"required"`
		OrderID     uint   `json:"order_id"`
		Content     string `json:"content" binding:"required"`
		Images      string `json:"images"`
		Rating      int    `json:"rating"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var orderItem models.OrderItem
	if models.DB.Where("id = ?", data.OrderItemID).First(&orderItem).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单项不存在",
		})
		return
	}

	if orderItem.Commented == 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "已评价",
		})
		return
	}

	rating := data.Rating
	if rating < 1 || rating > 5 {
		rating = 5
	}

	comment := models.Comment{
		UserID:      userID,
		ProductID:   data.ProductID,
		OrderID:     data.OrderID,
		OrderItemID: data.OrderItemID,
		Content:     data.Content,
		Images:      data.Images,
		Rating:      rating,
		Status:      1,
	}

	err := models.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}
		if err := tx.Model(&orderItem).Update("commented", 1).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "评价失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "评价成功",
	})
}

func CommentList(c *gin.Context) {
	productIDStr := c.Query("product_id")
	if productIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "缺少商品ID",
		})
		return
	}

	productID, _ := strconv.Atoi(productIDStr)
	page, pageSize := getPageInfo(c)

	var comments []models.Comment
	var total int64

	query := models.DB.Model(&models.Comment{}).Where("product_id = ? AND status = 1", productID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      comments,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
