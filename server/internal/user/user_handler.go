package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/reyhanyogs/realtime-chat/domain"
)

type Handler struct {
	domain.Service
}

func NewHandler(r *gin.Engine, s domain.Service) {
	handler := &Handler{
		Service: s,
	}
	r.POST("/signup", handler.CreateUser)
	r.POST("/login", handler.Login)
	r.GET("/logout", handler.Logout)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var u domain.CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.Service.CreateUser(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Login(c *gin.Context) {
	var user domain.LoginUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	u, err := h.Service.Login(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("jwt", u.AccessToken, 3600, "/", "localhost", false, true)

	res := &domain.LoginUserRes{
		Username: u.Username,
		ID:       u.ID,
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "logout successful",
	})
}
