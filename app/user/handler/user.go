package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	userRequest "washit-api/app/user/dto/request"
	userResource "washit-api/app/user/dto/resource"
	userService "washit-api/app/user/service"
	"washit-api/configs"
	"washit-api/redis"
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

func (h *UserHandler) Login(ctx *gin.Context) {
	var req userRequest.Login
	var res userResource.User

	if err := utils.ParseJson(ctx, &req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	user, accessToken, err := h.service.Login(ctx, &req)
	if err != nil {
		log.Println("Failed to login as user ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(user, &res)
	utils.WriteJson(ctx, http.StatusOK, map[string]interface{}{"user": res, "accessToken": accessToken})
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req userRequest.Register
	var res userResource.User

	if err := utils.ParseJson(ctx, &req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.WriteError(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.Register(ctx, &req)
	if err != nil {
		log.Println("Failed to register user ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(user, &res)
	utils.WriteJson(ctx, http.StatusCreated, map[string]interface{}{"user": res})
}

func (h *UserHandler) GetMe(ctx *gin.Context) {
	var res userResource.User

	cacheKey := ctx.Request.URL.RequestURI()
	err := h.cache.Get(cacheKey, &res)
	if err == nil {
		utils.WriteJson(ctx, http.StatusOK, map[string]interface{}{"user": res})
		return
	}

	userID := ctx.GetString("userId")

	user, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		log.Println("Failed to get user ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&user, &res)
	utils.WriteJson(ctx, http.StatusOK, map[string]interface{}{"user": res})
	_ = h.cache.SetWithExpiration(cacheKey, res, configs.ProductCachingTime)
}

func (h *UserHandler) GetUsers(ctx *gin.Context) {
	var res []userResource.User

	users, err := h.service.GetUsers(ctx)
	if err != nil {
		log.Println("Failed to get users ", err)
		utils.WriteError(ctx, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&users, &res)
	utils.WriteJson(ctx, http.StatusOK, map[string]interface{}{"data": res})
}
