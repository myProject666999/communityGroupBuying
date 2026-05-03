package controllers

import (
	"communityGroupBuying/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CartList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var carts []models.ShoppingCart
	models.DB.Where("user_id = ?", userID).Preload("Product").Find(&carts)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": carts,
	})
}

func AddCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cartData struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity"`
	}

	if err := c.ShouldBindJSON(&cartData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var product models.Product
	if models.DB.First(&product, cartData.ProductID).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "商品不存在",
		})
		return
	}

	if product.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "商品已下架",
		})
		return
	}

	var existingCart models.ShoppingCart
	if models.DB.Where("user_id = ? AND product_id = ?", userID, cartData.ProductID).First(&existingCart).RowsAffected > 0 {
		existingCart.Quantity += cartData.Quantity
		if existingCart.Quantity > product.Stock {
			existingCart.Quantity = product.Stock
		}
		models.DB.Save(&existingCart)
	} else {
		quantity := cartData.Quantity
		if quantity <= 0 {
			quantity = 1
		}
		if quantity > product.Stock {
			quantity = product.Stock
		}
		cart := models.ShoppingCart{
			UserID:    userID,
			ProductID: cartData.ProductID,
			Quantity:  quantity,
			Selected:  1,
		}
		models.DB.Create(&cart)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加成功",
	})
}

func UpdateCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var updateData struct {
		ID       uint `json:"id" binding:"required"`
		Quantity int  `json:"quantity"`
		Selected *int `json:"selected"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var cart models.ShoppingCart
	if models.DB.Where("id = ? AND user_id = ?", updateData.ID, userID).First(&cart).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "购物车项不存在",
		})
		return
	}

	if updateData.Quantity > 0 {
		cart.Quantity = updateData.Quantity
	}
	if updateData.Selected != nil {
		cart.Selected = *updateData.Selected
	}

	models.DB.Save(&cart)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

func DeleteCart(c *gin.Context) {
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

	result := models.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.ShoppingCart{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "购物车项不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

func ClearCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	models.DB.Where("user_id = ?", userID).Delete(&models.ShoppingCart{})

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "清空成功",
	})
}
