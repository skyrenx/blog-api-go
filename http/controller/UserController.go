package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/skyrenx/blog-api-go/http/entities"
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

func Register(c *gin.Context) {
	var user entities.User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to run handler: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process the request",
		})
		return
	}
	err := service.Register(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process the request",
		})
		return
	}
	c.Status(http.StatusCreated)
}
