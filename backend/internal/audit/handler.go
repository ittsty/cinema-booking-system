package audit

import (
	"context"
	"net/http"
	"time"

	"cinema-booking/internal/models"
	"cinema-booking/pkg/mongo"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetLogsHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := mongo.DB.Collection("audit_logs")

	filter := bson.M{}

	if event := c.Query("event"); event != "" {
		filter["event"] = event
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var logs []models.AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, logs)
}
