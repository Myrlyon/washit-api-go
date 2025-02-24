package user

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	fireBase "firebase.google.com/go"
	"github.com/gin-gonic/gin"

	userRequest "washit-api/internal/user/dto/request"
	userResource "washit-api/internal/user/dto/resource"
	userService "washit-api/internal/user/service"
	"washit-api/pkg/configs"
	"washit-api/pkg/redis"
	"washit-api/pkg/response"
	jwt "washit-api/pkg/token"
	"washit-api/pkg/utils"
)

type UserHandler struct {
	service userService.IUserService
	cache   redis.IRedis
	app     *fireBase.App
}

func NewUserHandler(service userService.IUserService, cache redis.IRedis, app *fireBase.App) *UserHandler {
	return &UserHandler{
		service: service,
		cache:   cache,
		app:     app,
	}
}

var MeCacheKey = "/api/v1/profile/me"

// RefreshToken refreshes the user's access token
//
//	@Summary	Refresh the user's access token
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	map[string]string	"accessToken"
//	@Router		/auth/refresh-token [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		log.Println("Failed to get userID from context")
		response.Error(c, http.StatusInternalServerError, "Failed to get userID from context", errors.New("Failed to get userID from context"))
		return
	}

	accessToken, err := h.service.RefreshToken(c, userID)
	if err != nil {
		log.Println("Failed to refresh token ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to refresh token", err)
		return
	}

	response.Success(c, http.StatusOK, "Successfully refreshed token", gin.H{"accessToken": accessToken}, nil)
}

