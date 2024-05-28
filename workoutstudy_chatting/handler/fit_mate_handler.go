package handler

import (
	"net/http"
	"strconv"
	"workoutstudy_chatting/service" // 서비스 패키지 경로에 맞게 수정

	"github.com/gin-gonic/gin"
)

type fitMateHandler struct {
	FitmateService service.FitMateService
}

func NewFitMateHandler(fitmateService service.FitMateService) *fitMateHandler {
	return &fitMateHandler{
		FitmateService: fitmateService,
	}
}

func (h *fitMateHandler) RetrieveFitGroupByMateID(c *gin.Context) {
	userID := c.Query("userId")

	// userID string -> int 변환
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}

	fitGroup, err := h.FitmateService.GetFitGroupsByUserID(userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, fitGroup)
}
