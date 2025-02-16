package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/skyrenx/blog-api-go/http/entities"
	"github.com/skyrenx/blog-api-go/http/service"

	"github.com/gin-gonic/gin"
)

func AuroraExampleHandler(c *gin.Context) {
	err := service.Example()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to run example: %v\n", err) //TODO ?
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process the request",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "done"})
}

func GetBlogEntries(c *gin.Context) {
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("pageNumber", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "1"))
	blogEntries, totalPages, err := service.GetBlogEntries(pageNumber, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process the request",
		})
		fmt.Fprintf(os.Stderr, "Unable to run handler: %v\n", err) //TODO ?
		return

	}
	c.JSON(http.StatusOK, gin.H{"blog_entries": blogEntries, "page_count": totalPages})
}

func GetBlogEntrySummaries(c *gin.Context) {
	pageNumber, _ := strconv.Atoi(c.DefaultQuery("pageNumber", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "1"))
	blogEntries, totalPages, err := service.GetBlogEntrySummaries(pageNumber, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process the request",
		})
		fmt.Fprintf(os.Stderr, "Unable to run handler: %v\n", err) //TODO ?
		return

	}
	c.JSON(http.StatusOK, gin.H{"blog_entry_summaries": blogEntries, "page_count": totalPages})
}

func GetBlogEntryById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	blogEntry, err := service.GetBlogEntryById(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to run handler: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process the request",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"blog_entry": blogEntry,
	})
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
		fmt.Fprintf(os.Stderr, "Unable to run handler: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process the request",
		})
		return
	}
}