// LoginWithGoogle handles user login via Google
//
//	@Summary	Login with Google
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		_	body		userRequest.Google	true	"Body"
//	@Success	200	{object}	userResource.WithToken
//	@Router		/auth/login/google [post]
func (h *UserHandler) LoginWithGoogle(c *gin.Context) {
	var res userResource.WithToken
	var req userRequest.Google

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		response.Error(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	client, err := h.app.Auth(context.Background())
	if err != nil {
		log.Println("Failed to initialize fireBase Auth ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to initialize fireBase Auth", err)
		return
	}

	token, err := client.VerifyIDToken(context.Background(), req.TokenID)
	if err != nil {
		log.Println("Failed to verify token ID ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to verify token ID ", err)
		return
	}

	userRecord, err := client.GetUser(context.Background(), token.UID)
	if err != nil {
		log.Println("Failed to get user record ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to get user record", err)
		return
	}

	user, accessToken, refreshToken, err := h.service.LoginWithGoogle(c, &req, userRecord.ProviderUserInfo[0])
	if err != nil {
		log.Println("Failed to login with Google ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to login with Google", err)
		return
	}

	utils.CopyTo(&user, &res.User)
	utils.CopyTo(&accessToken, &res.AccessToken)
	utils.CopyTo(&refreshToken, &res.RefreshToken)
	response.Success(c, http.StatusOK, "Successfully logged in with Google", &res, nil)
}

// Login handles user login
//
//	@Summary	Login as a user
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		_	body		userRequest.Login	true	"Body"
//	@Success	200	{object}	userResource.WithToken
//	@Router		/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req userRequest.Login
	var res userResource.WithToken

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		response.Error(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	user, accessToken, refreshToken, err := h.service.Login(c, &req)
	if err != nil {
		log.Println("Failed to login as user ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	c.SetCookie("jwt", accessToken, jwt.AccessTokenExpiredTime, "/", c.Request.Host, false, true)

	utils.CopyTo(&user, &res.User)
	utils.CopyTo(&accessToken, &res.AccessToken)
	utils.CopyTo(&refreshToken, &res.RefreshToken)
	response.Success(c, http.StatusOK, "Successfully logged in", &res, nil)
}

// Register handles user registration
//
//	@Summary	Register a new user
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Param		_	body		userRequest.Register	true	"Body"
//	@Success	201	{object}	userResource.User
//	@Router		/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req userRequest.Register
	var res userResource.User

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		response.Error(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	user, err := h.service.Register(c, &req)
	if err != nil {
		log.Println("Failed to register user ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to register user", err)
		return
	}

	utils.CopyTo(&user, &res)
	response.Success(c, http.StatusCreated, "Successfully registered", &res, nil)
}

// Logout handles user logout
//
//	@Summary	Logout the current logged-in user
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Success	200	{object}	userResource.User
//	@Router		/auth/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true)
	response.Success(c, http.StatusOK, "Successfully logged out", nil, nil)
}

// BanUser bans a user by ID
//
//	@Summary	Ban a user
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	userResource.User
//	@Router		/user/{id}/ban [put]
func (h *UserHandler) BanUser(c *gin.Context) {
	var res userResource.User

	user, err := h.service.BanUser(c, c.GetString("userID"))
	if err != nil {
		log.Println("Failed to ban user ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to ban user", err)
		return
	}

	utils.CopyTo(&user, &res)
	response.Success(c, http.StatusOK, user.FirstName+" is successfully banned", &res, links(res.ID))
}

// UnbanUser unbans a user by ID
//
//	@Summary	Unban a user
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	userResource.User
//	@Router		/user/{id}/unban [put]
func (h *UserHandler) UnbanUser(c *gin.Context) {
	var res userResource.User

	user, err := h.service.UnbanUser(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to unban user ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to unban user", err)
		return
	}

	utils.CopyTo(&user, &res)
	response.Success(c, http.StatusOK, user.FirstName+" is successfully unbanned", &res, links(res.ID))
}

// UpdateMe updates the current logged-in user's profile
//
//	@Summary	Update the current logged-in user
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body		userRequest.UpdateProfile	true	"Body"
//	@Success	201	{object}	userResource.User
//	@Router		/profile/update [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	var req userRequest.UpdateProfile
	var res userResource.User

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		response.Error(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	user, err := h.service.UpdateProfile(c, c.GetString("userID"), &req)
	if err != nil {
		log.Println("Failed to update user ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	_ = h.cache.Remove(MeCacheKey)

	utils.CopyTo(&user, &res)
	response.Success(c, http.StatusOK, "Successfully updated", &res, links(res.ID))
}

// UpdatePassword updates the current logged-in user's password
//
//	@Summary	Update the current logged-in user's password
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body		userRequest.UpdatePassword	true	"Body"
//	@Success	201	{object}	userResource.User
//	@Router		/profile/update/password [put]
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var req userRequest.UpdatePassword

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		response.Error(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	if err := h.service.UpdatePassword(c, c.GetString("userID"), &req); err != nil {
		log.Println("Failed to update password ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to update password", err)
		return
	}

	response.Success(c, http.StatusOK, "Successfully updated password", nil, nil)
}

// UpdatePicture updates the current logged-in user's profile picture
//
//	@Summary	Update the current logged-in user's profile picture
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		_	body		userRequest.UpdatePicture	true	"Body"
//	@Success	200	{object}	userResource.User
//	@Router		/profile/update/picture [put]
func (h *UserHandler) UpdatePicture(c *gin.Context) {
	var req userRequest.UpdatePicture

	if err := utils.ParseJson(c, &req); err != nil {
		log.Println("Failed to parse request ", err)
		response.Error(c, http.StatusBadRequest, "Failed to parse request", err)
		return
	}

	user, err := h.service.UpdatePicture(c, c.GetString("userID"), &req)
	if err != nil {
		log.Println("Failed to update profile picture. err: ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to update profile picture", err)
	}

	response.Success(c, http.StatusOK, "Successfully updated profile picture", user, nil)
}

// GetMe retrieves the current logged-in user's profile
//
//	@Summary	Get the current logged-in user
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Success	200	{object}	userResource.User
//	@Router		/profile/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	var res userResource.User

	if err := h.cache.Get(MeCacheKey, &res); err == nil {
		response.Success(c, http.StatusOK, "Successfully retrieved user", &res, links(res.ID))
		return
	}

	user, err := h.service.GetMe(c, c.GetString("userID"))
	if err != nil {
		log.Println("Failed to get user ", err)
		response.Error(c, http.StatusNotFound, "User not found", err)
		return
	}

	utils.CopyTo(&user, &res)
	response.Success(c, http.StatusOK, "Successfully retrieved user", &res, links(res.ID))

	_ = h.cache.SetWithExpiration(MeCacheKey, &res, configs.ProductCachingTime)
}

// GetUsers retrieves all users
//
//	@Summary	Get all users
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Success	200	{object}	userResource.User
//	@Router		/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	var res []userResource.User

	users, err := h.service.GetUsers(c)
	if err != nil {
		log.Println("Failed to get users ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to get users", err)
		return
	}

	utils.CopyTo(&users, &res)
	response.Success(c, http.StatusOK, "Successfully retrieved users", &res, nil)
}

// GetBannedUsers retrieves all banned users
//
//	@Summary	Get all banned users
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Success	200	{object}	userResource.User
//	@Router		/users/banned [get]
func (h *UserHandler) GetBannedUsers(c *gin.Context) {
	var res []userResource.User

	users, err := h.service.GetBannedUsers(c)
	if err != nil {
		log.Println("Failed to get users ", err)
		response.Error(c, http.StatusInternalServerError, "Failed to get banned users", err)
		return
	}

	utils.CopyTo(&users, &res)
	response.Success(c, http.StatusOK, "Successfully retrieved banned users", res, nil)
}

// GetUserByID retrieves a user by ID
//
//	@Summary	Get a user by ID
//	@Tags		User
//	@Accept		json
//	@Produce	json
//	@Security	ApiKeyAuth
//	@Param		id	path		string	true	"User ID"
//	@Success	200	{object}	userResource.User
//	@Router		/user/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	var res userResource.User

	user, err := h.service.GetUserByID(c, c.Param("id"))
	if err != nil {
		log.Println("Failed to get user ", err)
		response.Error(c, http.StatusNotFound, "User not found", err)
		return
	}

	utils.CopyTo(&user, &res)
	response.Success(c, http.StatusOK, "Successfully retrieved user", &res, nil)
}

var links = func(orderID int64) map[string]response.HypermediaLink {
	return map[string]response.HypermediaLink{
		"self": {
			Href:   "/profile/me",
			Method: "GET",
		},
		"self-alternative": {
			Href:   "/user/" + strconv.FormatInt(orderID, 10),
			Method: "GET",
		},
		"update": {
			Href:   "/profile/update",
			Method: "PUT",
		},
		"ban": {
			Href:   "/user/" + strconv.FormatInt(orderID, 10) + "/ban",
			Method: "PUT",
		},
		"unban": {
			Href:   "/user/" + strconv.FormatInt(orderID, 10) + "/unban",
			Method: "PUT",
		},
	}
}
