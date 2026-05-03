package controllers

import (
	"communityGroupBuying/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func KnowledgeList(c *gin.Context) {
	page, pageSize := getPageInfo(c)

	var knowledgeList []models.Knowledge
	var total int64

	query := models.DB.Model(&models.Knowledge{}).Where("status = 1")
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

func KnowledgeDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var knowledge models.Knowledge
	if models.DB.Where("id = ? AND status = 1", id).First(&knowledge).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "知识不存在",
		})
		return
	}

	models.DB.Model(&knowledge).Update("views", knowledge.Views+1)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": knowledge,
	})
}

func CommunityList(c *gin.Context) {
	page, pageSize := getPageInfo(c)

	var communities []models.Community
	var total int64

	query := models.DB.Model(&models.Community{}).Where("status = 1")
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

func AddCommunity(c *gin.Context) {
	userID := c.GetUint("user_id")

	var data struct {
		Content string `json:"content" binding:"required"`
		Images  string `json:"images"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	community := models.Community{
		UserID:   userID,
		Content:  data.Content,
		Images:   data.Images,
		Likes:    0,
		Comments: 0,
		Status:   1,
	}

	if err := models.DB.Create(&community).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "发布失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "发布成功",
		"data": community,
	})
}

func ForumList(c *gin.Context) {
	page, pageSize := getPageInfo(c)

	var forums []models.Forum
	var total int64

	query := models.DB.Model(&models.Forum{}).Where("status = 1")
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

func ForumDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var forum models.Forum
	if models.DB.Where("id = ? AND status = 1", id).Preload("User").First(&forum).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "帖子不存在",
		})
		return
	}

	models.DB.Model(&forum).Update("views", forum.Views+1)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": forum,
	})
}

func AddForum(c *gin.Context) {
	userID := c.GetUint("user_id")

	var data struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Images  string `json:"images"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误: " + err.Error(),
		})
		return
	}

	forum := models.Forum{
		UserID:  userID,
		Title:   data.Title,
		Content: data.Content,
		Images:  data.Images,
		Views:   0,
		Likes:   0,
		Status:  1,
	}

	if err := models.DB.Create(&forum).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "发布失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "发布成功",
		"data": forum,
	})
}

func NewsList(c *gin.Context) {
	page, pageSize := getPageInfo(c)

	var newsList []models.News
	var total int64

	query := models.DB.Model(&models.News{}).Where("status = 1")
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

func NewsDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	var news models.News
	if models.DB.Where("id = ? AND status = 1", id).First(&news).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "资讯不存在",
		})
		return
	}

	models.DB.Model(&news).Update("views", news.Views+1)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": news,
	})
}
