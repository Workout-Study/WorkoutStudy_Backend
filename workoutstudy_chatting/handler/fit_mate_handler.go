package handler

import (
	"net/http"
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
	fitMateID := c.Query("fitMateId")

	fitGroup, err := h.FitmateService.GetFitGroupByMateID(fitMateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, fitGroup)
}
