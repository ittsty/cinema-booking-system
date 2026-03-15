package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user_id": c.GetString("user_id"),
		"email":   c.GetString("user_email"),
		"role":    c.GetString("role"),
	})
}
