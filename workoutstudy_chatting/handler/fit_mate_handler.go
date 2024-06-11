package handler

import (
	"net/http"
	"strconv"
	"workoutstudy_chatting/service" // 서비스 패키지 경로에 맞게 수정

	"github.com/gin-gonic/gin"
)

type fitMateHandler struct {
	FitmateService service.FitMateUseCase
}

func NewFitMateHandler(fitmateService service.FitMateUseCase) *fitMateHandler {
	return &fitMateHandler{
		FitmateService: fitmateService,
	}
}

// @Summary 피트그룹 조회 API
// @Description userId 로 해당 사용자가 속해 있는 피트그룹들의 정보를 조희
// @Tags fitmate
// @Accept  json
// @Produce  json
// @Param userId query int true  "사용자 ID, fitMateId 가 아님.""
// @Success 200 {object} model.FitGroup
// @Router /retrieve/fit-group [get]
func (h *fitMateHandler) RetrieveFitGroupByUserID(c *gin.Context) {
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
