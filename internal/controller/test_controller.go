package controller

import (
	"df-post-maker/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TestController struct {
	uc *usecase.TestUseCase
}

func NewTestController(uc *usecase.TestUseCase) *TestController {
	return &TestController{uc: uc}
}

func (c *TestController) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/test")
	{
		group.GET("/ping", c.Ping)
	}
}

// Ping godoc
// @Summary Ping target host
// @Tags test
// @Success 200 {string} string "pong"
// @Failure 502 {string} string "error message"
// @Router /test/ping [get]
func (c *TestController) Ping(ctx *gin.Context) {
	if err := c.uc.Ping(); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	ctx.String(http.StatusOK, "pong")
}
