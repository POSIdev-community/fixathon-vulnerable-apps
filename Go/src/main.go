package main

import (
	"context"
	"phdays-app/src/api"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := engine()
	r.Use(gin.Logger())
	if err := engine().Run(":8080"); err != nil {
		logrus.WithContext(context.Background()).Fatal("Unable to start:", err)
	}
}

func engine() *gin.Engine {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.LoadHTMLGlob("./src/html/*")
	r.Static("/static", "./src/static")

	r.GET("/", api.Index)
	r.GET("/login", api.LoginPage)
	r.GET("/search", api.Search)
	r.GET("/articles/:id", api.Article)
	r.GET("/profile/:userId", api.Profile)

	apiGroup := r.Group("/api")
	apiGroup.POST("/logout", api.Logout)
	apiGroup.POST("/search", api.SearchArticles)
	apiGroup.GET("/articles", api.Articles)
	apiGroup.POST("/login", api.Login)

	gr := r.Group("")
	gr.Use(api.AuthRequired)
	{
		gr.GET("/article_create", api.ArticleCreate)
		gr.GET("/my_profile", api.MyProfile)
	}

	apiGroup.Use(api.AuthRequired)
	{
		apiGroup.POST("profile/upload_photo", api.UploadPhoto)
		apiGroup.POST("profile/upload_photo_url", api.UploadPhotoUrl)
		apiGroup.POST("/article_create", api.CreateArticle)
	}

	return r
}
