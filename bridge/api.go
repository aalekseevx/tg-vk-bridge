package bridge

import (
	"errors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/toorop/gin-logrus"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
)

func setDbErrorToContext(err error, context *gin.Context) {
	var status int
	if errors.Is(err, gorm.ErrRecordNotFound) {
		status = http.StatusNotFound
	} else {
		status = http.StatusInternalServerError
	}
	context.JSON(status, gin.H{
		"message": err,
	})
}

func setRedisErrorToContext(err error, context *gin.Context) {
	var status int
	if err == redis.Nil {
		status = http.StatusNotFound
	} else {
		status = http.StatusInternalServerError
	}
	context.JSON(status, gin.H{
		"message": err,
	})
}

func rootEndpoint(context *gin.Context) {
	context.Redirect(http.StatusMovedPermanently, "https://t.me/"+os.Getenv("BOT_USERNAME"))
}

func matchEndpoint(context *gin.Context) {
	// Bind JSON request
	var request struct {
		Channel       int64  `json:"channel" binding:"required"`
		Group         int    `json:"group" binding:"required"`
		Delay         *int   `json:"delay" binding:"required"`
		VkAccessToken string `json:"vk_access_token" binding:"required"`
		CreatorToken  string `json:"creator_token" binding:"required"`
	}
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// Verify that request comes from channel creator
	isOwner, err := isOwnerToken(request.CreatorToken, request.Channel)
	if err != nil {
		setDbErrorToContext(err, context)
		return
	}
	if !isOwner {
		context.JSON(http.StatusForbidden, gin.H{
			"message": "You are not a creator of this channel",
		})
		return
	}

	// Match channel to VK group
	err = matchChannel(Channel{
		ChatID:        request.Channel,
		GroupID:       request.Group,
		Delay:         *request.Delay,
		VkAccessToken: request.VkAccessToken,
	})
	if err != nil {
		setDbErrorToContext(err, context)
		return
	}

	context.JSON(http.StatusOK, gin.H{})
}

func channelsEndpoint(context *gin.Context) {
	// Bind JSON request
	var request struct {
		CreatorToken string `form:"token" binding:"required"`
	}
	if err := context.ShouldBindQuery(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	// Find channels by token
	channels, err := getChannelsByCreatorToken(request.CreatorToken)
	if err != nil {
		setDbErrorToContext(err, context)
		return
	}

	context.JSON(http.StatusOK, channels)
}

func configEndpoint(context *gin.Context) {
	vkAppId, err := strconv.Atoi(os.Getenv("VK_APP_ID"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"bot_username": os.Getenv("BOT_USERNAME"),
		"vk_app_id": vkAppId,
	})
}

func postEndpoint(context *gin.Context) {
	var request struct {
		PostToken string `form:"token" binding:"required"`
	}
	if err := context.ShouldBindQuery(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	post, err := getPostByToken(request.PostToken)
	if err != nil {
		setRedisErrorToContext(err, context)
		return
	}
	context.JSON(http.StatusOK, post)
}

func finaliseEndpoint(context *gin.Context) {
	var request struct {
		PostToken string `form:"token" binding:"required"`
		VkResponse string `form:"vk_response" binding:"required"`
	}
	err := context.ShouldBindQuery(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	post, err := getPostByToken(request.PostToken)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{
			"message": err,
		})
		return
	}

	err = sendFinaliseNotification(post, request.VkResponse)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{})
}

func mainApi() {
	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Panic(err)
	}
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DefaultWriter = log.WithField("service", "api").Writer()

	router := gin.New()
	router.Use(ginlogrus.Logger(log), gin.Recovery())
	router.Use(static.Serve("/", static.LocalFile("./static", false)))

	router.GET("/", rootEndpoint)
	router.GET("/api/config", configEndpoint)
	router.GET("/api/post", postEndpoint)
	router.GET("/api/channels", channelsEndpoint)
	router.POST("/api/match", matchEndpoint)
	router.POST("/api/finalise", finaliseEndpoint)


	if err = router.Run(); err != nil {
		panic(err)
	}
}
