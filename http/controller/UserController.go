package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skyrenx/blog-api-go/http/service"
)

func GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	r, err := service.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
		return
	}
	c.JSON(http.StatusOK, *r)
}
