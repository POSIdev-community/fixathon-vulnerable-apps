package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {

	render(c, "index.html", gin.H{
		"authorized": isAuthorized(c),
	})
}

func render(c *gin.Context, name string, data gin.H) {
	accept := c.Request.Header.Get("Accept")
	switch accept {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data)
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, name, data)
	}
}
