package controllers

import (
	"communityGroupBuying/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func generateOrderNo() string {
	now := time.Now()
	timestamp := now.Format("20060102150405")
	nanosecond := now.Nanosecond() / 1000
	return fmt.Sprintf("%s%06d", timestamp, nanosecond%1000000)
}

func CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	var orderData struct {
		AddressID uint   `json:"address_id" binding:"required"`
		CartIDs   []uint `json:"cart_ids"`
		Remark    string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&orderData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var address models.Address
	if models.DB.Where("id = ? AND user_id = ?", orderData.AddressID, userID).First(&address).RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "地址不存在",
		})
		return
	}

	var carts []models.ShoppingCart
	query := models.DB.Where("user_id = ? AND selected = 1", userID)
	if len(orderData.CartIDs) > 0 {
		query = query.Where("id IN ?", orderData.CartIDs)
	}
	query.Preload("Product").Find(&carts)

	if len(carts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "购物车为空",
		})
		return
	}

	var totalPrice float64
	var orderItems []models.OrderItem

	for _, cart := range carts {
		if cart.Product == nil {
			continue
		}
		if cart.Product.Status != 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "商品 " + cart.Product.Name + " 已下架",
			})
			return
		}
		if cart.Product.Stock < cart.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": 400,
				"msg":  "商品 " + cart.Product.Name + " 库存不足",
			})
			return
		}

		itemTotal := cart.Product.Price * float64(cart.Quantity)
		totalPrice += itemTotal

		orderItems = append(orderItems, models.OrderItem{
			ProductID:    cart.ProductID,
			ProductName:  cart.Product.Name,
			ProductImage: cart.Product.Image,
			Price:        cart.Product.Price,
			Quantity:     cart.Quantity,
			TotalPrice:   itemTotal,
		})
	}

	order := models.Order{
		OrderNo:     generateOrderNo(),
		UserID:      userID,
		AddressID:   orderData.AddressID,
		TotalPrice:  totalPrice,
		ActualPrice: totalPrice,
		Freight:     0,
		Status:      0,
		Remark:      orderData.Remark,
	}

	err := models.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for i := range orderItems {
			orderItems[i].OrderID = order.ID
			if err := tx.Create(&orderItems[i]).Error; err != nil {
				return err
			}

			if err := tx.Model(&models.Product{}).Where("id = ?", orderItems[i].ProductID).
				Update("sales", gorm.Expr("sales + ?", orderItems[i].Quantity)).Error; err != nil {
				return err
			}

			if err := tx.Model(&models.Product{}).Where("id = ?", orderItems[i].ProductID).
				Update("stock", gorm.Expr("stock - ?", orderItems[i].Quantity)).Error; err != nil {
				return err
			}
		}

		var cartIDs []uint
		for _, cart := range carts {
			cartIDs = append(cartIDs, cart.ID)
		}
		if err := tx.Delete(&models.ShoppingCart{}, cartIDs).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建订单失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": gin.H{
			"order_id":  order.ID,
			"order_no":  order.OrderNo,
			"total_price": totalPrice,
		},
	})
}

func OrderList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, pageSize := getPageInfo(c)
	status := c.Query("status")

	var orders []models.Order
	var total int64

	query := models.DB.Model(&models.Order{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Items").Preload("Address").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&orders)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      orders,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func OrderDetail(c *gin.Context) {
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

	var order models.Order
	if models.DB.Where("id = ? AND user_id = ?", id, userID).Preload("Items").Preload("Address").First(&order).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": order,
	})
}

func CancelOrder(c *gin.Context) {
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

	var order models.Order
	if models.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	if order.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "只能取消待支付订单",
		})
		return
	}

	order.Status = -1
	models.DB.Save(&order)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "取消成功",
	})
}

func PayOrder(c *gin.Context) {
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

	var order models.Order
	if models.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	if order.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "订单状态不正确",
		})
		return
	}

	var user models.User
	models.DB.First(&user, userID)

	if user.Balance < order.ActualPrice {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "余额不足",
		})
		return
	}

	now := time.Now()
	err = models.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Update("balance", user.Balance-order.ActualPrice).Error; err != nil {
			return err
		}

		order.Status = 1
		order.PayTime = &now
		if err := tx.Save(&order).Error; err != nil {
			return err
		}

		integral := int(order.ActualPrice)
		if integral > 0 {
			if err := tx.Model(&user).Update("integral", user.Integral+integral).Error; err != nil {
				return err
			}

			integralRecord := models.Integral{
				UserID:   userID,
				Type:     1,
				Integral: integral,
				Source:   "订单消费",
				Remark:   fmt.Sprintf("订单%s消费获得积分", order.OrderNo),
			}
			if err := tx.Create(&integralRecord).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "支付失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "支付成功",
	})
}

func ReceiveOrder(c *gin.Context) {
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

	var order models.Order
	if models.DB.Where("id = ? AND user_id = ?", id, userID).First(&order).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	if order.Status != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "订单状态不正确",
		})
		return
	}

	now := time.Now()
	order.Status = 3
	order.ReceivedTime = &now
	models.DB.Save(&order)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "确认收货成功",
	})
}
