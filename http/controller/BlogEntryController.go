package controller

import (
	"fmt"
	"net/http"
	"os"

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
