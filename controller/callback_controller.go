package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v30/github"

	"github.com/ryota-sakamoto/github-apps-substrate/service"
)

type CallbackController struct {
	installationService service.InstallationService
	subscribeService    service.SubscribeService
}

func NewCallbackController(is service.InstallationService, ss service.SubscribeService) CallbackController {
	return CallbackController{
		installationService: is,
		subscribeService:    ss,
	}
}

func (ca CallbackController) Endpoint(c *gin.RouterGroup) {
	c.POST("/callback", ca.callback)
}

func (ca CallbackController) callback(c *gin.Context) {
	var err error

	switch c.GetHeader("X-Github-Event") {
	case "installation", "installation_repositories":
		var event github.InstallationEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			log.Printf("%+v\n", err)
			c.AbortWithStatus(400)
			return
		}

		err = ca.installationService.Action(&event)
	case "push":
		var event github.PushEvent
		if err := c.ShouldBindJSON(&event); err != nil {
			log.Printf("%+v\n", err)
			c.AbortWithStatus(400)
			return
		}

		err = ca.subscribeService.SubscribePush(&event)
	default:
		c.AbortWithStatus(400)
		return
	}

	if err != nil {
		log.Printf("%+v\n", err)
		c.AbortWithStatus(500)
		return
	}
}
