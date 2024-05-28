package api

import (
	"fmt"
	"net/http"
	"phdays-app/src/models"

	"github.com/gin-gonic/gin"
)

func Articles(c *gin.Context) {
	articles := models.GetAllArticles()
	if len(articles) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, articles)
	}
}

func Article(c *gin.Context) {
	id := c.Param("id")
	article, err := models.GetArticle(id)
	if err != nil || article == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	render(c, "article_template.html", gin.H{
		"content": article.Content,
		"title":   article.Title,
		"userId":  article.UserId,
		"author":  article.Author,
	})
}

func CreateArticle(c *gin.Context) {
	userId := c.GetString(loggedInUser)
	title := c.PostForm("title")
	content := c.PostForm("content")

	articleId, err := models.AddNewArticle(content, title, userId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating article")
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/articles/%v", articleId))
}

func ArticleCreate(c *gin.Context) {
	render(c, "article_create.html", nil)
}
