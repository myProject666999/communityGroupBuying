package controllers

import (
	"communityGroupBuying/models"
	"communityGroupBuying/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AdminLogin(c *gin.Context) {
	var loginForm struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	var admin models.Admin
	if models.DB.Where("username = ?", loginForm.Username).First(&admin).RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	if admin.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "账户已被禁用",
		})
		return
	}

	if !utils.CheckPassword(loginForm.Password, admin.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	token, err := utils.GenerateToken(admin.ID, admin.Username, "admin")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token":    token,
			"admin_id": admin.ID,
			"username": admin.Username,
			"nickname": admin.Nickname,
			"role":     admin.Role,
		},
	})
}

func AdminInfo(c *gin.Context) {
	adminID := c.GetUint("user_id")

	var admin models.Admin
	if models.DB.First(&admin, adminID).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": admin,
	})
}

func DashboardStats(c *gin.Context) {
	var userCount, productCount, orderCount int64
	var todaySales, totalSales float64

	models.DB.Model(&models.User{}).Count(&userCount)
	models.DB.Model(&models.Product{}).Where("status = 1").Count(&productCount)
	models.DB.Model(&models.Order{}).Where("status > 0").Count(&orderCount)

	today := time.Now().Format("2006-01-02")
	models.DB.Model(&models.Order{}).Where("status >= 1 AND DATE(created_at) = ?", today).
		Select("COALESCE(SUM(actual_price), 0)").Scan(&todaySales)
	models.DB.Model(&models.Order{}).Where("status >= 1").
		Select("COALESCE(SUM(actual_price), 0)").Scan(&totalSales)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"user_count":     userCount,
			"product_count":  productCount,
			"order_count":    orderCount,
			"today_sales":    todaySales,
			"total_sales":    totalSales,
		},
	})
}

func AdminList(c *gin.Context) {
	page, pageSize := getPageInfo(c)

	var admins []models.Admin
	var total int64

	query := models.DB.Model(&models.Admin{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&admins)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      admins,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func AddAdmin(c *gin.Context) {
	var admin models.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	if admin.Username == "" || admin.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名和密码不能为空",
		})
		return
	}

	var existing models.Admin
	if models.DB.Where("username = ?", admin.Username).First(&existing).RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(admin.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}
	admin.Password = hashedPassword

	if err := models.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": admin,
	})
}

func UpdateAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var existing models.Admin
	if models.DB.First(&existing, id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	var updateData struct {
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
		Status   *int   `json:"status"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	updates := make(map[string]interface{})
	if updateData.Nickname != "" {
		updates["nickname"] = updateData.Nickname
	}
	if updateData.Role != "" {
		updates["role"] = updateData.Role
	}
	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}
	if updateData.Password != "" {
		hashedPassword, err := utils.HashPassword(updateData.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "密码加密失败",
			})
			return
		}
		updates["password"] = hashedPassword
	}

	if err := models.DB.Model(&existing).Updates(updates).Error; err != nil {
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

func DeleteAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var admin models.Admin
	if models.DB.First(&admin, id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "管理员不存在",
		})
		return
	}

	if admin.Role == "super" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "不能删除超级管理员",
		})
		return
	}

	models.DB.Delete(&admin)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

func AdminUserList(c *gin.Context) {
	page, pageSize := getPageInfo(c)
	keyword := c.Query("keyword")
	status := c.Query("status")

	var users []models.User
	var total int64

	query := models.DB.Model(&models.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR phone LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func AdminUpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var existing models.User
	if models.DB.First(&existing, id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	var updateData struct {
		Nickname string  `json:"nickname"`
		Phone    string  `json:"phone"`
		Balance  float64 `json:"balance"`
		Integral int     `json:"integral"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	updates := make(map[string]interface{})
	if updateData.Nickname != "" {
		updates["nickname"] = updateData.Nickname
	}
	if updateData.Phone != "" {
		updates["phone"] = updateData.Phone
	}
	if updateData.Balance >= 0 {
		updates["balance"] = updateData.Balance
	}
	if updateData.Integral >= 0 {
		updates["integral"] = updateData.Integral
	}

	if err := models.DB.Model(&existing).Updates(updates).Error; err != nil {
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

func AdminUpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var data struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if err := models.DB.Model(&models.User{}).Where("id = ?", id).Update("status", data.Status).Error; err != nil {
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

func AdminProductList(c *gin.Context) {
	page, pageSize := getPageInfo(c)
	keyword := c.Query("keyword")
	categoryID := c.Query("category_id")
	status := c.Query("status")

	var products []models.Product
	var total int64

	query := models.DB.Model(&models.Product{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("Category").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&products)

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

func AdminAddProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	if product.Name == "" || product.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "商品名称和价格不能为空",
		})
		return
	}

	if err := models.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "添加失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加成功",
		"data": product,
	})
}

func AdminUpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var existing models.Product
	if models.DB.First(&existing, id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "商品不存在",
		})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	product.ID = uint(id)

	if err := models.DB.Model(&existing).Updates(product).Error; err != nil {
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

func AdminDeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	result := models.DB.Delete(&models.Product{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "商品不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

func AdminUpdateProductStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var data struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if err := models.DB.Model(&models.Product{}).Where("id = ?", id).Update("status", data.Status).Error; err != nil {
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

func AdminCategoryList(c *gin.Context) {
	var categories []models.Category
	models.DB.Order("sort ASC, created_at DESC").Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": categories,
	})
}

func AdminAddCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if category.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "分类名称不能为空",
		})
		return
	}

	if err := models.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "添加失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加成功",
		"data": category,
	})
}

func AdminUpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var existing models.Category
	if models.DB.First(&existing, id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "分类不存在",
		})
		return
	}

	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	category.ID = uint(id)

	if err := models.DB.Model(&existing).Updates(category).Error; err != nil {
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

func AdminDeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var count int64
	models.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "该分类下还有商品，不能删除",
		})
		return
	}

	result := models.DB.Delete(&models.Category{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "分类不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

func AdminOrderList(c *gin.Context) {
	page, pageSize := getPageInfo(c)
	orderNo := c.Query("order_no")
	status := c.Query("status")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	var orders []models.Order
	var total int64

	query := models.DB.Model(&models.Order{})
	if orderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+orderNo+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != "" {
		query = query.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(created_at) <= ?", endDate)
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

func AdminOrderDetail(c *gin.Context) {
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
	if models.DB.Where("id = ?", id).Preload("Items").Preload("Address").First(&order).RowsAffected == 0 {
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

func AdminShipOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var data struct {
		CourierNo string `json:"courier_no"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var order models.Order
	if models.DB.First(&order, id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "订单不存在",
		})
		return
	}

	if order.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "订单状态不正确",
		})
		return
	}

	now := time.Now()
	order.Status = 2
	order.ShippingTime = &now
	order.CourierNo = data.CourierNo

	if err := models.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "发货失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "发货成功",
	})
}

func AdminBannerList(c *gin.Context) {
	var banners []models.Banner
	models.DB.Order("sort ASC, created_at DESC").Find(&banners)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": banners,
	})
}

func AdminAddBanner(c *gin.Context) {
	var banner models.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if banner.Image == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "图片不能为空",
		})
		return
	}

	if err := models.DB.Create(&banner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "添加失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加成功",
		"data": banner,
	})
}

func AdminUpdateBanner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var existing models.Banner
	if models.DB.First(&existing, id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "轮播图不存在",
		})
		return
	}

	var banner models.Banner
	if err := c.ShouldBindJSON(&banner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	banner.ID = uint(id)

	if err := models.DB.Model(&existing).Updates(banner).Error; err != nil {
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

func AdminDeleteBanner(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	result := models.DB.Delete(&models.Banner{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "轮播图不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

func AdminPromotionList(c *gin.Context) {
	page, pageSize := getPageInfo(c)

	var promotions []models.Promotion
	var total int64

	query := models.DB.Model(&models.Promotion{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&promotions)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      promotions,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func AdminAddPromotion(c *gin.Context) {
	var promotion models.Promotion
	if err := c.ShouldBindJSON(&promotion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if promotion.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "标题不能为空",
		})
		return
	}

	if err := models.DB.Create(&promotion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "添加失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加成功",
		"data": promotion,
	})
}

func AdminUpdatePromotion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var existing models.Promotion
	if models.DB.First(&existing, id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "促销不存在",
		})
		return
	}

	var promotion models.Promotion
	if err := c.ShouldBindJSON(&promotion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	promotion.ID = uint(id)

	if err := models.DB.Model(&existing).Updates(promotion).Error; err != nil {
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

func AdminDeletePromotion(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	result := models.DB.Delete(&models.Promotion{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "促销不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

func AdminKnowledgeList(c *gin.Context) {
	page, pageSize := getPageInfo(c)

	var knowledgeList []models.Knowledge
	var total int64

	query := models.DB.Model(&models.Knowledge{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&knowledgeList)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      knowledgeList,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func AdminAddKnowledge(c *gin.Context) {
	var knowledge models.Knowledge
	if err := c.ShouldBindJSON(&knowledge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if knowledge.Title == "" || knowledge.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "标题和内容不能为空",
		})
		return
	}

	if err := models.DB.Create(&knowledge).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "添加失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加成功",
		"data": knowledge,
	})
}

func AdminUpdateKnowledge(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var existing models.Knowledge
	if models.DB.First(&existing, id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "知识不存在",
		})
		return
	}

	var knowledge models.Knowledge
	if err := c.ShouldBindJSON(&knowledge); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	knowledge.ID = uint(id)

	if err := models.DB.Model(&existing).Updates(knowledge).Error; err != nil {
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

func AdminDeleteKnowledge(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	result := models.DB.Delete(&models.Knowledge{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "知识不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

func AdminCommunityList(c *gin.Context) {
	page, pageSize := getPageInfo(c)

	var communities []models.Community
	var total int64

	query := models.DB.Model(&models.Community{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&communities)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      communities,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func AdminUpdateCommunityStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var data struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if err := models.DB.Model(&models.Community{}).Where("id = ?", id).Update("status", data.Status).Error; err != nil {
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

func AdminDeleteCommunity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	result := models.DB.Delete(&models.Community{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "社区分享不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

func AdminForumList(c *gin.Context) {
	page, pageSize := getPageInfo(c)

	var forums []models.Forum
	var total int64

	query := models.DB.Model(&models.Forum{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&forums)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      forums,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func AdminUpdateForumStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var data struct {
		Status int `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if err := models.DB.Model(&models.Forum{}).Where("id = ?", id).Update("status", data.Status).Error; err != nil {
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

func AdminDeleteForum(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	result := models.DB.Delete(&models.Forum{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "帖子不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

func AdminNewsList(c *gin.Context) {
	page, pageSize := getPageInfo(c)

	var newsList []models.News
	var total int64

	query := models.DB.Model(&models.News{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&newsList)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":      newsList,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func AdminAddNews(c *gin.Context) {
	var news models.News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	if news.Title == "" || news.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "标题和内容不能为空",
		})
		return
	}

	if err := models.DB.Create(&news).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "添加失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "添加成功",
		"data": news,
	})
}

func AdminUpdateNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var existing models.News
	if models.DB.First(&existing, id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "资讯不存在",
		})
		return
	}

	var news models.News
	if err := c.ShouldBindJSON(&news); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	news.ID = uint(id)

	if err := models.DB.Model(&existing).Updates(news).Error; err != nil {
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

func AdminDeleteNews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	result := models.DB.Delete(&models.News{}, id)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "资讯不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
