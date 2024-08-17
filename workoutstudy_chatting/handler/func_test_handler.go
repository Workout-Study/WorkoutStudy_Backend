package handler

import (
	"log"
	"net/http"
	"time"
	"workoutstudy_chatting/model"
	"workoutstudy_chatting/service"

	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	UserService     service.UserUseCase
	FitGroupService service.FitGroupUseCase
	FitMateService  service.FitMateUseCase
}

// NewTestHandler는 TestHandler의 새로운 인스턴스를 반환합니다.
func NewTestHandler(userService service.UserUseCase, fitGroupService service.FitGroupUseCase, fitMateService service.FitMateUseCase) *TestHandler {
	return &TestHandler{
		UserService:     userService,
		FitGroupService: fitGroupService,
		FitMateService:  fitMateService,
	}
}

// POST /test/create/user
func (h *TestHandler) CreateUser(c *gin.Context) {
	var user model.Users
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	savedUser, err := h.UserService.SaveUser(&user)
	if err != nil {
		log.Printf("Error saving user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusCreated, savedUser)
}

// POST /test/create/fit-group
func (h *TestHandler) CreateFitGroup(c *gin.Context) {
	var fitGroup model.FitGroup
	if err := c.BindJSON(&fitGroup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fitGroup.CreatedAt = time.Now()
	fitGroup.UpdatedAt = time.Now()

	savedFitGroup, err := h.FitGroupService.SaveFitGroup(&fitGroup)
	if err != nil {
		log.Printf("Error saving fit group: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save fit group"})
		return
	}

	c.JSON(http.StatusCreated, savedFitGroup)
}

// POST /test/create/fit-mate
func (h *TestHandler) CreateFitMate(c *gin.Context) {
	var fitMate model.FitMate
	if err := c.BindJSON(&fitMate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	fitMate.CreatedAt = time.Now()
	fitMate.UpdatedAt = time.Now()

	savedFitMate, err := h.FitMateService.SaveFitMate(&fitMate)
	if err != nil {
		log.Printf("Error saving fit mate: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save fit mate"})
		return
	}

	c.JSON(http.StatusCreated, savedFitMate)
}
