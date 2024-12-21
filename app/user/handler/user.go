package user

import (
	"context"
	"errors"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"

	userRequest "washit-api/app/user/dto/request"
	userResource "washit-api/app/user/dto/resource"
	userService "washit-api/app/user/service"
	"washit-api/configs"
	"washit-api/redis"
	jwt "washit-api/token"
	"washit-api/utils"
)

type UserHandler struct {
	service   userService.UserServiceInterface
	cache     redis.RedisInterface
	app       *firebase.App
	validator *validator.Validate
}

func NewUserHandler(service userService.UserServiceInterface, cache redis.RedisInterface, app *firebase.App, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		service:   service,
		cache:     cache,
		app:       app,
		validator: validator,
	}
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		log.Println("Failed to get userId from context")
		utils.WriteError(c, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	accessToken, err := h.service.RefreshToken(c, userID)
	if err != nil {
		log.Println("Failed to refresh token ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(c, http.StatusOK, gin.H{"accessToken": accessToken})
}

// @Summary	Login with Google
// @Tags		User
// @Accept		json
// @Produce	json
// @Param		_	body		userRequest.Google	true	"Body"
// @Success	200	{object}	userResource.User
// @Router		/auth/login/google [post]
func (h *UserHandler) LoginWithGoogle(c *gin.Context) {
	var res userResource.User
	var req userRequest.Google

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	client, err := h.app.Auth(context.Background())
	if err != nil {
		log.Println("Failed to initialize Firebase Auth ", err)
		utils.WriteError(c, http.StatusInternalServerError, errors.New("Failed to initialize Firebase Auth"))
		return
	}

	token, err := client.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		log.Println("Failed to verify ID token ", err)
		utils.WriteError(c, http.StatusUnauthorized, errors.New("Invalid ID token"))
		return
	}

	userRecord, err := client.GetUser(context.Background(), token.UID)
	if err != nil {
		log.Println("Failed to get user record ", err)
		utils.WriteError(c, http.StatusInternalServerError, errors.New("Failed to retrieve user profile"))
		return
	}

	user, accessToken, refreshToken, err := h.service.LoginWithGoogle(c, &req, userRecord.ProviderUserInfo[0])
	if err != nil {
		log.Println("Failed to login with Google ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&user, &res.User)
	utils.CopyTo(&accessToken, &res.AccessToken)
	utils.CopyTo(&refreshToken, &res.RefreshToken)
	res.Message = "Successfully logged in with Google"
	utils.WriteJson(c, http.StatusOK, &res)
}

// @Summary	Login as a user
// @Tags		User
// @Accept		json
// @Produce	json
// @Param		_	body		userRequest.Login	true	"Body"
// @Success	200	{object}	userResource.User
// @Router		/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req userRequest.Login
	var res userResource.User

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		log.Println("Failed to validate request ", err)
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(c, &req)
	if err != nil {
		log.Println("Failed to login as user ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	tokenString, ok := accessToken.(string)
	if !ok {
		log.Println("Failed to assert accessToken as string")
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("jwt", tokenString, jwt.AccessTokenExpiredTime, "/", "localhost", false, true)

	utils.CopyTo(&user, &res.User)
	utils.CopyTo(&accessToken, &res.AccessToken)
	utils.CopyTo(&refreshToken, &res.RefreshToken)
	res.Message = "Successfully logged in"
	utils.WriteJson(c, http.StatusOK, &res)
}

// @Summary	Register a new user
// @Tags		User
// @Accept		json
// @Produce	json
// @Param		_	body		userRequest.Register	true	"Body"
// @Success	201	{object}	userResource.User
// @Router		/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req userRequest.Register
	var res userResource.User

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		log.Println("Failed to validate request ", err)
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.service.Register(c, &req)
	if err != nil {
		log.Println("Failed to register user ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&user, &res.User)
	res.Message = "Successfully registered"
	utils.WriteJson(c, http.StatusCreated, &res)
}

func (h *UserHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
	utils.WriteJson(c, http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func (h *UserHandler) BanUser(c *gin.Context) {
	var res userResource.User

	user, err := h.service.BanUser(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to ban user ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&user, &res.User)
	res.Message = user.FirstName + " is successfully banned"
	utils.WriteJson(c, http.StatusOK, &res)
}

func (h *UserHandler) UnbanUser(c *gin.Context) {
	var res userResource.User

	user, err := h.service.UnbanUser(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to unban user ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&user, &res.User)
	res.Message = user.FirstName + " is successfully unbanned"
	utils.WriteJson(c, http.StatusOK, &res)
}

// @Summary	Update the current logged-in user
// @Tags		User
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Param		_	body		userRequest.Update	true	"Body"
// @Success	201	{object}	userResource.User
// @Router		/profile/update [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	var req userRequest.Update
	var res userResource.User

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		log.Println("Failed to validate request ", err)
		utils.WriteError(c, http.StatusBadRequest, err)
		return
	}

	userID := c.GetString("userId")
	user, err := h.service.UpdateMe(c, userID, &req)
	if err != nil {
		log.Println("Failed to update user ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	cacheKey := "/api/v1/profile/me"
	_ = h.cache.Remove(cacheKey)

	utils.CopyTo(&user, &res.User)
	res.Message = "Successfully updated"
	utils.WriteJson(c, http.StatusOK, &res)
}

// @Summary	Get the current logged-in user
// @Tags		User
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	userResource.User
// @Router		/profile/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	var res userResource.User

	cacheKey := c.Request.URL.RequestURI()
	if err := h.cache.Get(cacheKey, &res); err == nil {
		utils.WriteJson(c, http.StatusOK, &res)
		return
	}

	userID := c.GetString("userId")
	user, err := h.service.GetMe(c, userID)
	if err != nil {
		log.Println("Failed to get user ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&user, &res.User)
	res.Message = "user are successfully retrieved"
	utils.WriteJson(c, http.StatusOK, &res)

	_ = h.cache.SetWithExpiration(cacheKey, &res, configs.ProductCachingTime)
}

// @Summary	Get all users
// @Tags		User
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	userResource.UserList
// @Router		/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	var res userResource.UserList

	users, err := h.service.GetUsers(c)
	if err != nil {
		log.Println("Failed to get users ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&users, &res.Users)
	res.Message = "Successfully retrieved users"
	utils.WriteJson(c, http.StatusOK, &res)
}

// @Summary	Get all banned users
// @Tags		User
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	userResource.UserList
// @Router		/users/banned [get]
func (h *UserHandler) GetBannedUsers(c *gin.Context) {
	var res userResource.UserList

	users, err := h.service.GetBannedUsers(c)
	if err != nil {
		log.Println("Failed to get users ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&users, &res.Users)
	res.Message = "Successfully retrieved banned users"
	utils.WriteJson(c, http.StatusOK, &res)
}

// @Summary	Get a user by ID
// @Tags		User
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Param		id	path		string	true	"User ID"
// @Success	200	{object}	userResource.Base
// @Router		/user/{id} [get]
func (h *UserHandler) GetUserById(c *gin.Context) {
	var res userResource.User

	user, err := h.service.GetUserByID(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to get user ", err)
		utils.WriteError(c, http.StatusInternalServerError, err)
		return
	}

	utils.CopyTo(&user, &res.User)
	res.Message = "Successfully retrieved user"
	utils.WriteJson(c, http.StatusOK, &res)
}
