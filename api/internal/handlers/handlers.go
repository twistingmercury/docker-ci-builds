package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUUID(c *gin.Context) {
	t := uuid.New()
	c.JSON(http.StatusOK, gin.H{"uuid": t})
}
