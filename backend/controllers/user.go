package controllers

import (
	"communityGroupBuying/models"
	"communityGroupBuying/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名和密码不能为空",
		})
		return
	}

	var existingUser models.User
	if models.DB.Where("username = ?", user.Username).First(&existingUser).RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名已存在",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}
	user.Password = hashedPassword
	user.Nickname = user.Username

	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "注册失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "注册成功",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func Login(c *gin.Context) {
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

	var user models.User
	if models.DB.Where("username = ?", loginForm.Username).First(&user).RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 403,
			"msg":  "账户已被禁用",
		})
		return
	}

	if !utils.CheckPassword(loginForm.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"msg":  "用户名或密码错误",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, "user")
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
			"user_id":  user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
		},
	})
}

func UserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if models.DB.First(&user, userID).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": user,
	})
}

func UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")

	var updateData struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
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
	if updateData.Avatar != "" {
		updates["avatar"] = updateData.Avatar
	}
	if updateData.Phone != "" {
		updates["phone"] = updateData.Phone
	}
	if updateData.Email != "" {
		updates["email"] = updateData.Email
	}

	if err := models.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
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

func UpdatePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var passwordData struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&passwordData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var user models.User
	if models.DB.First(&user, userID).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "用户不存在",
		})
		return
	}

	if !utils.CheckPassword(passwordData.OldPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "原密码错误",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(passwordData.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "密码加密失败",
		})
		return
	}

	if err := models.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新密码失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "密码修改成功",
	})
}
