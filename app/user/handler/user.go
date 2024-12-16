package user

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"washit-api/app/user/dto/request"
	"washit-api/app/user/dto/resource"
	"washit-api/app/user/service"
	"washit-api/utils"
)

type UserHandler struct {
	service service.UserServiceInterface
}

func NewUserHandler(service service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Login(c *gin.Context) {
	var r request.Login
	if err := utils.ParseJson(c, &r); err != nil {
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(r); err != nil {
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(c, &r)
	if err != nil {
		log.Println("Failed to login as user ", err)
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	c.Set("userId", user.ID)
	c.Next()

	var res resource.User
	utils.CopyTo(user, &res)
	utils.WriteJson(c, http.StatusOK, map[string]interface{}{"user": res, "accessToken": accessToken, "refreshToken": refreshToken})
}

func (h *UserHandler) Register(c *gin.Context) {
	var r request.Register
	if err := utils.ParseJson(c, &r); err != nil {
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(r); err != nil {
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.Register(c, &r)
	if err != nil {
		log.Println("Failed to register user ", err)
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	var res resource.User
	utils.CopyTo(user, &res)
	utils.WriteJson(c, http.StatusCreated, map[string]interface{}{"user": res})
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("userId")
	println("Handler: ", userID)
	if userID == "" {
		utils.WriteError(c, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	user, err := h.service.GetUserByID(c, userID)
	if err != nil {
		log.Println("Failed to get user ", err)
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	var res resource.User
	utils.CopyTo(&user, &res)
	utils.WriteJson(c, http.StatusOK, map[string]interface{}{"user": res})
}
