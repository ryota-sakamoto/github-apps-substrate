package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v30/github"

	"github.com/ryota-sakamoto/github-apps-substrate/pkg/service"
)

type CallbackHandler struct {
	subscribeService service.SubscribeService
}

func NewCallbackHandler(ss service.SubscribeService) CallbackHandler {
	return CallbackHandler{
		subscribeService: ss,
	}
}

func (ca CallbackHandler) Endpoint(c *gin.RouterGroup) {
	c.POST("/callback", ca.callback)
}

func (ca CallbackHandler) callback(c *gin.Context) {
	var err error

	switch c.GetHeader("X-Github-Event") {
	case "push":
		var event github.PushEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			log.Printf("%+v\n", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		err = ca.subscribeService.SubscribePush(&event)
	default:
		c.AbortWithStatus(http.StatusNotImplemented)
		return
	}

	if err != nil {
		log.Printf("%+v\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
