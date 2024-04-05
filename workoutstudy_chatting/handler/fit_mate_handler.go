package handler

import (
	"net/http"
	"workoutstudy_chatting/service" // 서비스 패키지 경로에 맞게 수정

	"github.com/gin-gonic/gin"
)

type fitMateHandler struct {
	Service *service.FitMateService
}

func NewFitMateHandler(s *service.FitMateService) *fitMateHandler {
	return &fitMateHandler{Service: s}
}

// RetrieveFitGroupByMateID - fit_mate_id로 fit_group 조회
func (ctrl *fitMateHandler) RetrieveFitGroupByMateID(c *gin.Context) {
	fitMateID := c.Param("fitMateId") // URL에서 fit_mate_id 추출

	// 서비스를 통해 fit_group 조회
	fitGroup, err := ctrl.Service.GetFitGroupByMateID(fitMateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, fitGroup)
}
