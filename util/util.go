package util

import (
	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}
