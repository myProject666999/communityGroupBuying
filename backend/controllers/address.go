package controllers

import (
	"communityGroupBuying/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddressList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var addresses []models.Address
	models.DB.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": addresses,
	})
}

func AddressDetail(c *gin.Context) {
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

	var address models.Address
	if models.DB.Where("id = ? AND user_id = ?", id, userID).First(&address).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "地址不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": address,
	})
}

func AddAddress(c *gin.Context) {
	userID := c.GetUint("user_id")

	var address models.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	if address.Name == "" || address.Phone == "" || address.Address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "姓名、电话和详细地址不能为空",
		})
		return
	}

	address.UserID = userID

	if address.IsDefault == 1 {
		models.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", 0)
	} else {
		var count int64
		models.DB.Model(&models.Address{}).Where("user_id = ?", userID).Count(&count)
		if count == 0 {
			address.IsDefault = 1
		}
	}

	if err := models.DB.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "添加失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加成功",
		"data": address,
	})
}

func UpdateAddress(c *gin.Context) {
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

	var existingAddress models.Address
	if models.DB.Where("id = ? AND user_id = ?", id, userID).First(&existingAddress).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "地址不存在",
		})
		return
	}

	var address models.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	address.ID = uint(id)
	address.UserID = userID

	if address.IsDefault == 1 {
		models.DB.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", 0)
	}

	updates := make(map[string]interface{})
	if address.Name != "" {
		updates["name"] = address.Name
	}
	if address.Phone != "" {
		updates["phone"] = address.Phone
	}
	if address.Province != "" {
		updates["province"] = address.Province
	}
	if address.City != "" {
		updates["city"] = address.City
	}
	if address.District != "" {
		updates["district"] = address.District
	}
	if address.Address != "" {
		updates["address"] = address.Address
	}
	if address.IsDefault != existingAddress.IsDefault {
		updates["is_default"] = address.IsDefault
	}

	if err := models.DB.Model(&existingAddress).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

func DeleteAddress(c *gin.Context) {
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

	result := models.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Address{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "地址不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

func SetDefaultAddress(c *gin.Context) {
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

	var address models.Address
	if models.DB.Where("id = ? AND user_id = ?", id, userID).First(&address).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "地址不存在",
		})
		return
	}

	if err := models.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Address{}).Where("user_id = ?", userID).Update("is_default", 0).Error; err != nil {
			return err
		}
		if err := tx.Model(&address).Update("is_default", 1).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "设置失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "设置成功",
	})
}
