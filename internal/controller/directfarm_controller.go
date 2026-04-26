package controller

import (
	"df-post-maker/internal/dto"
	"df-post-maker/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DirectFarmController struct {
	uc *usecase.DirectFarmUseCase
}

func NewDirectFarmController(uc *usecase.DirectFarmUseCase) *DirectFarmController {
	return &DirectFarmController{uc: uc}
}

func (c *DirectFarmController) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/direct-farm")
	{
		group.POST("/auth", c.Auth)
		group.POST("/post/add", c.CreatePost)
		group.POST("/upload", c.Upload)
	}
}

// Auth godoc
// @Summary Authenticate via direct.farm
// @Tags direct-farm
// @Accept json
// @Produce json
// @Param request body dto.AuthRequestDto true "Login credentials"
// @Success 200 {object} dto.AuthResponseDto
// @Failure 400 {object} map[string]string
// @Failure 502 {object} map[string]string
// @Router /direct-farm/auth [post]
func (c *DirectFarmController) Auth(ctx *gin.Context) {
	var req dto.AuthRequestDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.uc.Auth(req)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreatePost godoc
// @Summary Create a post on direct.farm
// @Tags direct-farm
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body dto.CreatePostRequestDto true "Post data"
// @Success 200 {object} dto.CreatePostResponseDto
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 502 {object} map[string]string
// @Router /direct-farm/post/add [post]
func (c *DirectFarmController) CreatePost(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	var req dto.CreatePostRequestDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.SetDefaults()

	resp, err := c.uc.CreatePost(req, authHeader)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// Upload godoc
// @Summary Upload a file to direct.farm
// @Tags direct-farm
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param file formData file true "File to upload"
// @Success 200 {object} dto.UploadResponseDto
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 502 {object} map[string]string
// @Router /direct-farm/upload [post]
func (c *DirectFarmController) Upload(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required: " + err.Error()})
		return
	}
	defer file.Close()

	resp, err := c.uc.Upload(file, header, authHeader)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
