package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"cinema-booking/pkg/firebase"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})
			c.Abort()
			return
		}

		idToken := parts[1]

		client, err := firebase.App().Auth(context.Background())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to initialize firebase auth client",
			})
			c.Abort()
			return
		}

		token, err := client.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			c.Abort()
			return
		}

		email, _ := token.Claims["email"].(string)
		role := resolveRole(email)

		c.Set("user_id", token.UID)
		c.Set("user_email", email)
		c.Set("role", role)

		c.Next()
	}
}

func resolveRole(email string) string {
	adminEmails := strings.Split(os.Getenv("ADMIN_EMAILS"), ",")
	for _, adminEmail := range adminEmails {
		if strings.TrimSpace(strings.ToLower(adminEmail)) == strings.ToLower(email) {
			return "ADMIN"
		}
	}
	return "USER"
}
