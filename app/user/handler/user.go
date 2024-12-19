package user

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	userRequest "washit-api/app/user/dto/request"
	userResource "washit-api/app/user/dto/resource"
	userService "washit-api/app/user/service"
	"washit-api/configs"
	"washit-api/redis"
	jwt "washit-api/token"
	"washit-api/utils"
)

type UserHandler struct {
	service userService.UserServiceInterface
	cache   redis.RedisInterface
}

func NewUserHandler(service userService.UserServiceInterface, cache redis.RedisInterface) *UserHandler {
	return &UserHandler{
		service: service,
		cache:   cache,
	}
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	userID := ctx.GetString("userId")
	if userID == "" {
		utils.WriteError(ctx, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	accessToken, err := h.service.RefreshToken(ctx, userID)
	if err != nil {
		log.Println("Failed to refresh token ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(ctx, http.StatusOK, map[string]interface{}{"accessToken": accessToken})
}

// @Summary	Login as a user
// @Tags		User
// @Accept		json
// @Produce	json
// @Param		_	body		userRequest.Login	true	"Body"
// @Success	200	{object}	userResource.ShowToken
// @Router		/auth/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	var req userRequest.Login
	var res userResource.ShowToken

	if err := utils.ParseJson(ctx, &req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(&req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(ctx, &req)
	if err != nil {
		log.Println("Failed to login as user ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	tokenString, ok := accessToken.(string)
	if !ok {
		log.Println("Failed to assert accessToken as string")
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.SetCookie("jwt", tokenString, jwt.AccessTokenExpiredTime, "/", "localhost", false, true)

	utils.CopyTo(&user, &res.User)
	utils.CopyTo(&accessToken, &res.AccessToken)
	utils.CopyTo(&refreshToken, &res.RefreshToken)
	utils.WriteJson(ctx, http.StatusOK, &res)
}

// @Summary	Register a new user
// @Tags		User
// @Accept		json
// @Produce	json
// @Param		_	body		userRequest.Register	true	"Body"
// @Success	201	{object}	userResource.HideToken
// @Router		/auth/register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	var req userRequest.Register
	var res userResource.HideToken

	if err := utils.ParseJson(ctx, &req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(&req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.Register(ctx, &req)
	if err != nil {
		log.Println("Failed to register user ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&user, &res.User)
	utils.WriteJson(ctx, http.StatusCreated, &res)
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	ctx.SetCookie("jwt", "", -1, "/", "localhost", false, true)
	utils.WriteJson(ctx, http.StatusOK, map[string]interface{}{"message": "Successfully logged out"})
}

// @Summary	Update the current logged-in user
// @Tags		User
// @Accept		json
// @Produce	json
// @Param		_	body		userRequest.Update	true	"Body"
// @Success	201	{object}	userResource.HideToken
// @Router		/profile/update [post]
func (h *UserHandler) UpdateMe(ctx *gin.Context) {
	var req userRequest.Update
	var res userResource.HideToken

	if err := utils.ParseJson(ctx, &req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(&req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	userID := ctx.GetString("userId")

	user, err := h.service.UpdateMe(ctx, userID, &req)
	if err != nil {
		log.Println("Failed to update user ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&user, &res.User)
	utils.WriteJson(ctx, http.StatusOK, &res)
}

// @Summary	Get the current logged-in user
// @Tags		User
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	userResource.HideToken
// @Router		/profile/me [get]
func (h *UserHandler) GetMe(ctx *gin.Context) {
	var res userResource.HideToken

	cacheKey := ctx.Request.URL.RequestURI()
	err := h.cache.Get(cacheKey, &res)
	if err == nil {
		utils.WriteJson(ctx, http.StatusOK, &res)
		return
	}

	userID := ctx.GetString("userId")
	user, err := h.service.GetMe(ctx, userID)
	if err != nil {
		log.Println("Failed to get user ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&user, &res.User)
	utils.WriteJson(ctx, http.StatusOK, &res)
	_ = h.cache.SetWithExpiration(cacheKey, &res, configs.ProductCachingTime)
}

// @Summary	Get all users
// @Tags		User
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	userResource.Base
// @Router		/users [get]
func (h *UserHandler) GetUsers(ctx *gin.Context) {
	var res []userResource.Base

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		log.Println("Failed to get users ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&users, &res)
	utils.WriteJson(ctx, http.StatusOK, utils.ToData("users", &res))
}

// @Summary	Get a user by ID
// @Tags		User
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	userResource.Base
// @Router		/user/:id [get]
func (h *UserHandler) GetUserById(ctx *gin.Context) {
	var res userResource.Base

	user, err := h.service.GetUserByID(ctx, ctx.Param("id"))
	if err != nil {
		log.Println("Failed to get user ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&user, &res)
	utils.WriteJson(ctx, http.StatusOK, utils.ToData("user", &res))
}
