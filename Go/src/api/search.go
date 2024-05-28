package api

import (
	"net/http"
	"phdays-app/src/models"

	"github.com/gin-gonic/gin"
)

type SearchString struct {
	Keyword string `json:"search" binding:"required"`
}

func Search(c *gin.Context) {
	render(c, "search.html", gin.H{})
}

func SearchArticles(c *gin.Context) {
	var search SearchString
	if err := c.ShouldBindJSON(&search); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	articles := models.GetArticlesByKeyword(search.Keyword)
	c.JSON(http.StatusOK, articles)
}
