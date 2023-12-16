package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

// GetTestGinContext returns a gin context for testing
func GetTestGinContext() *gin.Context {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	return ctx
}
