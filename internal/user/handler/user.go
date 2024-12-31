package user

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	fireBase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"

	userRequest "washit-api/internal/user/dto/request"
	userResource "washit-api/internal/user/dto/resource"
	userService "washit-api/internal/user/service"
	"washit-api/pkg/configs"
	"washit-api/pkg/redis"
	jwt "washit-api/pkg/token"
	"washit-api/pkg/utils"
)

type UserHandler struct {
	service   userService.UserServiceInterface
	cache     redis.RedisInterface
	app       *fireBase.App
	validator *validator.Validate
}

func NewUserHandler(service userService.UserServiceInterface, cache redis.RedisInterface, app *fireBase.App, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		service:   service,
		cache:     cache,
		app:       app,
		validator: validator,
	}
}

var MeCacheKey = "/api/v1/profile/me"

func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		log.Println("Failed to get userId from context")
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get userId from context", errors.New("Failed to get userId from context"))
		return
	}

	accessToken, err := h.service.RefreshToken(c, userID)
	if err != nil {
		log.Println("Failed to refresh token ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to refresh token", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Successfully refreshed token", gin.H{"accessToken": accessToken}, nil)
}

// @Summary	Login with Google
// @Tags		User
// @Accept		json
// @Produce	json
// @Param		_	body		userRequest.Google	true	"Body"
// @Success	200	{object}	userResource.User
// @Router		/auth/login/google [post]
func (h *UserHandler) LoginWithGoogle(c *gin.Context) {
	var res userResource.WithToken
	var req userRequest.Google

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	client, err := h.app.Auth(context.Background())
	if err != nil {
		log.Println("Failed to initialize fireBase Auth ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to initialize fireBase Auth", err)
		return
	}

	token, err := client.VerifyIDToken(context.Background(), req.IDToken)
	if err != nil {
		log.Println("Failed to verify ID token ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to verify ID token", err)
		return
	}

	userRecord, err := client.GetUser(context.Background(), token.UID)
	if err != nil {
		log.Println("Failed to get user record ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get user record", err)
		return
	}

	user, accessToken, refreshToken, err := h.service.LoginWithGoogle(c, &req, userRecord.ProviderUserInfo[0])
	if err != nil {
		log.Println("Failed to login with Google ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to login with Google", err)
		return
	}

	utils.CopyTo(&user, &res.User)
	utils.CopyTo(&accessToken, &res.AccessToken)
	utils.CopyTo(&refreshToken, &res.RefreshToken)
	utils.SuccessResponse(c, http.StatusOK, "Successfully logged in with Google", &res, nil)
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
	var res userResource.WithToken

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		log.Println("Failed to validate request ", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to validate request", err)
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(c, &req)
	if err != nil {
		log.Println("Failed to login as user ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	tokenString, ok := accessToken.(string)
	if !ok {
		log.Println("Failed to assert accessToken as string")
		utils.WriteError(c, http.StatusInternalServerError, errors.New("Failed to assert accessToken as string"))
		return
	}

	c.SetCookie("jwt", tokenString, jwt.AccessTokenExpiredTime, "/", c.Request.Host, false, true)

	utils.CopyTo(&user, &res.User)
	utils.CopyTo(&accessToken, &res.AccessToken)
	utils.CopyTo(&refreshToken, &res.RefreshToken)
	utils.SuccessResponse(c, http.StatusOK, "Successfully logged in", &res, nil)
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
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		log.Println("Failed to validate request ", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to validate request", err)
		return
	}

	user, err := h.service.Register(c, &req)
	if err != nil {
		log.Println("Failed to register user ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	utils.CopyTo(&user, &res)
	utils.SuccessResponse(c, http.StatusCreated, "Successfully registered", &res, nil)
}

func (h *UserHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
	utils.SuccessResponse(c, http.StatusOK, "Successfully logged out", nil, nil)
}

func (h *UserHandler) BanUser(c *gin.Context) {
	var res userResource.User

	user, err := h.service.BanUser(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to ban user ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to ban user", err)
		return
	}

	utils.CopyTo(&user, &res)
	utils.SuccessResponse(c, http.StatusOK, user.FirstName+" is successfully banned", &res, links(res.ID))
}

func (h *UserHandler) UnbanUser(c *gin.Context) {
	var res userResource.User

	user, err := h.service.UnbanUser(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to unban user ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to unban user", err)
		return
	}

	utils.CopyTo(&user, &res)
	utils.SuccessResponse(c, http.StatusOK, user.FirstName+" is successfully unbanned", &res, links(res.ID))
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
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		log.Println("Failed to validate request ", err)
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to validate request", err)
		return
	}

	userID := c.GetString("userId")
	user, err := h.service.UpdateMe(c, userID, &req)
	if err != nil {
		log.Println("Failed to update user ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	_ = h.cache.Remove(MeCacheKey)

	utils.CopyTo(&user, &res)
	utils.SuccessResponse(c, http.StatusOK, "Successfully updated", &res, links(res.ID))
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

	if err := h.cache.Get(MeCacheKey, &res); err == nil {
		utils.SuccessResponse(c, http.StatusOK, "Successfully retrieved user", &res, links(res.ID))
		return
	}

	userID := c.GetString("userId")
	user, err := h.service.GetMe(c, userID)
	if err != nil {
		log.Println("Failed to get user ", err)
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	utils.CopyTo(&user, &res)
	utils.SuccessResponse(c, http.StatusOK, "Successfully retrieved user", &res, links(res.ID))

	_ = h.cache.SetWithExpiration(MeCacheKey, &res, configs.ProductCachingTime)
}

// @Summary	Get all users
// @Tags		User
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	userResource.UserList
// @Router		/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	var res []userResource.User

	users, err := h.service.GetUsers(c)
	if err != nil {
		log.Println("Failed to get users ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get users", err)
		return
	}

	utils.CopyTo(&users, &res)
	utils.SuccessResponse(c, http.StatusOK, "Successfully retrieved users", &res, nil)
}

// @Summary	Get all banned users
// @Tags		User
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Success	200	{object}	userResource.UserList
// @Router		/users/banned [get]
func (h *UserHandler) GetBannedUsers(c *gin.Context) {
	var res []userResource.User

	users, err := h.service.GetBannedUsers(c)
	if err != nil {
		log.Println("Failed to get users ", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get banned users", err)
		return
	}

	utils.CopyTo(&users, &res)
	utils.SuccessResponse(c, http.StatusOK, "Successfully retrieved banned users", res, nil)
}

// @Summary	Get a user by ID
// @Tags		User
// @Accept		json
// @Produce	json
// @Security	ApiKeyAuth
// @Param		id	path		string	true	"User ID"
// @Success	200	{object}	userResource.User
// @Router		/user/{id} [get]
func (h *UserHandler) GetUserById(c *gin.Context) {
	var res userResource.User

	user, err := h.service.GetUserByID(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to get user ", err)
		utils.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	utils.CopyTo(&user, &res)
	utils.SuccessResponse(c, http.StatusOK, "Successfully retrieved user", &res, nil)
}

var links = func(orderId int64) map[string]utils.HypermediaLink {
	return map[string]utils.HypermediaLink{
		"self": {
			Href:   "/profile/me",
			Method: "GET",
		},
		"self-alternative": {
			Href:   "/user/" + strconv.FormatInt(orderId, 10),
			Method: "GET",
		},
		"update": {
			Href:   "/profile/update",
			Method: "PUT",
		},
		"ban": {
			Href:   "/user/" + strconv.FormatInt(orderId, 10) + "/ban",
			Method: "PUT",
		},
		"unban": {
			Href:   "/user/" + strconv.FormatInt(orderId, 10) + "/unban",
			Method: "PUT",
		},
	}
}
