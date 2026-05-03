package main

import (
	"communityGroupBuying/config"
	"communityGroupBuying/models"
	"communityGroupBuying/routes"
	"communityGroupBuying/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.dbname"),
		viper.GetString("database.charset"),
	)

	var err error
	models.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	models.DB.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Category{},
		&models.Product{},
		&models.ShoppingCart{},
		&models.Address{},
		&models.Order{},
		&models.OrderItem{},
		&models.Favorite{},
		&models.Comment{},
		&models.Knowledge{},
		&models.Community{},
		&models.Banner{},
		&models.Promotion{},
		&models.Recharge{},
		&models.Integral{},
		&models.Forum{},
		&models.News{},
	)

	initDefaultData()

	gin.SetMode(viper.GetString("server.mode"))
	r := gin.Default()

	routes.InitRoutes(r)

	port := viper.GetString("server.port")
	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

func initDefaultData() {
	var adminCount int64
	models.DB.Model(&models.Admin{}).Count(&adminCount)
	if adminCount == 0 {
		hashedPassword, _ := utils.HashPassword("admin123")
		admin := models.Admin{
			Username: "admin",
			Password: hashedPassword,
			Nickname: "超级管理员",
			Role:     "super",
			Status:   1,
		}
		models.DB.Create(&admin)
		log.Println("默认管理员已创建: admin / admin123")
	}

	var categoryCount int64
	models.DB.Model(&models.Category{}).Count(&categoryCount)
	if categoryCount == 0 {
		categories := []models.Category{
			{Name: "新鲜蔬菜", Icon: "🥬", Sort: 1},
			{Name: "时令水果", Icon: "🍎", Sort: 2},
			{Name: "肉禽蛋奶", Icon: "🥚", Sort: 3},
			{Name: "水产海鲜", Icon: "🦐", Sort: 4},
			{Name: "粮油米面", Icon: "🍚", Sort: 5},
			{Name: "土特产", Icon: "🌰", Sort: 6},
		}
		models.DB.Create(&categories)
		log.Println("默认分类已创建")
	}

	var productCount int64
	models.DB.Model(&models.Product{}).Count(&productCount)
	if productCount == 0 {
		products := []models.Product{
			{Name: "有机西红柿", CategoryID: 1, Price: 5.99, MarketPrice: 7.99, Stock: 100, Unit: "斤", Description: "新鲜有机西红柿，自然成熟", Status: 1, IsRecommend: 1},
			{Name: "山东大苹果", CategoryID: 2, Price: 8.99, MarketPrice: 12.99, Stock: 200, Unit: "斤", Description: "红富士苹果，脆甜多汁", Status: 1, IsRecommend: 1},
			{Name: "土鸡蛋", CategoryID: 3, Price: 1.5, MarketPrice: 2.0, Stock: 500, Unit: "个", Description: "农家散养土鸡蛋", Status: 1, IsRecommend: 1},
			{Name: "新鲜菠菜", CategoryID: 1, Price: 4.99, MarketPrice: 6.99, Stock: 80, Unit: "斤", Description: "鲜嫩菠菜，营养丰富", Status: 1, IsRecommend: 0},
			{Name: "赣南脐橙", CategoryID: 2, Price: 6.99, MarketPrice: 9.99, Stock: 150, Unit: "斤", Description: "赣南脐橙，甜蜜多汁", Status: 1, IsRecommend: 1},
			{Name: "五常大米", CategoryID: 5, Price: 39.9, MarketPrice: 49.9, Stock: 300, Unit: "袋/5kg", Description: "正宗五常大米，香糯可口", Status: 1, IsRecommend: 1},
		}
		models.DB.Create(&products)
		log.Println("默认商品已创建")
	}

	var bannerCount int64
	models.DB.Model(&models.Banner{}).Count(&bannerCount)
	if bannerCount == 0 {
		banners := []models.Banner{
			{Image: "https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20vegetables%20farm%20market%20banner%20green%20healthy&image_size=landscape_16_9", Link: "", Sort: 1, Status: 1},
			{Image: "https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20fruits%20assortment%20colorful%20market%20banner&image_size=landscape_16_9", Link: "", Sort: 2, Status: 1},
			{Image: "https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=organic%20food%20farm%20products%20healthy%20lifestyle%20banner&image_size=landscape_16_9", Link: "", Sort: 3, Status: 1},
		}
		models.DB.Create(&banners)
		log.Println("默认轮播图已创建")
	}
}
