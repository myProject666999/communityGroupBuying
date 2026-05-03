package controllers

import (
	"communityGroupBuying/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func generateRechargeNo() string {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	nanosecond := now.Nanosecond() / 1000
	return fmt.Sprintf("R%s%06d", timestamp, nanosecond%1000000)
}

func CreateRecharge(c *gin.Context) {
	userID := c.GetUint("user_id")

	var data struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if data.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "充值金额必须大于0",
		})
		return
	}

	var bonus float64
	switch {
	case data.Amount >= 500:
		bonus = 50
	case data.Amount >= 200:
		bonus = 20
	case data.Amount >= 100:
		bonus = 10
	}

	recharge := models.Recharge{
		UserID:  userID,
		OrderNo: generateRechargeNo(),
		Amount:  data.Amount,
		Bonus:   bonus,
		Status:  0,
	}

	if err := models.DB.Create(&recharge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建充值订单失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": recharge,
	})
}

func RechargeList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, pageSize := getPageInfo(c)

	var recharges []models.Recharge
	var total int64

	query := models.DB.Model(&models.Recharge{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&recharges)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      recharges,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
