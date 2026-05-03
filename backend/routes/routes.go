package routes

import (
	"communityGroupBuying/controllers"
	"communityGroupBuying/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/register", controllers.Register)
			user.POST("/login", controllers.Login)
		}

		admin := api.Group("/admin")
		{
			admin.POST("/login", controllers.AdminLogin)
		}

		product := api.Group("/product")
		{
			product.GET("/list", controllers.ProductList)
			product.GET("/detail/:id", controllers.ProductDetail)
			product.GET("/recommend", controllers.RecommendProducts)
		}

		category := api.Group("/category")
		{
			category.GET("/list", controllers.CategoryList)
		}

		banner := api.Group("/banner")
		{
			banner.GET("/list", controllers.BannerList)
		}

		knowledge := api.Group("/knowledge")
		{
			knowledge.GET("/list", controllers.KnowledgeList)
			knowledge.GET("/detail/:id", controllers.KnowledgeDetail)
		}

		community := api.Group("/community")
		{
			community.GET("/list", controllers.CommunityList)
		}

		news := api.Group("/news")
		{
			news.GET("/list", controllers.NewsList)
			news.GET("/detail/:id", controllers.NewsDetail)
		}

		forum := api.Group("/forum")
		{
			forum.GET("/list", controllers.ForumList)
			forum.GET("/detail/:id", controllers.ForumDetail)
		}
	}

	auth := api.Group("")
	auth.Use(middleware.JWTAuth())
	{
		user := auth.Group("/user")
		{
			user.GET("/info", controllers.UserInfo)
			user.PUT("/info", controllers.UpdateUserInfo)
			user.PUT("/password", controllers.UpdatePassword)
		}

		cart := auth.Group("/cart")
		{
			cart.GET("/list", controllers.CartList)
			cart.POST("/add", controllers.AddCart)
			cart.PUT("/update", controllers.UpdateCart)
			cart.DELETE("/delete/:id", controllers.DeleteCart)
			cart.POST("/clear", controllers.ClearCart)
		}

		address := auth.Group("/address")
		{
			address.GET("/list", controllers.AddressList)
			address.GET("/detail/:id", controllers.AddressDetail)
			address.POST("/add", controllers.AddAddress)
			address.PUT("/update/:id", controllers.UpdateAddress)
			address.DELETE("/delete/:id", controllers.DeleteAddress)
			address.PUT("/default/:id", controllers.SetDefaultAddress)
		}

		order := auth.Group("/order")
		{
			order.POST("/create", controllers.CreateOrder)
			order.GET("/list", controllers.OrderList)
			order.GET("/detail/:id", controllers.OrderDetail)
			order.PUT("/cancel/:id", controllers.CancelOrder)
			order.PUT("/pay/:id", controllers.PayOrder)
			order.PUT("/receive/:id", controllers.ReceiveOrder)
		}

		favorite := auth.Group("/favorite")
		{
			favorite.GET("/list", controllers.FavoriteList)
			favorite.POST("/add", controllers.AddFavorite)
			favorite.DELETE("/delete/:id", controllers.DeleteFavorite)
		}

		comment := auth.Group("/comment")
		{
			comment.POST("/add", controllers.AddComment)
			comment.GET("/list", controllers.CommentList)
		}

		communityAuth := auth.Group("/community")
		{
			communityAuth.POST("/add", controllers.AddCommunity)
		}

		forumAuth := auth.Group("/forum")
		{
			forumAuth.POST("/add", controllers.AddForum)
		}

		recharge := auth.Group("/recharge")
		{
			recharge.POST("/create", controllers.CreateRecharge)
			recharge.GET("/list", controllers.RechargeList)
		}
	}

	adminAuth := api.Group("/admin")
	adminAuth.Use(middleware.JWTAuth(), middleware.AdminAuth())
	{
		adminAuth.GET("/info", controllers.AdminInfo)

		dashboard := adminAuth.Group("/dashboard")
		{
			dashboard.GET("/stats", controllers.DashboardStats)
		}

		adminUser := adminAuth.Group("/admin")
		{
			adminUser.GET("/list", controllers.AdminList)
			adminUser.POST("/add", controllers.AddAdmin)
			adminUser.PUT("/update/:id", controllers.UpdateAdmin)
			adminUser.DELETE("/delete/:id", controllers.DeleteAdmin)
		}

		userManage := adminAuth.Group("/user")
		{
			userManage.GET("/list", controllers.AdminUserList)
			userManage.PUT("/update/:id", controllers.AdminUpdateUser)
			userManage.PUT("/status/:id", controllers.AdminUpdateUserStatus)
		}

		productManage := adminAuth.Group("/product")
		{
			productManage.GET("/list", controllers.AdminProductList)
			productManage.POST("/add", controllers.AdminAddProduct)
			productManage.PUT("/update/:id", controllers.AdminUpdateProduct)
			productManage.DELETE("/delete/:id", controllers.AdminDeleteProduct)
			productManage.PUT("/status/:id", controllers.AdminUpdateProductStatus)
		}

		categoryManage := adminAuth.Group("/category")
		{
			categoryManage.GET("/list", controllers.AdminCategoryList)
			categoryManage.POST("/add", controllers.AdminAddCategory)
			categoryManage.PUT("/update/:id", controllers.AdminUpdateCategory)
			categoryManage.DELETE("/delete/:id", controllers.AdminDeleteCategory)
		}

		orderManage := adminAuth.Group("/order")
		{
			orderManage.GET("/list", controllers.AdminOrderList)
			orderManage.GET("/detail/:id", controllers.AdminOrderDetail)
			orderManage.PUT("/ship/:id", controllers.AdminShipOrder)
		}

		bannerManage := adminAuth.Group("/banner")
		{
			bannerManage.GET("/list", controllers.AdminBannerList)
			bannerManage.POST("/add", controllers.AdminAddBanner)
			bannerManage.PUT("/update/:id", controllers.AdminUpdateBanner)
			bannerManage.DELETE("/delete/:id", controllers.AdminDeleteBanner)
		}

		promotionManage := adminAuth.Group("/promotion")
		{
			promotionManage.GET("/list", controllers.AdminPromotionList)
			promotionManage.POST("/add", controllers.AdminAddPromotion)
			promotionManage.PUT("/update/:id", controllers.AdminUpdatePromotion)
			promotionManage.DELETE("/delete/:id", controllers.AdminDeletePromotion)
		}

		knowledgeManage := adminAuth.Group("/knowledge")
		{
			knowledgeManage.GET("/list", controllers.AdminKnowledgeList)
			knowledgeManage.POST("/add", controllers.AdminAddKnowledge)
			knowledgeManage.PUT("/update/:id", controllers.AdminUpdateKnowledge)
			knowledgeManage.DELETE("/delete/:id", controllers.AdminDeleteKnowledge)
		}

		communityManage := adminAuth.Group("/community")
		{
			communityManage.GET("/list", controllers.AdminCommunityList)
			communityManage.PUT("/status/:id", controllers.AdminUpdateCommunityStatus)
			communityManage.DELETE("/delete/:id", controllers.AdminDeleteCommunity)
		}

		forumManage := adminAuth.Group("/forum")
		{
			forumManage.GET("/list", controllers.AdminForumList)
			forumManage.PUT("/status/:id", controllers.AdminUpdateForumStatus)
			forumManage.DELETE("/delete/:id", controllers.AdminDeleteForum)
		}

		newsManage := adminAuth.Group("/news")
		{
			newsManage.GET("/list", controllers.AdminNewsList)
			newsManage.POST("/add", controllers.AdminAddNews)
			newsManage.PUT("/update/:id", controllers.AdminUpdateNews)
			newsManage.DELETE("/delete/:id", controllers.AdminDeleteNews)
		}
	}
}
