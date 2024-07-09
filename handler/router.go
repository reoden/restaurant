package handler

import (
	"restaurant/pkgs"

	"github.com/gin-gonic/gin"
)

func Router(e *gin.Engine) {
	// 前台操作
	restaurant := e.Group("/restaurant/v1")
	{
		accountHandler := NewAccountHandler()
		restaurant.GET("/self", pkgs.WrapperHandler(accountHandler.Self))
	}

	{
		loginHandler := NewLoginHandler()
		restaurant.POST("/login", pkgs.WrapperHandler(loginHandler.UserLogin))
		restaurant.POST("/sendcode", pkgs.WrapperHandler(loginHandler.SendCode))
	}

	{
		appHandler := NewAppHandler()

		// restaurant.DELETE("/apps/:id", pkgs.WrapperHandler(appHandler.Delete))
		// restaurant.PUT("/apps/:id", pkgs.WrapperHandler(appHandler.Update))
		restaurant.GET("/apps/search", pkgs.WrapperHandler(appHandler.Search))
		restaurant.GET("/apps/self", pkgs.WrapperHandler(appHandler.Self))
		restaurant.POST("/apps/pic", pkgs.WrapperHandler(UploadPicture))
		restaurant.GET("/apps/:id", pkgs.WrapperHandler(appHandler.Read))
	}

	{
		bannerHandler := NewBannerHandler()
		restaurant.GET("/banner", pkgs.WrapperHandler(bannerHandler.List))
	}

	{
		chefHandler := NewChefHandler()
		// restaurant.GET("/chef", pkgs.WrapperHandler(chefHandler.List))
		restaurant.GET("/chef/:id", pkgs.WrapperHandler(chefHandler.Read))
	}

	{
		chefShowListHandler := NewChefShowListHandler()
		restaurant.GET("/chefShowList", pkgs.WrapperHandler(chefShowListHandler.List))
	}

	{
		chefTotalHandler := NewChefTotalHandler()
		restaurant.GET("/chef", pkgs.WrapperHandler(chefTotalHandler.List))
	}

	{
		guestHandler := NewGuestHandler()
		// restaurant.GET("/guest", pkgs.WrapperHandler(guestHandler.List))
		restaurant.GET("/guest/:id", pkgs.WrapperHandler(guestHandler.Read))
	}

	{
		guestShowListHandler := NewGuestShowListHandler()
		restaurant.GET("/guestShowList", pkgs.WrapperHandler(guestShowListHandler.List))
	}

	{
		guestTotalHandler := NewGuestTotalHandler()
		restaurant.GET("/guest", pkgs.WrapperHandler(guestTotalHandler.List))
	}

	{
		newsHandler := NewNewsHandler()
		// restaurant.GET("/news", pkgs.WrapperHandler(newsHandler.List))
		restaurant.GET("/news/:id", pkgs.WrapperHandler(newsHandler.Read))
	}

	{
		newsShowListHandler := NewNewsShowListHandler()
		restaurant.GET("/newsShowList", pkgs.WrapperHandler(newsShowListHandler.List))
	}

	{
		newsTotalHandler := NewNewsTotalHandler()
		restaurant.GET("/news", pkgs.WrapperHandler(newsTotalHandler.List))
	}

	// 后台操作
	admin := e.Group("/restaurant/admin")
	{
		loginHandler := NewLoginHandler()
		admin.POST("/login", pkgs.WrapperHandler(loginHandler.AdminLogin))
	}

	{
		appHandler := NewAppHandler()
		admin.GET("/apps", pkgs.WrapperHandler(appHandler.List))
		admin.POST("/apps/submit", pkgs.WrapperHandler(appHandler.Create))
		admin.POST("/apps/save", pkgs.WrapperHandler(appHandler.Save))
		admin.GET("/apps/:id/detail", pkgs.WrapperHandler(appHandler.Get))
		admin.PUT("/apps/:id", pkgs.WrapperHandler(appHandler.UpdateStatus)) //管理员执行的审核操作
		admin.GET("/apps/:id/download", pkgs.WrapperHandler(appHandler.DownloadDoc))
	}

	{
		bannerHandler := NewBannerHandler()
		admin.POST("/banner", pkgs.WrapperHandler(bannerHandler.Create))
		admin.PUT("/banner", pkgs.WrapperHandler(bannerHandler.Update))
		admin.DELETE("/banner/:id", pkgs.WrapperHandler(bannerHandler.Delete))
	}

	{
		chefTotalHandler := NewChefTotalHandler()
		admin.POST("/chef/submit", pkgs.WrapperHandler(chefTotalHandler.Create))
		admin.DELETE("/chef/:id", pkgs.WrapperHandler(chefTotalHandler.Delete))
	}

	{
		chefShowListHandler := NewChefShowListHandler()
		admin.PUT("/chefShowList", pkgs.WrapperHandler(chefShowListHandler.Update))
	}

	{
		guestTotalHandler := NewGuestTotalHandler()
		admin.POST("/guest/submit", pkgs.WrapperHandler(guestTotalHandler.Create))
		admin.DELETE("/guest/:id", pkgs.WrapperHandler(guestTotalHandler.Delete))
	}

	{
		guestShowListHandler := NewGuestShowListHandler()
		admin.PUT("/guestShowList", pkgs.WrapperHandler(guestShowListHandler.Update))
	}

	{
		newsTotalHandler := NewNewsTotalHandler()
		admin.POST("/news/submit", pkgs.WrapperHandler(newsTotalHandler.Create))
		admin.DELETE("/news/:id", pkgs.WrapperHandler(newsTotalHandler.Delete))
	}

	{
		newsShowListHandler := NewNewsShowListHandler()
		admin.PUT("/newsShowList", pkgs.WrapperHandler(newsShowListHandler.Update))
	}
}
