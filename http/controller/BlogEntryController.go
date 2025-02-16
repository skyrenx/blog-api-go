package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/skyrenx/blog-api-go/http/entities"
	"github.com/skyrenx/blog-api-go/http/service"

	"github.com/gin-gonic/gin"
)

func AuroraExampleHandler(c *gin.Context) {
	err := service.Example()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to run example: %v\n", err) //TODO ?
		os.Exit(1)
	}

	c.JSON(http.StatusOK, gin.H{"message": "done"})
}

func CreateBlogEntry(c *gin.Context) {
	// Declare a BlogEntry instance
	var entry entities.BlogEntry

	// Bind the JSON body to the BlogEntry struct
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := service.CreateBlogEntry(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to run handler: %v\n", err) //TODO ?
		os.Exit(1)
	}
}
